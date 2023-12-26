package main

import (
	"context"
	"fmt"
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

func TranslateInto(text string, language string) (result string) {

	request := strings.Join([]string{
		fmt.Sprintf("Translate text into %s", language),
		`"` + text + `"`,
	}, "/n")

	result = SendRequestToChatGPT(request)

	return result
}

func SendRequestToChatGPT(text string) string {

	<-delay_ChatGPT

	resp, err := clientOpenAI.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: text,
				},
			},
		},
	)

	if err != nil {
		Logs <- Log{"ChatGPT", "request: " + text + "\nerror: " + err.Error(), true}
		return "Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже."
	}

	return resp.Choices[0].Message.Content

}
