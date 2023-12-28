package main

import (
	"context"
	"log"
	"strings"

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

		if strings.Contains(errString, "This model's maximum context length is 4097 tokens") && firstLaunch { //чтобы в рекурсию не уйти

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

	//This model's maximum context length is 4097 tokens.

}
