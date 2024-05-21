package main

import (
	"fmt"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
)

// После отправки картинок пользователем
func gen_image(user *UserInfo, message *tgbotapi.Message) {

	if gen_DailyLimitOfRequestsIsOver(user, gen10) {
		return
	}

	// Проверяем наличие картинок в сообщении
	if message.Photo == nil {
		msgText := GetText(MsgText_UploadImages, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	// Потом лучше переделать, а то могут быть баги
	if len(user.Gen_LocalFiles) >= 10 {
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
		Logs <- NewLog(user, "gemini{img}", Error, err.Error())
		msgText := GetText(MsgText_FailedLoadImages, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		user.WG.Done()
		return
	}

	// Собираем в map, где key = MsgID, value = путь к файлу
	// Тут важно собрать горутины в очередь, чтобы самая первая стала основной
	user.Mutex.Lock()
	if user.Gen_LocalFiles == nil {
		// инициализируем мапу с файлами картинок (хотя обычно она != nil)
		user.Gen_LocalFiles = map[int]string{}
	}
	ImageNumber := len(user.Gen_LocalFiles) // количество уже добавленных
	IsMainGorutine := ImageNumber == 0      // определяем главную горутину

	// не нужно +++
	//	newName := fmt.Sprintf("img_%d_gen_%d", user.ChatID, ImageNumber) // создаем новое имя с индексом в массиве фото
	//	newFilename := strings.ReplaceAll(filename, name, newName)        // получаем полный путь с новым именем
	//	os.Rename(filename, newFilename)                                  // заменяем имя у уже созданного
	//	user.Gen_LocalFiles[message.MessageID] = newFilename               // указываем в мапе новый путь до файла
	// не нужно ---

	user.Gen_LocalFiles[message.MessageID] = filename // указываем в мапе путь до файла

	user.Mutex.Unlock()

	user.WG.Done()

	// Если не основная горутина (первая), то завершаем
	if !IsMainGorutine {
		return
	}

	SendMessage(user, GetText(MsgText_LoadingImages, user.Language), nil, "")

	// В основной горутине встаём на ожидание, чтобы остальные картинки успели загрузиться
	user.WG.Wait()

	user.Gen_LocalFiles = SortMap(user.Gen_LocalFiles)

	// Просим написать запрос к ним
	msgText := fmt.Sprintf(GetText(MsgText_PhotosUploadedWriteQuestion, user.Language), len(user.Gen_LocalFiles))
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

	user.Path = "gemini/type/image/text"

}

// После ввода вопроса пользователем
func gen_imgtext(user *UserInfo, text string) {

	if gen_DailyLimitOfRequestsIsOver(user, gen10) {
		return
	}

	// Проверяем наличие текста в сообщении
	if text == "" {
		msgText := GetText(MsgText_WriteQuestionToImages, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	<-delay_Gemini

	prompt := []genai.Part{genai.Text(text)}
	for _, v := range user.Gen_LocalFiles {
		imgData, err := os.ReadFile(v)
		if err != nil {
			Logs <- NewLog(user, "gemini", Error, err.Error())
			continue
		}
		prompt = append(prompt, genai.ImageData("jpeg", imgData))
	}

	resp, err := gen_ImageModel.GenerateContent(gen_ctx, prompt...)

	if err != nil {
		Logs <- NewLog(user, "gemini{img}", Error, err.Error())
		msgText := GetText(MsgText_BadRequest1, user.Language)
		SendMessage(user, msgText, GetButton(btn_GenNewgen, user.Language), "")
		user.Path = "gemini/type/image/text/newgen"
		return
	}

	if resp.Candidates[0].Content == nil {
		Logs <- NewLog(user, "gemini{img}", Error, "resp.Candidates[0].Content = nil")
		msgText := GetText(MsgText_BadRequest1, user.Language)
		SendMessage(user, msgText, GetButton(btn_GenNewgen, user.Language), "")
		user.Path = "gemini/type/image/text/newgen"
		return
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	SendMessage(user, string(result), GetButton(btn_GenNewgen, user.Language), "")

	user.Usage.Gen10++
	Operation := SQL_NewOperation(user, "gemini", "img", text)
	SQL_AddOperation(Operation)

	user.Path = "gemini/type/image/text/newgen"

}

// После ответа пользователя на результат по вопросу и картинкам
func gen_imgtext_newgen(user *UserInfo, text string) {

	if gen_DailyLimitOfRequestsIsOver(user, gen10) {
		return
	}

	switch text {

	// ИЗМЕНИТЬ ЗАПРОС
	case GetText(BtnText_ChangeQuerryText, user.Language):
		SendMessage(user, GetText(MsgText_WriteQuestionToImages, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "gemini/type/image/text"

	// ЗАГРУЗИТЬ НОВЫЕ КАРТИНКИ
	case GetText(BtnText_UploadNewImages, user.Language):
		user.GenDeleteFiles(false) // на всякий почистим, если что-то осталось
		SendMessage(user, GetText(MsgText_UploadImages, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "gemini/type/image"

	// НАЧАТЬ ДИАЛОГ
	case GetText(BtnText_StartDialog, user.Language):
		user.GenDeleteFiles(false) // на всякий почистим, если что-то осталось
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GenEndDialog, user.Language), "")
		user.Path = "gemini/type/dialog"

	// ОБРАБОТКА НОВОГО ЗАПРОСА
	default:
		// Предполагаем, что там новый вопрос к загруженным картинкам
		gen_imgtext(user, text)
	}

}
