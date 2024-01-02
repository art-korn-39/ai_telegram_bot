package main

import (
	"context"
	"log"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
)

const (
	OpenAIToken = "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" //Коли
	//OpenAIToken = "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" // мой
)

var (
	clientOpenAI *openai.Client
)

func init() {

	clientOpenAI = openai.NewClient(OpenAIToken)

}

//для оценки токенов
//2 символа RU = 1 токен
//4 символа EN = 1 токен
//utf8.RuneCountInString(str)

func SendRequestToChatGPT(text string, user *UserInfo, firstLaunch bool) string {

	<-delay_ChatGPT

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

			Message := tgbotapi.NewMessage(user.ChatID, "Достингут лимит в 4097 токенов, контекст диалога очищен.")
			Bot.Send(Message)

			log.Println("limit 4097 tokens, refresh context")
			user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
			content = SendRequestToChatGPT(text, user, false)

		} else {
			Logs <- Log{"ChatGPT", "request: " + text + "\nerror: " + errString, true}
			return "Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже."
		}

	} else {
		content = resp.Choices[0].Message.Content
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content},
	)

	user.Messages_ChatGPT = messages

	return content

}
