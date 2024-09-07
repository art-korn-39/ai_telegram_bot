package main

import "slices"

func gen15_start(user *UserInfo) {

	if Cfg.Gen_Rip && !slices.Contains(Cfg.Admins, user.Username) {
		gen_rip(user)
		return
	}

	msgText := GetText(MsgText_Gemini15Hello, user.Language)
	SendMessage(user, msgText, nil, "")

	msgText = GetText(MsgText_SelectOption, user.Language)
	SendMessage(user, msgText, GetButton(btn_Gen15Types, user.Language), "")

	user.Path = "gen15/type"

}

// После выбора пользователем типа взаимодействия
func gen15_type(user *UserInfo, text string) {

	if gen_DailyLimitOfRequestsIsOver(user, gen15) {
		return
	}

	switch text {

	// НАЧАТЬ ДИАЛОГ
	case GetText(BtnText_StartDialog, user.Language):
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GenEndDialog, user.Language), "")
		user.Path = "gen15/type/dialog"

	// DATA ANALYSIS
	case GetText(BtnText_DataAnalysis, user.Language):
		SendMessage(user, GetText(MsgText_UploadFiles, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "gen15/type/file"

	// Вы можете отправить для обработки: картинку / видео / текстовый файл / аудио / голосовое сообщение

	// ОБРАБОТКА НОВОГО ЗАПРОСА (чат)
	default:
		SendMessage(user, GetText(MsgText_ProcessingRequest, user.Language), GetButton(btn_GenEndDialog, user.Language), "")
		gen15_dialog(user, text)
		user.Path = "gen15/type/dialog"
	}

}
