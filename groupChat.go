package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// ID всей группы: 1087968824
// ChatID: -1002105691474 //общий

func MessageInGroupChat(upd tgbotapi.Update) bool {

	if upd.Message == nil {
		return false
	}

	if upd.Message.From.IsBot {
		if upd.Message.From.ID == 1087968824 {
			return true
		}
	}

	return false

}

func HandleGroupChatMessage(m *tgbotapi.Message) {

	if SubString(m.Text, 0, 5) != "/img " {
		return
	}

	user := NewUserInfo(m)
	prompt := SubString(m.Text, 5, utf8.RuneCountInString(m.Text))

	if strings.TrimSpace(prompt) == "" {
		return
	}

	res, err := SendRequestToKandinsky(prompt, "UHD", user)
	if err != nil {
		Logs <- NewLog(user, "kandinsky{group}", Error, err.Error())
		SendMessage(user, res, nil, "") //GetButton(btn_ImgNewgen, user.Language), "")
	} else {
		caption := fmt.Sprintf(GetText(MsgText_ResultImageGeneration, user.Language), prompt, "4K")
		err := SendPhotoMessage(user, res, caption, nil) //GetButton(btn_ImgNewgenFull, user.Language))
		if err != nil {
			Logs <- NewLog(user, "kandinsky{group}", Error, "{ImgSend} "+err.Error())
			SendMessage(user, GetText(MsgText_ErrorWhileSendingPicture, user.Language), nil, "") //GetButton(btn_ImgNewgen, user.Language), "")
		} else {
			//user.Usage.Kand++
			//Operation := SQL_NewOperation(user, "kandinsky", text, inputText)
			//SQL_AddOperation(Operation)
			//user.Options["image"] = res
		}
	}

}
