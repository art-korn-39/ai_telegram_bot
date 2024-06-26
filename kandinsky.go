package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"
)

// https://fusionbrain.ai/docs/doc/api-dokumentaciya/

// При большой нагрузке или технических работах сервис может быть временно недоступен для приема новых задач.
// Можно заранее проверить текущее состояние с помощью GET запроса на URL /key/api/v1/text2image/availability.
// Во время недоступности задачи не будут приниматься и в ответе на запрос к модели вместо uuid вашего задания
// будет возвращен текущий статус работы сервиса.
// Пример:
// {
//   "model_status": "DISABLED_BY_QUEUE"
// }

var (
	kand_Styles   = map[string]string{"No style": "DEFAULT", "Art": "KANDINSKY", "4K": "UHD", "Anime": "ANIME"}
	kand_Model_id = "4"
)

// После команды "/kandinsky" или при вводе текста = "kandinsky"
func kand_start(user *UserInfo) {

	msgText := GetText(MsgText_EnterDescriptionOfPicture, user.Language)
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

	user.Path = "kandinsky/text"

}

// После ввода запроса пользователем
func kand_text(user *UserInfo, text string) {

	if utf8.RuneCountInString(text) >= 900 {
		SendMessage(user, GetText(MsgText_DescriptionTextNotExceed900Char, user.Language), nil, "")
		return
	}

	user.Options["text"] = text

	msgText := GetText(MsgText_SelectStyleForImage, user.Language)
	SendMessage(user, msgText, GetButton(btn_KandStyles, ""), "")

	user.Path = "kandinsky/text/style"

}

// После выбора стиля пользователем
func kand_style(user *UserInfo, text string) {

	style, ok := kand_Styles[text]
	if !ok {
		msgText := GetText(MsgText_SelectStyleFromOptions, user.Language)
		SendMessage(user, msgText, GetButton(btn_KandStyles, ""), "")
		return
	}

	user.Options["style"] = style
	inputText := user.Options["text"]

	msgText := GetText(MsgText_ImageGenerationStarted1, user.Language)
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

	<-delay_Kandinsky

	res, err := SendRequestToKandinsky(inputText, style, user)
	if err != nil {
		Logs <- NewLog(user, "kandinsky", Error, err.Error())
		SendMessage(user, res, GetButton(btn_ImgNewgen, user.Language), "")
	} else {
		caption := fmt.Sprintf(GetText(MsgText_ResultImageGeneration, user.Language), inputText, text)
		err := SendPhotoMessage(user, res, caption, GetButton(btn_ImgNewgenFull, user.Language))
		if err != nil {
			Logs <- NewLog(user, "kandinsky", Error, "{ImgSend} "+err.Error())
			SendMessage(user, GetText(MsgText_ErrorWhileSendingPicture, user.Language), GetButton(btn_ImgNewgen, user.Language), "")
		} else {
			user.Usage.Kand++
			Operation := SQL_NewOperation(user, "kandinsky", "image", text, inputText)
			SQL_AddOperation(Operation)
			user.Options["image"] = res
		}
	}

	user.Path = "kandinsky/text/style/newgen"

}

// После получения результата генерации
func kand_newgen(user *UserInfo, text string) {

	switch text {

	// ИЗМЕНИТЬ ЗАПРОС
	case GetText(BtnText_ChangeQuerryText, user.Language):
		SendMessage(user, GetText(MsgText_EnterDescriptionOfPicture, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "kandinsky/text"

	// ИЗМЕНИТЬ СТИЛЬ
	case GetText(BtnText_ChooseAnotherStyle, user.Language):
		SendMessage(user, GetText(MsgText_SelectStyleForImage, user.Language), GetButton(btn_KandStyles, ""), "")
		user.Path = "kandinsky/text/style"

	// UPSCALE
	case GetText(BtnText_Upscale, user.Language):
		sdxl_upscale(user, btn_ImgNewgen)
		user.Path = "kandinsky/text/style/newgen"

	// ОБРАБОТКА НОВОГО ЗАПРОСА
	default:
		// Предполагаем, что там новый запрос
		user.Path = "kandinsky/text"
		kand_text(user, text)
	}

}

// Максимальный размер текстового описания - 1000 символов.
// "censored": true (возвращая ошибку, когда запрос или изображение не соответствует)
func SendRequestToKandinsky(text string, style string, user *UserInfo) (result string, err error) {

	scriptPath := WorkDir + "/scripts/generate_image.py"
	dataFolder := WorkDir + "/data"

	cmd := exec.Command(`python`,
		scriptPath,
		dataFolder,
		text,
		style,
		strconv.Itoa(int(user.ChatID)),
		kand_Model_id,
		Cfg.Kandinsky_Key,
		Cfg.Kandinsky_Secret,
	)

	if cmd.Err != nil {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, cmd.Err.Error())
		err = errors.New("{cmd} " + description)
		return GetText(MsgText_FailedGenerateImage1, user.Language), err
	}

	// Получение результата команды
	res, err2 := cmd.Output()

	if err2 != nil {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, err2.Error())
		err = errors.New("{cmd.Output()} " + description)
		return GetText(MsgText_FailedGenerateImage2, user.Language), err
	}

	pathToImage := strings.TrimSpace(string(res[:]))

	if pathToImage == "" {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, "скрипт вернул пустой путь до картинки (не успел получить ответ по api)")
		err = errors.New("{API} " + description)
		return GetText(MsgText_FailedGenerateImage1, user.Language), err
	}

	return pathToImage, nil

}

// func kand_CheckAvailability() bool {

// }
