package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	WebHook           = "https://telegram-ai-bot-art-korn-39.amvera.io/"
	TelegramBotToken  = "6894745722:AAHMfRNCNM_T9P5-zUgV4OU_DLfzOOUEFmg"
	ChannelChatID     = -1001997602646
	ChannelURL        = "https://t.me/+6ZMACWRgFdRkNGEy"
	CheckSubscription = true
)

var (
	Bot         *tgbotapi.BotAPI
	Logs        chan Log
	ListOfUsers = map[int64]*UserInfo{}
	WhiteList   = []string{"anastasoid", "Eleng39", "Jokorn"}
	arrayCMD    = []string{"gemini", "kandinsky", "chatgpt"}
)

//ID chat (my) = 403059287
//ID chat (my second) = 609614322

// Реализация через webhook
// if _, err := bot.SetWebhook(tgbotapi.NewWebhook(WebHook)); err != nil {
// 	log.Fatalf("setting webhook %v; error: %v", WebHook, err)
// }

//bot.ListenForWebhook("/")

func init() {

	// используя токен создаем новый инстанс бота
	b, err := tgbotapi.NewBotAPI(TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	Bot = b

}

//go:embed scripts/generate_image.py
var script_py embed.FS

func main() {

	// data, _ := script_py.ReadFile("scripts/generate_image.py")
	// fmt.Println(string(data))

	log.Printf("added by git")
	log.Printf("Authorized on account %s", Bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, _ := Bot.GetUpdatesChan(u)

	// В отдельно горутине обрабатываем информацию по логам
	Logs = make(chan Log, 10)
	go func() {
		for v := range Logs {
			log.Printf("[%s] %s", v.UserName, v.Text)
			if v.IsError { // дополнительно: ошибку записываем в файл
				err_file, _ := os.OpenFile("Errors.txt", os.O_APPEND|os.O_WRONLY, 0600)
				err_file.WriteString(fmt.Sprintf("[%s] %s\n", v.UserName, v.Text))
				err_file.Close()
			}
		}
	}()

	// Читаем входящие запросы из канала
	for update := range updates {

		go func(upd tgbotapi.Update) {

			if upd.Message == nil {
				return
			}

			// Проверка подписки пользователя на канал
			if !accessIsAllowed(upd) {
				return
			}

			// Получаем данные пользователя
			User, ok := ListOfUsers[upd.Message.Chat.ID]
			if !ok {
				User = &UserInfo{}
				ListOfUsers[upd.Message.Chat.ID] = User
			}
			defer User.SetIsRunning(false)

			// Фиксируем пользователя и входящее сообщение
			Logs <- Log{upd.Message.From.UserName, upd.Message.Text, false}

			// Если предыдущий запрос ещё выполняется, то новые команды не обрабатываем
			if alreadyRunning(User, upd) {
				return
			}

			// Блокируем новые запросы от пользователя
			User.SetIsRunning(true)

			// Обработка запроса от пользователя
			var result ResultOfRequest
			if MsgIsCommand(upd.Message) {
				cmd := MsgCommand(upd.Message)
				User.LastCommand = cmd
				result = processCommand(cmd, upd)
			} else {
				result = processText(upd.Message.Text, User.LastCommand, upd)
			}

			// Отправка сообщения
			Bot.Send(result.Message)
			Logs <- Log{result.Log_author, result.Log_message, false}

		}(update)
	}
}

func processCommand(cmd string, upd tgbotapi.Update) ResultOfRequest {

	var result ResultOfRequest
	result.Log_author = "bot"

	switch cmd {
	case "start":
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, start(upd.Message.Chat.FirstName))
		msg.ParseMode = "HTML"

		var buttons = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Gemini"),
				tgbotapi.NewKeyboardButton("Kandinsky"),
				tgbotapi.NewKeyboardButton("ChatGPT"),
			),
		)
		msg.ReplyMarkup = buttons

		result.Message = msg
		result.Log_message = "/start for " + upd.Message.Chat.UserName
	case "stop":
		if upd.Message.From.UserName == "Art_Korn_39" {
			os.Exit(1)
		}
	case "chatgpt":
		msg_text := "Напишите свой вопрос:"
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_message = msg_text
	case "gemini":
		msg_text := "Напишите свой вопрос:"
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_message = msg_text
	case "kandinsky":
		msg_text := "Введите свой запрос:"
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_message = msg_text
	}

	return result

}

