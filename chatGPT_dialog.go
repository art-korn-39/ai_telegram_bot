package main

import (
	"context"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// После ввода сообщения пользователем
func gpt_dialog(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	if text == GetText(BtnText_ClearContext, user.Language) {
		user.Gpt_History = []openai.ChatCompletionMessage{}
		SendMessage(user, GetText(MsgText_DialogContextCleared, user.Language), nil, "")
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), nil, "")
		return
	}

	if text == GetText(BtnText_EndDialog, user.Language) {
		user.Gpt_History = []openai.ChatCompletionMessage{}
		SendMessage(user, GetText(MsgText_SelectOption, user.Language), GetButton(btn_GptTypes, user.Language), "")
		user.Path = "chatgpt/type"
		return
	}

	<-delay_ChatGPT

	Operation := SQL_NewOperation(user, "chatgpt", "dialog", text)
	SQL_AddOperation(Operation)

	gpt_DialogSendMessage(user, text, true)

}

func gpt_DialogSendMessage(user *UserInfo, text string, firstLaunch bool) {

	messages := append(user.Gpt_History, openai.ChatCompletionMessage{
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

			SendMessage(user, GetText(MsgText_LimitOf4097TokensReached, user.Language), nil, "")

			Logs <- NewLog(user, "chatgpt", Warning, "request: "+text+"\nwarning: "+errString)

			user.Gpt_History = []openai.ChatCompletionMessage{}
			gpt_DialogSendMessage(user, text, false)
			return

		} else {
			Logs <- NewLog(user, "chatgpt", Error, "request: "+text+"\nerror: "+errString)
			SendMessage(user, GetText(MsgText_ErrorWhileProcessingRequest, user.Language), nil, "")
			return
		}
	}

	user.Tokens_used_gpt = user.Tokens_used_gpt + resp.Usage.TotalTokens
	user.Usage.GPT = user.Usage.GPT + resp.Usage.TotalTokens

	content = resp.Choices[0].Message.Content

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content},
	)

	user.Gpt_History = messages

	SendMessage(user, content, nil, "")

}
