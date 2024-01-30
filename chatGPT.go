package main

import (
	"github.com/sashabaranov/go-openai"
)

//для оценки токенов
//2 символа RU = 1 токен
//4 символа EN = 1 токен
//utf8.RuneCountInString(str)

var (
	clientOpenAI *openai.Client
	voices       = []string{"fable", "nova", "echo", "onyx"}
)

// После команды "/chatgpt" или при вводе текста = "chatgpt"
func gpt_start(user *UserInfo) {

	msgText := gpt_WelcomeTextMessage(user)
	SendMessage(user, msgText, nil, "HTML")

	msgText = GetText(MsgText_SelectOption, user.Language)
	SendMessage(user, msgText, GetButton(btn_GptTypes, user.Language), "")

	user.Path = "chatgpt/type"

}

// После выбора пользователем типа взаимодействия
func gpt_type(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	switch text {
	case GetText(BtnText_StartDialog, user.Language):
		SendMessage(user, GetText(MsgText_ChatGPTDialogStarted, user.Language), nil, "")
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GptClearContext, user.Language), "")
		user.Path = "chatgpt/type/dialog"
	case GetText(BtnText_GenerateAudioFromText, user.Language):
		SendMessage(user, GetText(MsgText_EnterTextForAudio, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "chatgpt/type/speech_text"
	case GetText(BtnText_SendPictureWithText, user.Language):
		SendMessage(user, GetText(MsgText_UploadImage, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "chatgpt/type/image"
	default:
		SendMessage(user, "Обработка запроса...", GetButton(btn_GptClearContext, user.Language), "")
		gpt_dialog(user, text)
		user.Path = "chatgpt/type/dialog"
	}

}
