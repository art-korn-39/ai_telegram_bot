package main

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

//добавил текст в локал
//условие в main
//убрал кнопку
//заменил на 1.5
//добавил задержку по запросам

// https://ai.google.dev/tutorials/go_quickstart?hl=ru
// https://ai.google.dev/models/gemini?hl=ru

// - FinishReasonSafety означает, что потенциальное содержимое было помечено по соображениям безопасности.
// - BlockReasonSafety означает, что промт был заблокирован по соображениям безопасности. Вы можете проверить
// `safety_ratings`, чтобы понять, какая категория безопасности заблокировала его.

var (
	gen_ctx                 context.Context
	gen_client              *genai.Client
	gen_TextModel           *genai.GenerativeModel
	gen_TextModelWithCensor *genai.GenerativeModel
	gen_ImageModel          *genai.GenerativeModel
)

func NewConnectionGemini() {

	gen_ctx = context.Background()
	gen_client, _ = genai.NewClient(gen_ctx, option.WithAPIKey(Cfg.GeminiKey))
	gen_TextModel = gen_client.GenerativeModel("gemini-1.0-pro") //gemini-1.5-pro-latest
	gen_TextModelWithCensor = gen_client.GenerativeModel("gemini-1.0-pro")
	gen_ImageModel = gen_client.GenerativeModel("gemini-pro-vision")

	// 1 - блокировать всё
	// 2 - допускается с незначимым и низким
	// 3 - допускается незначительные, низкие и средние значения
	// 4 - не блокировать совсем
	SafetySettings := []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment, // домогательство, преследование
			Threshold: genai.HarmBlockNone,          // 4
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit, // откровенно сексуального характера
			Threshold: genai.HarmBlockNone,                // 4
		},
		{
			Category:  genai.HarmCategoryHateSpeech,  // разжигание ненависти
			Threshold: genai.HarmBlockMediumAndAbove, // 2
		},
		{
			Category:  genai.HarmCategoryDangerousContent, // опасный контент
			Threshold: genai.HarmBlockMediumAndAbove,      // 2
		},
	}

	//	gen_TextModel.SafetySettings = SafetySettings
	//	gen_ImageModel.SafetySettings = SafetySettings

	Unused(SafetySettings)

}

// После команды "/gemini" или при вводе текста = "gemini"
func gen_start(user *UserInfo) {

	msgText := GetText(MsgText_GeminiHello, user.Language)
	SendMessage(user, msgText, nil, "")

	msgText = GetText(MsgText_SelectOption, user.Language)
	SendMessage(user, msgText, GetButton(btn_GenTypes, user.Language), "")

	user.Path = "gemini/type"

}

// После команды "/gemini" или при вводе текста = "gemini"
func gen_rip(user *UserInfo) {

	msgText := GetText(MsgText_GeminiRIP, user.Language)
	SendMessage(user, msgText, GetButton(btn_Models, ""), "")

	user.Path = ""

}

// После выбора пользователем типа взаимодействия
func gen_type(user *UserInfo, text string) {

	if gen_DailyLimitOfRequestsIsOver(user) {
		return
	}

	switch text {

	// НАЧАТЬ ДИАЛОГ
	case GetText(BtnText_StartDialog, user.Language):
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GenEndDialog, user.Language), "")
		user.Path = "gemini/type/dialog"

	// AI VISION
	case GetText(BtnText_SendPictureWithText, user.Language):
		SendMessage(user, GetText(MsgText_UploadImages, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "gemini/type/image"

	// ОБРАБОТКА НОВОГО ЗАПРОСА (чат)
	default:
		SendMessage(user, GetText(MsgText_ProcessingRequest, user.Language), GetButton(btn_GenEndDialog, user.Language), "")
		gen_dialog(user, text)
		user.Path = "gemini/type/dialog"
	}

}
