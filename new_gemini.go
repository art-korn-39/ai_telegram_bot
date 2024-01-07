package main

import (
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
)

var (
	buttons_geminiTypes = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Начать диалог"),
			tgbotapi.NewKeyboardButton("Отправить картинку с текстом"),
		),
	)
	buttons_geminiNewgen = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить текст вопроса"),
			tgbotapi.NewKeyboardButton("Загрузить новые изображения"),
			tgbotapi.NewKeyboardButton("Начать диалог"),
		),
	)
	buttons_geminiEndDialog = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Завершить диалог")),
	)
)

// После команды "/gemini" или при вводе текста = "gemini"
func gen_start(user *UserInfo) {

	msgText := `Выберите один из предложенных вариантов:`
	SendMessage(user, msgText, buttons_geminiTypes, "")

	user.Path = "gemini/type"

}

// После выбора пользователем типа взаимодействия
func gen_type(user *UserInfo, text string) {

	switch text {
	case "Начать диалог":
		SendMessage(user, "Привет! Чем могу помочь?", buttons_geminiEndDialog, "")
		user.Path = "gemini/type/dialog"
	case "Отправить картинку с текстом":
		SendMessage(user, "Загрузите одну или несколько картинок", button_RemoveKeyboard, "")
		user.Path = "gemini/type/image"
	default:
		//SendMessage(user, "Неизвестная команда", buttons_geminiTypes, "")
		gen_dialog(user, text)
	}

}

// После ввода сообщения пользователем
func gen_dialog(user *UserInfo, text string) {

	if text == "Завершить диалог" {
		user.History_Gemini = []*genai.Content{}
		msgText := `Выберите один из предложенных вариантов:`
		SendMessage(user, msgText, buttons_geminiTypes, "")
		user.Path = "gemini/type"
		return
	}

	<-delay_Gemini

	Operation := SQL_NewOperation(user, "gemini", "dialog", text)
	SQL_AddOperation(Operation)

	var msgText string
	cs := model_Gemini.StartChat()
	cs.History = user.History_Gemini

	resp, err := cs.SendMessage(ctx_Gemini, genai.Text(text))
	if err != nil {
		errorString := err.Error()
		Logs <- NewLog(user, "Gemini", Error, errorString)

		if errorString == "blocked: candidate: FinishReasonSafety" {
			NewConnectionGemini() // В случае данного вида ошибки - запускаем новый клиент соединения
			msgText = "Не удалось получить ответ от сервиса. Попробуйте изменить текст вопроса или начать новый диалог."
		} else if errorString == "blocked: prompt: BlockReasonSafety" {
			msgText = "Запрос был заблокирован по соображениям безопасности. Попробуйте изменить текст запроса."
		} else {
			msgText = "Произошла непредвиденная ошибка. Попробуйте позже."
		}
		SendMessage(user, msgText, buttons_geminiEndDialog, "")
		return
	}

	if resp.Candidates[0].Content == nil {
		Logs <- NewLog(user, "Gemini", Error, "resp.Candidates[0].Content = nil")
		msgText = "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса."
		SendMessage(user, msgText, buttons_geminiEndDialog, "")
		return
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	history := append(user.History_Gemini,
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(text),
			},
			Role: "user",
		},
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(result),
			},
			Role: "model",
		},
	)

	user.History_Gemini = history

	msgText = string(result)
	SendMessage(user, msgText, buttons_geminiEndDialog, "")

}