func processText(text string, cmd string, upd tgbotapi.Update) ResultOfRequest {

	var result ResultOfRequest

	switch cmd {
	case "chatgpt":
		msg_text := SendRequestToChatGPT(upd.Message.Text)
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "ChatGPT"
		result.Log_message = msg_text

	case "gemini":
		msg_text := SendRequestToGemini(upd.Message.Text)
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "Gemini"
		result.Log_message = msg_text

	case "kandinsky":
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Запущена генерация картинки, она может занять 1-2 минуты.")
		Bot.Send(msg)

		pathToImage, err := SendRequestToKandinsky(upd.Message.Text)
		if err != nil {
			result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, "Не удалось сгенерировать изображение. Попробуйте позже.")
			result.Log_author = "Kandinsky"
			result.Log_message = "Ошибка при генерации картинки."
			Logs <- Log{"Kandinsky", err.Error(), true}
		} else {
			result.Message = tgbotapi.NewPhotoUpload(upd.Message.Chat.ID, pathToImage)
			result.Log_author = "Kandinsky"
			result.Log_message = pathToImage
		}
	case "":
		msg_text := "Не выбрана нейросеть для обработки запроса."
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "bot"
		result.Log_message = msg_text

	case "start":
		msg_text := "Не выбрана нейросеть для обработки запроса."
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "bot"
		result.Log_message = msg_text
	}

	return result

}

func accessIsAllowed(upd tgbotapi.Update) bool {

	if !CheckSubscription {
		return true
	}

	if slices.Contains(WhiteList, upd.Message.Chat.UserName) {
		return true
	}

	result := true

	conf := tgbotapi.ChatConfigWithUser{ChatID: ChannelChatID, UserID: int(upd.Message.Chat.ID)}
	chatMember, err := Bot.GetChatMember(conf)
	if err != nil {
		Logs <- Log{"bot{GetChatMember}", err.Error(), true}
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Произошла непредвиденная ошибка, попробуйте позже.")
		Bot.Send(msg)
		result = false
	}

	if chatMember.Status != "member" && upd.Message.Text != "/start" {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Для использования бота необходимо подписаться на канал👇")
		var button = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("✅Подписаться", ChannelURL),
			),
		)
		msg.ReplyMarkup = button
		Bot.Send(msg)
		result = false
	}

	return result

}

func alreadyRunning(User *UserInfo, upd tgbotapi.Update) bool {

	if User.IsRunning {
		warning := "Последняя операция ещё выполняется, дождитесь её завершения перед отправкой новых запросов."
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, warning)
		Bot.Send(msg)
		Logs <- Log{"bot", warning, false}
		return true
	}

	return false

}

func start(user string) string {

	t := `
Привет, %user! 👋

Я бот для работы с нейросетями.
С моей помощью ты можешь использовать следующие модели:

<b>ChatGPT</b> - используется для генерации текста.
<b>Gemini</b> - аналог ChatGPT от компании Google.
<b>Kandinsky</b> - используется для создания изображений по текстовому описанию.

Чтобы начать - просто выбери подходящую нейросеть и задай ей вопрос (или попроси сделать картинку), удачи!🔥`

	return strings.ReplaceAll(t, "%user", user)

}

// var buttons = tgbotapi.NewInlineKeyboardMarkup(
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonData("Gemini", "gemini"),
// 		tgbotapi.NewInlineKeyboardButtonData("Kandinsky", "kandinsky"),
// 		tgbotapi.NewInlineKeyboardButtonData("Chat GPT", "chat_gpt"),
// 	),
// )

// msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
// msg.ReplyMarkup = buttons
// Bot.Send(msg)

// Hello, %user!

// ChatGPT is a neural network that is used to generate text.
// Gemini is an analogue of ChatGPT from Google.
// Kandinsky is a neural network that is used to create images based on text descriptions.

// You can use the bot for various purposes: it will help you write a course or essay, write a poem, draw a picture and answer questions that interest you.
// To get started, select a neural network and simply write down a question.
// --------------------------------------------------------------------------------

// var reply string
// var author string
// command := upd.Message.Command() // комманда - сообщение, начинающееся с "/"
// switch command {
// case "stop":
// 	if upd.Message.From.UserName == "Art_Korn_39" {
// 		return
// 	}
// case "start":
// 	reply = start(upd.Message.Chat.UserName)
// 	author = "bot"
// case "":
// 	// если команды не поступало, то смотрим какая была последней
// 	lastCommand := User.LastCommand
// 	switch lastCommand {
// 	case "":
// 		reply = "Не выбрана нейросеть для обработки запроса."
// 		author = "bot"
// 	case "chat_gpt":
// 		reply = SendRequestToChatGPT(upd.Message.Text)
// 		author = "ChatGPT"
// 	case "gemini":
// 		reply = SendRequestToGemini_GetText(upd.Message.Text)
// 		author = "Gemini"
// 	case "kandinsky":
// 		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Запущена генерация картинки, она может занять 1-2 минуты.")
// 		Bot.Send(msg)
// 		pathToImage, err := SendRequestToKandinsky(upd.Message.Text)
// 		if err != nil {
// 			reply = err.Error()
// 		} else {
// 			PhotoConfig := tgbotapi.NewPhotoUpload(upd.Message.Chat.ID, pathToImage)
// 			Bot.Send(PhotoConfig)
// 			reply = pathToImage
// 		}
// 		author = "Kandinsky"
// 	}
// default:
// 	// сохраняем актуальную команду по юзеру
// 	User.LastCommand = command
// }

// if reply != "" {
// 	if author != "Kandinsky" {
// 		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, reply)
// 		Bot.Send(msg)
// 	}
// 	Logs <- Log{author, reply}
// }
