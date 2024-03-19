package main

import (
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
		max(Get_TPD_gpt(u)-u.Usage.GPT, 0),
		hours,
		minutes)
}

func gpt_DailyLimitOfTokensIsOver(u *UserInfo) bool {

	if slices.Contains(Cfg.WhiteList, u.Username) {
		return false
	}

	if u.Usage.GPT >= Get_TPD_gpt(u) {
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

func gpt_SendSpeechSamples(user *UserInfo) {

	for _, voice := range voices {
		filename := WorkDir + "/samples/gpt_voice_" + voice + ".mp3"
		Message := tgbotapi.NewAudioUpload(user.ChatID, filename)
		Message.Caption = voice
		Bot.Send(Message)
	}

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
