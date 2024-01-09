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
	SendMessage(user, msgText, button_RemoveKeyboard, "HTML")

	msgText = `Выберите один из предложенных вариантов:`
	SendMessage(user, msgText, buttons_gptTypes, "")

	user.Path = "chatgpt/type"

}

// После выбора пользователем типа взаимодействия
func gpt_type(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	switch text {
	case "Начать диалог":
		msgText := `Запущен диалог с СhatGPT, чтобы очистить контекст от предыдущих сообщений - нажмите кнопку "Очистить контекст". Это позволяет сократить расход токенов.`
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		SendMessage(user, "Привет! Чем могу помочь?", buttons_gptClearContext, "")
		user.Path = "chatgpt/type/dialog"
	case "Сгенерировать аудио из текста":
		SendMessage(user, "Введите текст для аудио:", button_RemoveKeyboard, "")
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

	if text == "Очистить контекст" {
		user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
		SendMessage(user, "Контекст диалога очищен.", nil, "")
		SendMessage(user, "Привет! Чем могу помочь?", nil, "")
		return
	}

	if text == "Завершить диалог" {
		user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
		SendMessage(user, "Выберите один из предложенных вариантов:", buttons_gptTypes, "")
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
		msgText := "Напишите текст для озвучивания."
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		return
	}

	length := utf8.RuneCountInString(text)
	if length*20 > (Cfg.DailyLimitTokens - user.Tokens_used_gpt) {
		msgText := "Недостаточно токенов, укажите текст меньшей длины."
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		return
	}

	user.Options["text"] = text

	msgText := `Примеры звучания голосов`
	SendMessage(user, msgText, buttons_gptSampleSpeech, "")

	msgText = `Выберите голос для озвучивания текста:`
	SendMessage(user, msgText, buttons_gptVoices, "")

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

	voice, isErr := gptGetVoice(text)
	if isErr {
		msgText := "Выберите голос из предложенных вариантов."
		SendMessage(user, msgText, buttons_gptVoices, "")
		return
	}

	model := openai.TTSModel1HD
	inputText := user.Options["text"]

	Operation := SQL_NewOperation(user, "chatgpt", "speech", text)
	SQL_AddOperation(Operation)

	SendMessage(user, "Запущено создание аудиофайла ...", nil, "")

	// Отправка запроса в openai
	res, err := clientOpenAI.CreateSpeech(context.Background(),
		openai.CreateSpeechRequest{
			Model: model,
			Input: inputText,
			Voice: voice,
		})

	if err != nil {
		Logs <- NewLog(user, "chatgpt{speech1}", 1, err.Error())
		msgText := "Произошла непредвиденная ошибка. Попробуйте позже."
		SendMessage(user, msgText, buttons_gptTypes, "")
		user.Path = "chatgpt/type"
		return
	}

	// Сохранение в файл
	filename := fmt.Sprintf(WorkDir+"/data/speech_%d.mp3", user.ChatID)
	err = CreateFile(filename, res)
	if err != nil {
		Logs <- NewLog(user, "chatgpt{speech2}", 1, err.Error())
		msgText := "Произошла непредвиденная ошибка. Попробуйте позже."
		SendMessage(user, msgText, buttons_gptTypes, "")
		user.Path = "chatgpt/type"
		return
	}

	tokensForRequest := utf8.RuneCountInString(inputText) * 20
	user.Tokens_used_gpt = user.Tokens_used_gpt + tokensForRequest

	caption := fmt.Sprintf(`Результат генерации по тексту "%s", голос: "%s"`, inputText, voice)
	SendAudioMessage(user, filename, caption, buttons_gptSpeechNewgen)

	user.Path = "chatgpt/type/speech_text/voice/newgen"

}

// После получения результата генерации
func gpt_speech_newgen(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	switch text {
	case "Изменить текст":
		SendMessage(user, "Напишите текст для озвучивания:", button_RemoveKeyboard, "")
		user.Path = "chatgpt/type/speech_text"
	case "Выбрать другой голос":
		SendMessage(user, "Выберите голос для озвучивания текста:", buttons_gptVoices, "")
		user.Path = "chatgpt/type/speech_text/voice"
	case "Начать диалог":
		user.ClearUserData()
		SendMessage(user, "Привет! Чем могу помочь?", buttons_gptClearContext, "")
		user.Path = "chatgpt/type/dialog"
	default:
		SendMessage(user, "Неизвестная команда.", buttons_gptSpeechNewgen, "")
	}

}
