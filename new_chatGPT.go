package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
)

// OpenAIToken = "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" //Коли
// OpenAIToken = "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" //Мой

//для оценки токенов
//2 символа RU = 1 токен
//4 символа EN = 1 токен
//utf8.RuneCountInString(str)

var (
	clientOpenAI     *openai.Client
	voices           = []string{"fable", "nova", "echo", "onyx"}
	buttons_gptTypes = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Начать диалог"),
			tgbotapi.NewKeyboardButton("Сгенерировать аудио из текста"),
		),
	)
	buttons_gptVoices = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("onyx"),
			tgbotapi.NewKeyboardButton("nova"),
			tgbotapi.NewKeyboardButton("echo"),
			tgbotapi.NewKeyboardButton("fable"),
		),
	)
	buttons_gptClearContext = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Очистить контекст"),
			tgbotapi.NewKeyboardButton("Завершить диалог"),
		),
	)
	buttons_gptAudioNewgen = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить текст"),
			tgbotapi.NewKeyboardButton("Выбрать другой голос"),
			tgbotapi.NewKeyboardButton("Начать диалог"),
		),
	)
	buttons_sampleAudio = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Audio samples", "gpt_audio_samples"),
		),
	)
)

// После команды "/chatgpt" или при вводе текста = "chatgpt"
func gpt_start(user *UserInfo) {

	msgText := fmt.Sprintf(`Вас приветствует ChatGPT 3.5 Turbo 🤖

*Каждые сутки вам начисляется <b>%d</b> токенов для использования.
- 1 токен равен примерно 4 символам английского алфавита или 2 символам русского алфавита.
- При создании аудио из текста расходуется приблизительно 1000 токенов на каждые 50 символов.

Текущий остаток токенов: <b>%d</b>`, Cfg.DailyLimitTokens, max(Cfg.DailyLimitTokens-user.TokensUsed_ChatGPT, 0))

	SendMessage(user, msgText, button_RemoveKeyboard, "HTML")

	msgText = `Выберите один из предложенных вариантов:`
	SendMessage(user, msgText, buttons_gptTypes, "")

	user.Path = "chatgpt/type"

}

// После выбора пользователем типа взаимодействия
func gpt_type(user *UserInfo, text string) {

	if DailyLimitOfTokensIsOver(user) {
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
		user.Path = "chatgpt/type/audio_text"
	default:
		gpt_dialog(user, text, true)
		user.Path = "chatgpt/type/dialog"
	}

}

// После ввода сообщения пользователем
func gpt_dialog(user *UserInfo, text string, firstLaunch bool) {

	if DailyLimitOfTokensIsOver(user) {
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

	messages := append(user.Messages_ChatGPT, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})

	resp, err := clientOpenAI.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	var content string
	if err != nil {
		errString := err.Error()

		// превышен лимит токенов, очищаем сообщения и отправляем запрос ещё раз
		if strings.Contains(errString, "This model's maximum context length is 4097 tokens") && firstLaunch { //чтобы в рекурсию не уйти

			SendMessage(user, "Достингут лимит в 4097 токенов, контекст диалога очищен.", nil, "")

			Logs <- NewLog(user, "chatgpt", Warning, "request: "+text+"\nwarning: "+errString)

			user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
			gpt_dialog(user, text, false)
			return

		} else {
			Logs <- NewLog(user, "chatgpt", Error, "request: "+text+"\nerror: "+errString)
			SendMessage(user, "Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже.", nil, "")
			return
		}
	}

	user.TokensUsed_ChatGPT = user.TokensUsed_ChatGPT + resp.Usage.TotalTokens

	content = resp.Choices[0].Message.Content

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content},
	)

	user.Messages_ChatGPT = messages

	//SendMessage(user, content, buttons_gptClearContext, "")
	SendMessage(user, content, nil, "")

}

