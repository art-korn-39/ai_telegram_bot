package main

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/sashabaranov/go-openai"
)

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
	if length*20 > (Get_TPD_gpt(user) - user.Tokens_used_gpt) {
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
	// ИЗМЕНИТЬ ТЕКСТ
	case GetText(BtnText_ChangeText, user.Language):
		SendMessage(user, GetText(MsgText_WriteTextForVoicing, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "chatgpt/type/speech_text"

	// ИЗМЕНИТЬ ГОЛОС
	case GetText(BtnText_ChooseAnotherVoice, user.Language):
		SendMessage(user, GetText(MsgText_SelectVoice, user.Language), GetButton(btn_GptVoices, ""), "")
		user.Path = "chatgpt/type/speech_text/voice"

	// НАЧАТЬ ДИАЛОГ
	case GetText(BtnText_StartDialog, user.Language):
		user.ClearUserData()
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GptClearContext, user.Language), "")
		user.Path = "chatgpt/type/dialog"

	default:
		SendMessage(user, GetText(MsgText_UnknownCommand, user.Language), GetButton(btn_GptSpeechNewgen, user.Language), "")
	}

}