// После отправки картинок пользователем
func gen_image(user *UserInfo, message *tgbotapi.Message) {

	// Проверяем наличие картинок в сообщении
	if message.Photo == nil {
		msgText := "Загрузите одну или несколько картинок."
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		return
	}

	// Потом лучше переделать, а то могут быть баги
	if len(user.Images_Gemini) >= 10 {
		return
	}

	photos := *message.Photo

	// Добавляем счётчик, что началась загрузка фотографии
	user.WG.Add(1)

	// Сохраняем картинку в файловую систему
	// data/img_ChatID_MsgID.jpg
	name := fmt.Sprintf("img_%d_gen_%d", user.ChatID, message.MessageID)
	filename, err := DownloadFile(photos[len(photos)-1].FileID, name)
	if err != nil {
		Logs <- NewLog(user, "Gemini", Error, err.Error())
		msgText := "Не удалось загрузить, попробуйте ещё раз."
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		user.WG.Done()
		return
	}

	// Собираем в map, где key = MsgID, value = путь к файлу
	// Тут важно собрать горутины в очередь, чтобы самая первая стала основной
	user.Mutex.Lock()
	if user.Images_Gemini == nil {
		// инициализируем мапу с файлами картинок (хотя обычно она != nil)
		user.Images_Gemini = map[int]string{}
	}
	ImageNumber := len(user.Images_Gemini)                            // количество уже добавленных
	newName := fmt.Sprintf("img_%d_gen_%d", user.ChatID, ImageNumber) // создаем новое имя с индексом в массиве фото
	newFilename := strings.ReplaceAll(filename, name, newName)        // получаем полный путь с новым именем
	os.Rename(filename, newFilename)                                  // заменяем имя у уже созданного
	IsMainGorutine := ImageNumber == 0                                // определяем главную горутину
	user.Images_Gemini[message.MessageID] = newFilename               // указываем в мапе новый путь до файла
	user.Mutex.Unlock()

	user.WG.Done()

	// Если это последующие горутины, то завершаем их
	if !IsMainGorutine {
		return
	}

	SendMessage(user, "Выполняется загрузка изображений ...", nil, "")

	// В основной горутине встаём на ожидание, чтобы остальные картинки успели загрузиться
	user.WG.Wait()

	user.Images_Gemini = SortMap(user.Images_Gemini)

	// Просим написать запрос к ним
	msgText := fmt.Sprintf(
		`Загружено фотографий: %d
Напишите свой вопрос.
Например:
"Кто на фотографии?"
"Чем отличаются эти картинки?"`, len(user.Images_Gemini))
	SendMessage(user, msgText, button_RemoveKeyboard, "")

	user.Path = "gemini/type/image/text"

}

// После ввода вопроса пользователем
func gen_imgtext(user *UserInfo, text string) {

	// Проверяем наличие текста в сообщении
	if text == "" {
		msgText := "Напишите свой вопрос к загруженным изображениям."
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		return
	}

	Operation := SQL_NewOperation(user, "gemini", "img", text)
	SQL_AddOperation(Operation)

	model := client_Gemini.GenerativeModel("gemini-pro-vision")

	prompt := []genai.Part{genai.Text(text)}
	for _, v := range user.Images_Gemini {
		imgData, err := os.ReadFile(v)
		if err != nil {
			Logs <- NewLog(user, "Gemini", Error, err.Error())
			continue
		}
		prompt = append(prompt, genai.ImageData("jpeg", imgData))
	}

	resp, err := model.GenerateContent(ctx_Gemini, prompt...)

	if err != nil {
		Logs <- NewLog(user, "Gemini{img}", Error, err.Error())
		msgText := "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса или использовать другие изображения."
		SendMessage(user, msgText, buttons_geminiNewgen, "")
		user.Path = "gemini/type/image/text/newgen"
		return
	}

	if resp.Candidates[0].Content == nil {
		Logs <- NewLog(user, "Gemini{img}", Error, "resp.Candidates[0].Content = nil")
		msgText := "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса или использовать другие изображения."
		SendMessage(user, msgText, buttons_geminiNewgen, "")
		user.Path = "gemini/type/image/text/newgen"
		return
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	SendMessage(user, string(result), buttons_geminiNewgen, "")

	user.Path = "gemini/type/image/text/newgen"

}

// После ответа пользователя на результат по вопросу и картинкам
func gen_imgtext_newgen(user *UserInfo, text string) {

	var msgText string

	switch text {
	case "Изменить текст вопроса":
		msgText := "Напишите свой вопрос к загруженным изображениям."
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		user.Path = "gemini/type/image/text"
	case "Загрузить новые изображения":
		user.DeleteImages() // на всякий почистим, если что-то осталось
		SendMessage(user, "Загрузите одну или несколько картинок.", button_RemoveKeyboard, "")
		user.Path = "gemini/type/image"
	case "Начать диалог":
		user.DeleteImages() // на всякий почистим, если что-то осталось
		SendMessage(user, "Привет! Чем могу помочь?", buttons_geminiEndDialog, "")
		user.Path = "gemini/type/dialog"
	default:
		// Предполагаем, что там новый вопрос к загруженным картинкам
		gen_imgtext(user, text)
	}

	SendMessage(user, msgText, nil, "")

}
