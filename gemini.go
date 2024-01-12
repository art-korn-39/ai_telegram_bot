package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

//https://ai.google.dev/tutorials/go_quickstart?hl=ru
//https://ai.google.dev/models/gemini?hl=ru

// Gemini_APIKEY = "AIzaSyC0myz4bPIDyx6pPtW0PBZqmJW37A5VJ_k"

// - FinishReasonSafety означает, что потенциальное содержимое было помечено по соображениям безопасности.
// - BlockReasonSafety означает, что промт был заблокирован по соображениям безопасности. Вы можете проверить
// `safety_ratings`, чтобы понять, какая категория безопасности заблокировала его.

var (
	ctx_Gemini    context.Context
	client_Gemini *genai.Client
	model_Gemini  *genai.GenerativeModel
)

func NewConnectionGemini() {
	ctx_Gemini = context.Background()
	client_Gemini, _ = genai.NewClient(ctx_Gemini, option.WithAPIKey(Cfg.GeminiKey))
	model_Gemini = client_Gemini.GenerativeModel("gemini-pro")
}

// После команды "/gemini" или при вводе текста = "gemini"
func gen_start(user *UserInfo) {

	msgText := `Вас приветствует Gemini Pro от компании Google 🚀`
	//На текущий момент я умею вести диалог и отвечать на вопросы по картинкам.`
	//В отличии от моего конкурента (ChatGPT) - у меня нет ограничений на использование, так что можешь развлекаться сколько пожелаешь 😎`

	SendMessage(user, msgText, nil, "")

	msgText = `Выберите один из предложенных вариантов:`
	SendMessage(user, msgText, buttons_genTypes, "")

	user.Path = "gemini/type"

}

// После выбора пользователем типа взаимодействия
func gen_type(user *UserInfo, text string) {

	if gen_DailyLimitOfRequestsIsOver(user) {
		return
	}

	switch text {
	case "Начать диалог":
		SendMessage(user, "Привет! Чем могу помочь?", buttons_genEndDialog, "")
		user.Path = "gemini/type/dialog"
	case "Отправить картинку с текстом":
		SendMessage(user, "Загрузите одну или несколько картинок", button_RemoveKeyboard, "")
		user.Path = "gemini/type/image"
	default:
		gen_dialog(user, text)
		user.Path = "gemini/type/dialog"
	}

}

