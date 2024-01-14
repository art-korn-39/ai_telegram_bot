package main

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/sashabaranov/go-openai"
)

// OpenAIToken = "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" //Коли
// OpenAIToken = "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" //Мой

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
		msgText := GetText(MsgText_ChatGPTDialogStarted, user.Language)
		SendMessage(user, msgText, nil, "")
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GptClearContext, user.Language), "")
		user.Path = "chatgpt/type/dialog"
	case GetText(BtnText_GenerateAudioFromText, user.Language):
		SendMessage(user, GetText(MsgText_EnterTextForAudio, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "chatgpt/type/speech_text"
	default:
		gpt_dialog(user, text)
		user.Path = "chatgpt/type/dialog"
	}

}

// После ввода сообщения пользователем
func gpt_dialog(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	if text == GetText(BtnText_ClearContext, user.Language) {
		user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
		SendMessage(user, GetText(MsgText_DialogContextCleared, user.Language), nil, "")
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), nil, "")
		return
	}

	if text == GetText(BtnText_EndDialog, user.Language) {
		user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
		SendMessage(user, GetText(MsgText_SelectOption, user.Language), GetButton(btn_GptTypes, user.Language), "")
		user.Path = "chatgpt/type"
		return
	}

	<-delay_ChatGPT

	Operation := SQL_NewOperation(user, "chatgpt", "dialog", text)
	SQL_AddOperation(Operation)

	gpt_DialogSendMessage(user, text, true)

}

// После ввода текста пользователем
func gpt_speech_text(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	// Проверяем наличие текста в сообщении
	if text == "" {
		msgText := GetText(MsgText_WriteTextForVoicing, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	length := utf8.RuneCountInString(text)
	if length*20 > (Cfg.TPD_gpt - user.Tokens_used_gpt) {
		msgText := GetText(MsgText_NotEnoughTokensWriteShorterTextLength, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	user.Options["text"] = text

	msgText := GetText(MsgText_VoiceExamples, user.Language)
	SendMessage(user, msgText, GetButton(btn_GptSampleSpeech, ""), "")

	msgText = GetText(MsgText_SelectVoice, user.Language)
	SendMessage(user, msgText, GetButton(btn_GptVoices, ""), "")

	user.Path = "chatgpt/type/speech_text/voice"

}

// После выбора голоса
func gpt_speech_voice(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	if text == "gpt_speech_samples" {
		gpt_SendSpeechSamples(user)
		return
	}

	voice, isErr := gpt_GetSpeechVoice(text)
	if isErr {
		msgText := GetText(MsgText_SelectVoiceFromOptions, user.Language)
		SendMessage(user, msgText, GetButton(btn_GptVoices, ""), "")
		return
	}

	model := openai.TTSModel1HD
	inputText := user.Options["text"]

	Operation := SQL_NewOperation(user, "chatgpt", "speech", text)
	SQL_AddOperation(Operation)

	SendMessage(user, GetText(MsgText_AudioFileCreationStarted, user.Language), nil, "")

	// Отправка запроса в openai
	res, err := clientOpenAI.CreateSpeech(context.Background(),
		openai.CreateSpeechRequest{
			Model: model,
			Input: inputText,
			Voice: voice,
		})

	if err != nil {
		Logs <- NewLog(user, "chatgpt", 1, "{speech1} "+err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		SendMessage(user, msgText, GetButton(btn_GptTypes, user.Language), "")
		user.Path = "chatgpt/type"
		return
	}

	// Сохранение в файл
	filename := fmt.Sprintf(WorkDir+"/data/speech_%d.mp3", user.ChatID)
	err = CreateFile(filename, res)
	if err != nil {
		Logs <- NewLog(user, "chatgpt", 1, "{speech2} "+err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		SendMessage(user, msgText, GetButton(btn_GptTypes, user.Language), "")
		user.Path = "chatgpt/type"
		return
	}

	tokensForRequest := utf8.RuneCountInString(inputText) * 20
	user.Tokens_used_gpt = user.Tokens_used_gpt + tokensForRequest

	caption := fmt.Sprintf(GetText(MsgText_ResultAudioGeneration, user.Language), inputText, voice)
	err = SendAudioMessage(user, filename, caption, GetButton(btn_GptSpeechNewgen, user.Language))
	if err != nil {
		Logs <- NewLog(user, "chatgpt", Error, "{AudioSend} "+err.Error())
		SendMessage(user, GetText(MsgText_ErrorSendingAudioFile, user.Language), nil, "")
	}

	user.Path = "chatgpt/type/speech_text/voice/newgen"

}

// После получения результата генерации
func gpt_speech_newgen(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	switch text {
	case GetText(BtnText_ChangeText, user.Language):
		SendMessage(user, GetText(MsgText_WriteTextForVoicing, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "chatgpt/type/speech_text"
	case GetText(BtnText_ChooseAnotherVoice, user.Language):
		SendMessage(user, GetText(MsgText_SelectVoice, user.Language), GetButton(btn_GptVoices, ""), "")
		user.Path = "chatgpt/type/speech_text/voice"
	case GetText(BtnText_StartDialog, user.Language):
		user.ClearUserData()
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GptClearContext, user.Language), "")
		user.Path = "chatgpt/type/dialog"
	default:
		SendMessage(user, GetText(MsgText_UnknownCommand, user.Language), GetButton(btn_GptSpeechNewgen, user.Language), "")
	}

}
