package main

import (
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
)

func gen_DialogSendMessage(user *UserInfo, text string) {

	var msgText string
	cs := model_Gemini.StartChat()
	cs.History = user.Messages_Gemini

	resp, err := cs.SendMessage(ctx_Gemini, genai.Text(text))
	if err != nil {
		errorString := err.Error()
		Logs <- NewLog(user, "gemini", Error, errorString)

		if errorString == "blocked: candidate: FinishReasonSafety" {

			// В случае данного вида ошибки - запускаем новый клиент соединения
			NewConnectionGemini()

			// Очищаем контекст
			cs.History = []*genai.Content{}
			user.Messages_Gemini = []*genai.Content{}

			// Отправляем повторно
			Logs <- NewLog(user, "gemini", Error, "Повторная отправка запроса ...")
			resp, err = cs.SendMessage(ctx_Gemini, genai.Text(text))
			if err != nil {
				Logs <- NewLog(user, "gemini", Error, err.Error())
				msgText = GetText(MsgText_BadRequest3, user.Language)
			} else {
				Logs <- NewLog(user, "gemini", Error, "После очистки контекста - запрос ушел.")
			}

		} else if errorString == "blocked: prompt: BlockReasonSafety" {

			msgText = GetText(MsgText_BadRequest4, user.Language)

		} else if errorString == "googleapi: Error 500:" {

			// Отправляем сообщение повторно
			time.Sleep(time.Millisecond * 200)
			Logs <- NewLog(user, "gemini", Error, "Повторная отправка запроса ...")
			resp, err = cs.SendMessage(ctx_Gemini, genai.Text(text))
			if err != nil {
				Logs <- NewLog(user, "gemini", Error, err.Error())
				msgText = GetText(MsgText_UnexpectedError, user.Language)
			}

		} else {
			msgText = GetText(MsgText_UnexpectedError, user.Language)
		}

		// Отправляем сообщение и завершаем процедуру если получили ошибку в ответ
		if err != nil {
			SendMessage(user, msgText, nil, "")
			return
		}
	}

	if resp.Candidates[0].Content == nil {
		Logs <- NewLog(user, "gemini", Error, "resp.Candidates[0].Content = nil")
		msgText = GetText(MsgText_BadRequest2, user.Language)
		SendMessage(user, msgText, nil, "")
		return
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	history := append(user.Messages_Gemini,
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(text),
			},
			Role: "user",
		},
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(result),
			},
			Role: "model",
		},
	)

	user.Messages_Gemini = history

	msgText = string(result)
	SendMessage(user, msgText, nil, "")

}

func gen_DailyLimitOfRequestsIsOver(u *UserInfo) bool {

	if u.Requests_today_gen >= Cfg.RPD_gen {
		duration := GetDurationToNextDay()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) - hours*60
		msgText := fmt.Sprintf(GetText(MsgText_DailyRequestLimitExceeded, u.Language), hours, minutes)
		SendMessage(u, msgText, GetButton(btn_Models, u.Language), "")
		u.Path = "start"
		return true
	}

	return false

}
