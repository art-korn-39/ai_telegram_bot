package main

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

const (
	OpenAIToken = "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"
)

var (
	clientOpenAI *openai.Client
)

func init() {

	clientOpenAI = openai.NewClient(OpenAIToken)

}

func SendRequestToChatGPT(text string, user *UserInfo) string {

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

	if err != nil {
		Logs <- Log{"ChatGPT", "request: " + text + "\nerror: " + err.Error(), true}
		return "Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже."
	}

	content := resp.Choices[0].Message.Content

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content},
	)

	user.Messages_ChatGPT = messages

	return content

}
