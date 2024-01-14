package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
)

func gpt_WelcomeTextMessage(u *UserInfo) string {

	duration := GetDurationToNextDay()
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - hours*60

	return fmt.Sprintf(GetText(MsgText_ChatGPTHello, u.Language),
		max(Cfg.TPD_gpt-u.Tokens_used_gpt, 0),
		hours,
		minutes)
}

func gpt_DialogSendMessage(user *UserInfo, text string, firstLaunch bool) {

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

			SendMessage(user, GetText(MsgText_LimitOf4097TokensReached, user.Language), nil, "")

			Logs <- NewLog(user, "chatgpt", Warning, "request: "+text+"\nwarning: "+errString)

			user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
			gpt_DialogSendMessage(user, text, false)
			return

		} else {
			Logs <- NewLog(user, "chatgpt", Error, "request: "+text+"\nerror: "+errString)
			SendMessage(user, GetText(MsgText_ErrorWhileProcessingRequest, user.Language), nil, "")
			return
		}
	}

	user.Tokens_used_gpt = user.Tokens_used_gpt + resp.Usage.TotalTokens

	content = resp.Choices[0].Message.Content

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content},
	)

	user.Messages_ChatGPT = messages

	SendMessage(user, content, nil, "")

}

func gpt_SendSpeechSamples(user *UserInfo) {

	for _, voice := range voices {
		filename := WorkDir + "/samples/gpt_voice_" + voice + ".mp3"
		Message := tgbotapi.NewAudioUpload(user.ChatID, filename)
		Message.Caption = voice
		Bot.Send(Message)
	}

}

func gpt_DailyLimitOfTokensIsOver(u *UserInfo) bool {

	if u.Tokens_used_gpt >= Cfg.TPD_gpt {
		duration := GetDurationToNextDay()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) - hours*60
		msgText := fmt.Sprintf(GetText(MsgText_DailyTokenLimitExceeded, u.Language), hours, minutes)
		SendMessage(u, msgText, GetButton(btn_Models, u.Language), "")
		u.Path = "start"
		return true
	}

	return false

}

func gpt_GetSpeechVoice(voice string) (v openai.SpeechVoice, isError bool) {

	array := []openai.SpeechVoice{openai.VoiceAlloy, openai.VoiceEcho, openai.VoiceFable,
		openai.VoiceOnyx, openai.VoiceNova, openai.VoiceShimmer}

	SV := openai.SpeechVoice(strings.ToLower(voice))

	i := slices.Index(array, SV)
	if i == -1 {
		return "", true
	} else {
		return array[i], false
	}

}
