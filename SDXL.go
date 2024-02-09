package main

import (
	"fmt"
	"unicode/utf8"
)

// После команды "/sdxl" или при вводе текста = "SDXL"
func sdxl_start(user *UserInfo) {

	if sdxl_DailyLimitOfRequestsIsOver(user, btn_Models) {
		return
	}

	msgText := sdxl_WelcomeTextMessage(user)
	SendMessage(user, msgText, nil, "HTML")

	msgText = GetText(MsgText_EnterDescriptionOfPicture, user.Language)
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

	user.Path = "sdxl/text"

}

// После ввода запроса пользователем
func sdxl_text(user *UserInfo, text string) {

	if sdxl_DailyLimitOfRequestsIsOver(user, btn_Models) {
		return
	}

	if utf8.RuneCountInString(text) >= 2000 {
		SendMessage(user, GetText(MsgText_DescriptionTextNotExceed2000Char, user.Language), nil, "")
		return
	}

	eng_text, err := gen_TranslateToEnglish(text)
	if err != nil {
		Logs <- NewLog(user, "SDXL", Error, err.Error())
		msgText := GetText(MsgText_ErrorTranslatingIntoEnglish, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		msgText = GetText(MsgText_EnterDescriptionOfPicture, user.Language)
		SendMessage(user, msgText, nil, "")
		user.Path = "sdxl/text"
		return
	}

	user.Options["text"] = eng_text

	msgText := GetText(MsgText_SelectStyleForImage, user.Language)
	SendMessage(user, msgText, GetButton(btn_SDXLStyles, ""), "")

	user.Path = "sdxl/text/style"

}

// STYLE
//3d-model analog-film anime cinematic
//comic-book digital-art enhance fantasy-art
//isometric line-art low-poly modeling-compound
//neon-punk origami photographic pixel-art tile-texture

// После выбора стиля пользователем
func sdxl_style(user *UserInfo, text string) {

	if sdxl_DailyLimitOfRequestsIsOver(user, btn_Models) {
		return
	}

	style, ok := SDXL_Styles[text]
	if !ok {
		msgText := GetText(MsgText_SelectStyleFromOptions, user.Language)
		SendMessage(user, msgText, GetButton(btn_SDXLStyles, ""), "")
		return
	}

	user.Options["style"] = style
	user.Options["styleName"] = text

	msgText := GetText(MsgText_ImageGenerationStarted2, user.Language)
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

	res, err := sdxl_Text2image(user)
	if err != nil {
		Logs <- NewLog(user, "SDXL", Error, err.Error())
		SendMessage(user, res, GetButton(btn_ImgNewgen, user.Language), "")
	} else {
		caption := fmt.Sprintf(GetText(MsgText_ResultImageGeneration, user.Language), user.Options["text"], user.Options["styleName"])
		err = SendPhotoMessage(user, res, caption, GetButton(btn_ImgNewgenFull, user.Language))
		if err != nil {
			Logs <- NewLog(user, "SDXL", Error, "{ImgSend} "+err.Error())
			SendMessage(user, GetText(MsgText_UnexpectedError, user.Language), GetButton(btn_ImgNewgen, user.Language), "")
		} else {
			user.Requests_today_sdxl++
			Operation := SQL_NewOperation(user, "sdxl", text, user.Options["text"])
			SQL_AddOperation(Operation)
			user.Options["image"] = res
		}
	}

	user.Path = "sdxl/text/style/newgen"

}

func sdxl_Upscale(user *UserInfo) {

	if sdxl_DailyLimitOfRequestsIsOver(user, 0) {
		return
	}

	filepath, ok := user.Options["image"]
	if !ok {
		Logs <- NewLog(user, "SDXL", Error, `{Upscale} пусто в user.Options["image"]`)
		msgText := GetText(MsgText_NoImageFoundToProcess, user.Language)
		SendMessage(user, msgText, GetButton(btn_ImgNewgen, ""), "")
		return
	}

	msgText := GetText(MsgText_ImageProcessingStarted, user.Language)
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

	res, err := sdxl_Image2ImageUpscale(user, filepath)
	if err != nil {
		if err.Error() != "CONTENT_FILTERED" {
			Logs <- NewLog(user, "SDXL", Error, err.Error())
		}
		SendMessage(user, res, GetButton(btn_ImgNewgen, user.Language), "")
	} else {
		err = SendFileMessage(user, res, "", GetButton(btn_ImgNewgen, user.Language))
		if err != nil {
			Logs <- NewLog(user, "SDXL", Error, "{FileSend} "+err.Error())
			SendMessage(user, GetText(MsgText_UnexpectedError, user.Language), GetButton(btn_ImgNewgen, user.Language), "")
		} else {
			user.Requests_today_sdxl++
			Operation := SQL_NewOperation(user, "sdxl", "upscale", user.Options["text"])
			SQL_AddOperation(Operation)
		}
	}

	//path задается выше по стеку, т.к. метод используется в кандинском и sdxl
}

// После получения результата генерации
func sdxl_newgen(user *UserInfo, text string) {

	if sdxl_DailyLimitOfRequestsIsOver(user, btn_Models) {
		return
	}

	switch text {
	case GetText(BtnText_ChangeQuerryText, user.Language):
		SendMessage(user, GetText(MsgText_EnterDescriptionOfPicture, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "sdxl/text"
	case GetText(BtnText_ChooseAnotherStyle, user.Language):
		SendMessage(user, GetText(MsgText_SelectStyleForImage, user.Language), GetButton(btn_SDXLStyles, ""), "")
		user.Path = "sdxl/text/style"
	case GetText(BtnText_Upscale, user.Language):
		sdxl_Upscale(user)
		user.Path = "sdxl/text/style/newgen"
	default:
		// Предполагаем, что там новый запрос
		user.Path = "sdxl/text"
		sdxl_text(user, text)
	}

}
