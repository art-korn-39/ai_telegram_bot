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

	return fmt.Sprintf(`Вас приветствует ChatGPT 3.5 Turbo 🤖

Текущий остаток токенов: <b>%d</b> <i>(обновится через: %d ч. %d мин.)</i>`,
		max(Cfg.DailyLimitTokens-u.Tokens_used_gpt, 0),
		hours,
		minutes)

	// return fmt.Sprintf(`Вас приветствует ChatGPT 3.5 Turbo 🤖

	// *Каждые сутки вам начисляется <b>%d</b> токенов для использования.
	// - 1 токен равен примерно 4 символам английского алфавита или 2 символам русского алфавита.
	// - При создании аудио из текста расходуется приблизительно 1000 токенов на каждые 50 символов.

	// Текущий остаток токенов: <b>%d</b>`, Cfg.DailyLimitTokens, max(Cfg.DailyLimitTokens-u.Tokens_used_gpt, 0))

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

			SendMessage(user, "Достингут лимит в 4097 токенов, контекст диалога очищен.", nil, "")

			Logs <- NewLog(user, "chatgpt", Warning, "request: "+text+"\nwarning: "+errString)

			user.Messages_ChatGPT = []openai.ChatCompletionMessage{}
			gpt_DialogSendMessage(user, text, false)
			return

		} else {
			Logs <- NewLog(user, "chatgpt", Error, "request: "+text+"\nerror: "+errString)
			SendMessage(user, "Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже.", nil, "")
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

	if u.Tokens_used_gpt >= Cfg.DailyLimitTokens {
		duration := GetDurationToNextDay()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) - hours*60
		msgText := fmt.Sprintf("Превышен дневной лимит токенов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.", hours, minutes)
		SendMessage(u, msgText, buttons_gptTypes, "")
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
