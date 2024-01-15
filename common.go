package main

import (
	"fmt"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func start(user *UserInfo, message *tgbotapi.Message) {

	name := message.From.FirstName
	if name == "" {
		name = user.Username
	}

	msgtxt := fmt.Sprintf(GetText(MsgText_Start, user.Language), name, Version)
	SendMessage(user, msgtxt, GetButton(btn_Models, ""), "HTML")

}

func language_start(user *UserInfo) {

	msgText := "Select language / Выберите язык"
	SendMessage(user, msgText, GetButton(btn_Languages, ""), "")

	user.Path = "language/type"

}

// После выбора языка
func language_type(user *UserInfo, text string) {

	switch text {
	case "English":
		user.Language = "en"
	case "Русский":
		user.Language = "ru"
	default:
		msgText := GetText(MsgText_SelectOption, user.Language)
		SendMessage(user, msgText, nil, "")
		return
	}

	msgText := GetText(MsgText_LanguageChanged, user.Language)
	SendMessage(user, msgText, GetButton(btn_Models, user.Language), "")

	user.Path = "start"

}