// После ввода сообщения пользователем
func gen_dialog(user *UserInfo, text string) {

	if text == "Завершить диалог" {
		user.Messages_Gemini = []*genai.Content{}
		SendMessage(user, `Выберите один из предложенных вариантов:`, buttons_genTypes, "")
		user.Path = "gemini/type"
		return
	}

	if gen_DailyLimitOfRequestsIsOver(user) {
		return
	}

	<-delay_Gemini

	user.Requests_today_gen++

	Operation := SQL_NewOperation(user, "gemini", "dialog", text)
	SQL_AddOperation(Operation)

	var msgText string
	cs := model_Gemini.StartChat()
	cs.History = user.Messages_Gemini

	resp, err := cs.SendMessage(ctx_Gemini, genai.Text(text))
	if err != nil {
		errorString := err.Error()
		Logs <- NewLog(user, "gemini", Error, errorString)

		if errorString == "blocked: candidate: FinishReasonSafety" {

			// В случае данного вида ошибки - запускаем новый клиент соединения
			NewConnectionGemini()
			msgText = "Не удалось получить ответ от сервиса. Попробуйте изменить текст вопроса или начать новый диалог."

		} else if errorString == "blocked: prompt: BlockReasonSafety" {

			msgText = "Запрос был заблокирован по соображениям безопасности. Попробуйте изменить текст запроса."

		} else if errorString == "googleapi: Error 500:" {

			// Отправляем сообщение повторно
			time.Sleep(time.Millisecond * 200)
			Logs <- NewLog(user, "gemini", Error, "Повторная отправка запроса ...")
			resp, err = cs.SendMessage(ctx_Gemini, genai.Text(text))
			if err != nil {
				msgText = "Произошла непредвиденная ошибка. Попробуйте позже."
			}

		} else {
			msgText = "Произошла непредвиденная ошибка. Попробуйте позже."
		}

		// Отправляем сообщение и завершаем процедуру если получили ошибку в ответ
		if err != nil {
			SendMessage(user, msgText, nil, "")
			return
		}
	}

	if resp.Candidates[0].Content == nil {
		Logs <- NewLog(user, "gemini", Error, "resp.Candidates[0].Content = nil")
		msgText = "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса."
		SendMessage(user, msgText, nil, "")
		return
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	history := append(user.Messages_Gemini,
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

	user.Messages_Gemini = history

	msgText = string(result)
	SendMessage(user, msgText, nil, "")

}

// После отправки картинок пользователем
func gen_image(user *UserInfo, message *tgbotapi.Message) {

	if gen_DailyLimitOfRequestsIsOver(user) {
		return
	}

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
		Logs <- NewLog(user, "gemini", Error, err.Error())
		msgText := "Не удалось загрузить изображение, попробуйте ещё раз."
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

	if gen_DailyLimitOfRequestsIsOver(user) {
		return
	}

	// Проверяем наличие текста в сообщении
	if text == "" {
		msgText := "Напишите свой вопрос к загруженным изображениям."
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		return
	}

	<-delay_Gemini

	user.Requests_today_gen++

	Operation := SQL_NewOperation(user, "gemini", "img", text)
	SQL_AddOperation(Operation)

	model := client_Gemini.GenerativeModel("gemini-pro-vision")

	prompt := []genai.Part{genai.Text(text)}
	for _, v := range user.Images_Gemini {
		imgData, err := os.ReadFile(v)
		if err != nil {
			Logs <- NewLog(user, "gemini", Error, err.Error())
			continue
		}
		prompt = append(prompt, genai.ImageData("jpeg", imgData))
	}

	resp, err := model.GenerateContent(ctx_Gemini, prompt...)

	if err != nil {
		Logs <- NewLog(user, "gemini{img}", Error, err.Error())
		msgText := "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса или использовать другие изображения."
		SendMessage(user, msgText, buttons_genNewgen, "")
		user.Path = "gemini/type/image/text/newgen"
		return
	}

	if resp.Candidates[0].Content == nil {
		Logs <- NewLog(user, "gemini{img}", Error, "resp.Candidates[0].Content = nil")
		msgText := "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса или использовать другие изображения."
		SendMessage(user, msgText, buttons_genNewgen, "")
		user.Path = "gemini/type/image/text/newgen"
		return
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	SendMessage(user, string(result), buttons_genNewgen, "")

	user.Path = "gemini/type/image/text/newgen"

}

// После ответа пользователя на результат по вопросу и картинкам
func gen_imgtext_newgen(user *UserInfo, text string) {

	if gen_DailyLimitOfRequestsIsOver(user) {
		return
	}

	switch text {
	case "Изменить текст вопроса":
		SendMessage(user, "Напишите свой вопрос к загруженным изображениям.", button_RemoveKeyboard, "")
		user.Path = "gemini/type/image/text"
	case "Загрузить новые изображения":
		user.DeleteImages() // на всякий почистим, если что-то осталось
		SendMessage(user, "Загрузите одну или несколько картинок.", button_RemoveKeyboard, "")
		user.Path = "gemini/type/image"
	case "Начать диалог":
		user.DeleteImages() // на всякий почистим, если что-то осталось
		SendMessage(user, "Привет! Чем могу помочь?", buttons_genEndDialog, "")
		user.Path = "gemini/type/dialog"
	default:
		// Предполагаем, что там новый вопрос к загруженным картинкам
		gen_imgtext(user, text)
	}

}
