package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
)

// После ввода сообщения пользователем
func gen_dialog(user *UserInfo, text string) {

	if text == GetText(BtnText_EndDialog, user.Language) {
		user.Gen_History = []*genai.Content{}
		SendMessage(user, GetText(MsgText_SelectOption, user.Language), GetButton(btn_GenTypes, user.Language), "")
		user.Path = "gemini/type"
		return
	}

	if gen_DailyLimitOfRequestsIsOver(user) {
		return
	}

	<-delay_Gemini

	user.Requests_today_gen++

	Operation := SQL_NewOperation(user, "gemini", "dialog", text)
	SQL_AddOperation(Operation)

	gen_DialogSendMessage(user, text)

}

func gen_DialogSendMessage(user *UserInfo, text string) {

	var msgText string
	cs := gen_TextModel.StartChat()
	cs.History = user.Gen_History

	resp, err := cs.SendMessage(gen_ctx, genai.Text(text))
	if err != nil {

		resp, msgText, err = gen_ProcessErrorsOfResponse(user, err, cs, text)

		// Отправляем сообщение и завершаем процедуру если ошибка осталась
		if err != nil {
			Logs <- NewLog(user, "gemini", Error, err.Error())
			SendMessage(user, msgText, nil, "")
			return
		}
	}

	if resp.Candidates[0].Content == nil {
		user.Gen_History = []*genai.Content{}
		Logs <- NewLog(user, "gemini", Error, "resp.Candidates[0].Content = nil")
		msgText = GetText(MsgText_BadRequest2, user.Language)
		SendMessage(user, msgText, nil, "")
		return
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	history := append(user.Gen_History,
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

	user.Gen_History = history

	msgText = string(result)
	SendMessage(user, msgText, nil, "")

}

func gen_ProcessErrorsOfResponse(user *UserInfo, errIn error, cs *genai.ChatSession, text string) (resp *genai.GenerateContentResponse, msgText string, err error) {

	err = errIn
	errString := errIn.Error()

	switch errString {
	case "blocked: candidate: FinishReasonSafety":
		NewConnectionGemini()
		cs = gen_TextModel.StartChat()
		user.Gen_History = []*genai.Content{}
		resp, err = cs.SendMessage(gen_ctx, genai.Text(text))
		if err != nil {
			msgText = GetText(MsgText_BadRequest3, user.Language)
		}

	case "googleapi: Error 500:":
		time.Sleep(time.Millisecond * 200)
		resp, err = cs.SendMessage(gen_ctx, genai.Text(text))
		if err != nil {
			msgText = GetText(MsgText_UnexpectedError, user.Language)
		}

	case "blocked: prompt: BlockReasonSafety":
		msgText = GetText(MsgText_BadRequest4, user.Language)
	default:
		msgText = GetText(MsgText_UnexpectedError, user.Language)
	}

	return
}

func gen_TranslateToEnglish(text string) (string, error) {

	if IsEngByLoop(text) {
		return text, nil
	}

	prompt := fmt.Sprintf(`translate to english next text:
%s`, text)

	chatSession := gen_TextModelWithCensor.StartChat()
	resp, err := chatSession.SendMessage(gen_ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if resp.Candidates[0].Content == nil {
		err := errors.New("{gemini} resp.Candidates[0].Content = nil")
		return "", err
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	return string(result), nil

}
