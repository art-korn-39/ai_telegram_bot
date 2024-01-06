package main

import (
	"context"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
)

var (
	button_gptClearContext = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Очистить контекст"),
		),
	)
)

// После команды "/chatgpt" или при вводе текста = "chatgpt"
func gpt_start(user *UserInfo) {

	msgText := `Запущен диалог с СhatGPT, чтобы очистить контекст от предыдущих сообщений - нажмите кнопку "Очистить контекст" или это произойдёт автоматически, когда будет превышен лимит токенов.`
	SendMessage(user, msgText, button_gptClearContext, "")

	msgText = `Привет! Чем могу помочь?`
	SendMessage(user, msgText, button_gptClearContext, "")

	user.Path = "chatgpt/dialog"

}

// После ввода сообщения пользователем
func gpt_dialog(user *UserInfo, text string, firstLaunch bool) {

	if text == "Очистить контекст" {
		user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
		SendMessage(user, "Контекст диалога очищен.", nil, "")
		SendMessage(user, "Привет! Чем могу помочь?", nil, "")
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

			Logs <- NewLog(user, "ChatGPT", Warning, "request: "+text+"\nwarning: "+errString)

			user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
			gpt_dialog(user, text, false)
			return

		} else {
			Logs <- NewLog(user, "ChatGPT", Error, "request: "+text+"\nerror: "+errString)
			SendMessage(user, "Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже.", nil, "")
			return
		}
	}

	content = resp.Choices[0].Message.Content

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content},
	)

	user.Messages_ChatGPT = messages

	SendMessage(user, content, button_gptClearContext, "")

}
