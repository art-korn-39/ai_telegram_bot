package main

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

//https://ai.google.dev/tutorials/go_quickstart?hl=ru
//https://ai.google.dev/models/gemini?hl=ru

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

	msgText := GetText(MsgText_GeminiHello, user.Language)
	SendMessage(user, msgText, nil, "")

	msgText = GetText(MsgText_SelectOption, user.Language)
	SendMessage(user, msgText, GetButton(btn_GenTypes, user.Language), "")

	user.Path = "gemini/type"

}

// После выбора пользователем типа взаимодействия
func gen_type(user *UserInfo, text string) {

	if gen_DailyLimitOfRequestsIsOver(user) {
		return
	}

	switch text {
	case GetText(BtnText_StartDialog, user.Language):
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GenEndDialog, user.Language), "")
		user.Path = "gemini/type/dialog"
	case GetText(BtnText_SendPictureWithText, user.Language):
		SendMessage(user, GetText(MsgText_UploadImages, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "gemini/type/image"
	default:
		gen_dialog(user, text)
		user.Path = "gemini/type/dialog"
	}

}