// После ввода текста пользователем
func gpt_audio_text(user *UserInfo, text string) {

	if DailyLimitOfTokensIsOver(user) {
		return
	}

	// Проверяем наличие текста в сообщении
	if text == "" {
		msgText := "Напишите текст для озвучивания."
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		return
	}

	length := utf8.RuneCountInString(text)
	if length*20 > (Cfg.DailyLimitTokens - user.TokensUsed_ChatGPT) {
		msgText := "Недостаточно токенов, укажите текст меньшей длины."
		SendMessage(user, msgText, button_RemoveKeyboard, "")
		return
	}

	user.Options["text"] = text

	msgText := `Примеры звучания голосов`
	SendMessage(user, msgText, buttons_sampleAudio, "")

	msgText = `Выберите голос для озвучивания текста:`
	SendMessage(user, msgText, buttons_gptVoices, "")

	user.Path = "chatgpt/type/audio_text/voice"

}

// После выбора голоса
func gpt_speech(user *UserInfo, text string) {

	if DailyLimitOfTokensIsOver(user) {
		return
	}

	if text == "gpt_audio_samples" {
		gpt_SendAudioSamples(user)
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

	res, err := clientOpenAI.CreateSpeech(context.Background(),
		openai.CreateSpeechRequest{
			Model: model,
			Input: inputText,
			Voice: voice,
		})

	if err != nil {
		Logs <- NewLog(user, "chatgpt{audio}", 1, err.Error())
		msgText := "Произошла непредвиденная ошибка. Попробуйте позже."
		SendMessage(user, msgText, buttons_gptTypes, "")
		user.Path = "chatgpt/type"
		return
	}

	user.TokensUsed_ChatGPT = user.TokensUsed_ChatGPT + utf8.RuneCountInString(inputText)*20

	filename := fmt.Sprintf(WorkDir+"/data/audio_%d.mp3", user.ChatID)
	outFile, _ := os.Create(filename)
	defer outFile.Close()
	_, err = io.Copy(outFile, res)
	if err != nil {
		Logs <- NewLog(user, "chatgpt{audio}", 1, err.Error())
		msgText := "Произошла непредвиденная ошибка. Попробуйте позже."
		SendMessage(user, msgText, buttons_gptTypes, "")
		user.Path = "chatgpt/type"
		return
	}

	Message := tgbotapi.NewAudioUpload(user.ChatID, filename)
	Message.Caption = fmt.Sprintf(`Результат генерации по тексту "%s", голос: "%s"`, inputText, voice)
	Message.ReplyMarkup = buttons_gptAudioNewgen
	Bot.Send(Message)

	Logs <- NewLog(user, "chatgpt", Info, filename)

	user.Path = "chatgpt/type/audio_text/voice/newgen"

}

// После получения результата генерации
func gpt_audio_newgen(user *UserInfo, text string) {

	if DailyLimitOfTokensIsOver(user) {
		return
	}

	switch text {
	case "Изменить текст":
		SendMessage(user, "Напишите текст для озвучивания:", button_RemoveKeyboard, "")
		user.Path = "chatgpt/type/audio_text"
	case "Выбрать другой голос":
		SendMessage(user, "Выберите голос для озвучивания текста:", buttons_gptVoices, "")
		user.Path = "chatgpt/type/audio_text/voice"
	case "Начать диалог":
		user.ClearUserData()
		SendMessage(user, "Привет! Чем могу помочь?", buttons_gptClearContext, "")
		user.Path = "chatgpt/type/dialog"
	default:
		SendMessage(user, "Неизвестная команда.", buttons_gptAudioNewgen, "")
	}

}

func gpt_SendAudioSamples(user *UserInfo) {

	for _, voice := range voices {
		filename := WorkDir + "/samples/gpt_voice_" + voice + ".mp3"
		Message := tgbotapi.NewAudioUpload(user.ChatID, filename)
		Message.Caption = voice
		Bot.Send(Message)
	}

}

func DailyLimitOfTokensIsOver(u *UserInfo) bool {

	if u.TokensUsed_ChatGPT >= Cfg.DailyLimitTokens {
		duration := GetDurationToNextDay()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) - hours*60
		msgText := fmt.Sprintf("Превышен дневной лимит токенов, дождитесь обновления лимита (%d ч %d мин) или воспользуйтесь другой нейросетью.", hours, minutes)
		SendMessage(u, msgText, buttons_gptTypes, "")
		return true
	}

	return false

}
