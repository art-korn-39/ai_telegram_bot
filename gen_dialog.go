package main

import (
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

// После ввода сообщения пользователем
func gen_dialog(user *UserInfo, text string) {

	if text == GetText(BtnText_EndDialog, user.Language) {
		user.Gen_History = []*genai.Content{}
		SendMessage(user, GetText(MsgText_SelectOption, user.Language), GetButton(btn_GenTypes, user.Language), "")
		user.Path = "gemini/type"
		return
	}

	if gen_DailyLimitOfRequestsIsOver(user, gen10) {
		return
	}

	<-delay_Gemini

	if Cfg.Gen_UseStream {
		gen_DialogSendMessageStream(user, text)
	} else {
		gen_DialogSendMessage(user, text)
	}

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

	// Получение результата из ответа
	result, err := gen_GetResultFromResponse(user, resp)
	if err != nil {
		Logs <- NewLog(user, "gemini", Error, err.Error())
		SendMessage(user, result, nil, "")
		return
	}

	gen_AddToHistory(user, text, result)

	SendMessage(user, result, nil, "")

	user.Usage.Gen10++
	Operation := SQL_NewOperation(user, "gemini", "dialog", text)

	SQL_AddOperation(Operation)

}

func gen_DialogSendMessageStream(user *UserInfo, text string) {

	//Ошибка: "напиши первые 3000 сиволов из библии"

	cs := gen_TextModel.StartChat()
	cs.History = user.Gen_History
	iter := cs.SendMessageStream(gen_ctx, genai.Text(text))

	var resp *genai.GenerateContentResponse
	var M tgbotapi.Message
	var resultFull string
	var err error
	withoutStream := false

	for {
		resp, err = iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

			// Пробуем ещё раз отправить запрос с другим контекстом
			var msgText string
			resp, msgText, err = gen_ProcessErrorsOfResponse(user, err, cs, text)

			// Завершаем процедуру если ошибка осталась
			if err != nil {
				Logs <- NewLog(user, "gemini", Error, err.Error())
				SendMessage(user, msgText, nil, "")
				return
				// Если удалось получить ответ, то выходим из цикла, т.к. результат уже есть
			} else {
				withoutStream = true
				break
			}

		}

		// Получение результата из ответа
		result, err := gen_GetResultFromResponse(user, resp)
		if err != nil {
			Logs <- NewLog(user, "gemini", Error, err.Error())
			SendMessage(user, result, nil, "")
			return
		}

		resultFull = resultFull + result

		if (M == tgbotapi.Message{}) {
			M = SendMessage(user, resultFull, nil, "")
		} else {
			SendEditMessage(user, M.MessageID, resultFull)
		}

		// чтобы избежать ошибки "Content = nil"
		time.Sleep(1 * time.Second)

	}

	// Было прерывание потока
	if withoutStream {
		result, err := gen_GetResultFromResponse(user, resp)
		if err != nil {
			Logs <- NewLog(user, "gemini", Error, err.Error())
			SendMessage(user, result, nil, "")
			return
		}

		resultFull = result

		if (M == tgbotapi.Message{}) {
			M = SendMessage(user, resultFull, nil, "")
		} else {
			SendEditMessage(user, M.MessageID, resultFull)
		}
	}

	gen_AddToHistory(user, text, resultFull)

	user.Usage.Gen10++
	Operation := SQL_NewOperation(user, "gemini", "dialog", text)
	SQL_AddOperation(Operation)

}

func gen_AddToHistory(user *UserInfo, t1, t2 string) {

	history := append(user.Gen_History,
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(t1),
			},
			Role: "user",
		},
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(t2),
			},
			Role: "model",
		},
	)

	user.Gen_History = history

}

// Обработку ошибок выполняем тут же
func gen_GetResultFromResponse(user *UserInfo, resp *genai.GenerateContentResponse) (string, error) {

	if resp.Candidates[0].Content == nil {

		user.Gen_History = []*genai.Content{}
		msgText := GetText(MsgText_BadRequest2, user.Language)
		err := errors.New("resp.Candidates[0].Content = nil")
		return msgText, err

	}

	result := string(resp.Candidates[0].Content.Parts[0].(genai.Text))

	return result, nil

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

	case "googleapi: Error 400:":
		time.Sleep(time.Millisecond * 200)
		resp, err = cs.SendMessage(gen_ctx, genai.Text(text))
		if err != nil {
			msgText = GetText(MsgText_GenGeoError, user.Language)
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
