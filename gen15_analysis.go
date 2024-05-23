package main

import (
	"fmt"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
)

// https://medium.com/google-cloud/generating-texts-using-files-uploaded-by-gemini-1-5-api-5777f1c902ab
// В настоящее время загрузка текстовых файлов, файлов CSV и PDF приводит к появлению
// сообщения об ошибке типа «Запрос содержит недопустимый аргумент». Похоже, что на
// данном этапе поддерживаются только файлы изображений и фильмов. Мы ожидаем, что это
// ограничение будет устранено в будущем обновлении.

func gen15_file(user *UserInfo, message *tgbotapi.Message) {

	if gen_DailyLimitOfRequestsIsOver(user, gen15) {
		return
	}

	fileID := gen15_GetFileID(message)
	if fileID == "" {
		msgText := GetText(MsgText_UploadFiles, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	if len(user.Gen_LocalFiles) >= 10 {
		return
	}

	// Добавляем счётчик, что началась загрузка файла
	user.WG.Add(1)

	// Сохраняем файл в файловую систему (data/file_ChatID_MsgID.ext)
	name := fmt.Sprintf("file_%d_gen_%d", user.ChatID, message.MessageID)
	filename, err := DownloadFile(fileID, name)
	if err != nil {
		Logs <- NewLog(user, "gen15{file}", Error, err.Error())
		msgText := ProcessErrorOfDownloadingFile(user, err, MsgText_FailedLoadFiles)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		user.WG.Done()
		return
	}

	// блокируем пользователя для изменений, чтобы в других горутинах в мапу ничего параллельно не писалось
	user.Mutex.Lock()

	// инициализируем мапу с файлами (хотя обычно она != nil)
	if user.Gen_LocalFiles == nil {
		user.Gen_LocalFiles = map[int]string{}
	}

	CountOfFiles := len(user.Gen_LocalFiles) // количество уже добавленных
	IsMainGorutine := CountOfFiles == 0      // определяем главную горутину

	user.Gen_LocalFiles[message.MessageID] = filename // указываем в мапе путь до файла

	user.Mutex.Unlock()

	user.WG.Done()

	// Если не основная горутина (первая), то завершаем
	if !IsMainGorutine {
		return
	}

	SendMessage(user, GetText(MsgText_LoadingFiles, user.Language), nil, "")

	// В основной горутине встаём на ожидание, чтобы остальные картинки успели загрузиться
	user.WG.Wait()

	user.Gen_LocalFiles = SortMap(user.Gen_LocalFiles)

	// Просим написать запрос к ним
	msgText := fmt.Sprintf(GetText(MsgText_FilesUploadedWriteText, user.Language), len(user.Gen_LocalFiles))
	SendMessage(user, msgText, GetButton(btn_SendWithoutText, ""), "")

	user.Path = "gen15/type/file/text"

}

// После ввода текста пользователем
func gen15_filetext(user *UserInfo, text string) {

	if gen_DailyLimitOfRequestsIsOver(user, gen15) {
		return
	}

	// Проверяем наличие текста в сообщении
	if text == "" {
		return
	}

	if text == GetText(BtnText_SendWithoutText, user.Language) {
		text = ""
	}

	msgText := GetText(MsgText_ProcessingRequest, user.Language)
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

	<-delay_Gen15

	prompt := []genai.Part{}

	// инициализируем мапу с файлами (хотя обычно она != nil)
	if user.Gen_CloudFiles == nil {
		user.Gen_CloudFiles = []*genai.File{}
	}

	for _, filename := range user.Gen_LocalFiles {

		file, err := gen15_UploadFileToCloudStorage(filename)
		if err != nil {
			Logs <- NewLog(user, "gen15{UploadFile}", Error, err.Error())
			continue
		}

		user.Gen_CloudFiles = append(user.Gen_CloudFiles, file)

	}

	// удаляем локальные файлы
	user.GenDeleteFiles(false)

	// собираем промпт из ссылок на облачные файлы
	for _, file := range user.Gen_CloudFiles {
		prompt = append(prompt, genai.FileData{URI: file.URI})
	}

	if text != "" {
		prompt = append(prompt, genai.Text(text))
	}

	user.Path = "gen15/type/file/text/newgen"

	resp, err := gen15_Model.GenerateContent(gen_ctx, prompt...)
	if err != nil {
		Logs <- NewLog(user, "gen15{GenerateContent}", Error, err.Error())
		msgText := gen15_ProcessErrorsOfResponse(user, err)
		SendMessage(user, msgText, GetButton(btn_Gen15Newgen, user.Language), "")
		return
	}

	if resp.Candidates[0].Content == nil {
		Logs <- NewLog(user, "gen15{file}", Error, "resp.Candidates[0].Content = nil")
		msgText := GetText(MsgText_BadRequest5, user.Language)
		SendMessage(user, msgText, GetButton(btn_Gen15Newgen, user.Language), "")
		return
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	SendMessage(user, string(result), GetButton(btn_Gen15Newgen, user.Language), "")

	user.Usage.Gen15++
	Operation := SQL_NewOperation(user, "gen15", "file", gen15_GetMIME(user), text)
	SQL_AddOperation(Operation)

}

// После ответа пользователя на результат по вопросу и картинкам
func gen15_filetext_newgen(user *UserInfo, text string) {

	if gen_DailyLimitOfRequestsIsOver(user, gen15) {
		return
	}

	switch text {

	// ИЗМЕНИТЬ ЗАПРОС
	case GetText(BtnText_ChangeText, user.Language):
		SendMessage(user, GetText(MsgText_WriteTextToFiles, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "gen15/type/file/text"

	// ЗАГРУЗИТЬ НОВЫЙ ФАЙЛ
	case GetText(BtnText_UploadNewFile, user.Language):
		user.GenDeleteFiles(true) // на всякий почистим, если что-то осталось
		SendMessage(user, GetText(MsgText_UploadFiles, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "gen15/type/file"

	// НАЧАТЬ ДИАЛОГ
	case GetText(BtnText_StartDialog, user.Language):
		user.GenDeleteFiles(true) // на всякий почистим, если что-то осталось
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GenEndDialog, user.Language), "")
		user.Path = "gen15/type/dialog"

	// ОБРАБОТКА НОВОГО ЗАПРОСА
	default:
		// Предполагаем, что там новый вопрос к загруженным картинкам
		gen15_filetext(user, text)
	}

}

// gen_client.UploadFile(gen_ctx, "", f, nil)
// Автоматическое удаление файлов через 2 дня. Максимум 2 ГБ на файл, ограничение 20 ГБ на проект. Загрузка не разрешена.
