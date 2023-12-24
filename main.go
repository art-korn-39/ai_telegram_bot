package main

import (
	"log"
	"os"
	"slices"

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
	WhiteList   = []string{"anastasoid", "Eleng39", "Jokorn", "apolo39"}
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

func main() {

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, _ := Bot.GetUpdatesChan(u)

	// В отдельно горутине обрабатываем информацию по логам
	Logs = make(chan Log, 10)
	go saveLogs()

	// Читаем входящие запросы из канала
	for update := range updates {

		go func(upd tgbotapi.Update) {

			// Запишем panic если горутина завершилась с ошибкой
			defer logPanic()

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
			if User.CheckUserLock(upd) {
				return
			}

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
