package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

// После ввода сообщения пользователем
func gen15_dialog(user *UserInfo, text string) {

	if text == GetText(BtnText_EndDialog, user.Language) {
		user.Gen_History = []*genai.Content{}
		SendMessage(user, GetText(MsgText_SelectOption, user.Language), GetButton(btn_GenTypes, user.Language), "")
		user.Path = "gen15/type"
		return
	}

	if gen_DailyLimitOfRequestsIsOver(user, gen15) {
		return
	}

	<-delay_Gen15

	if Cfg.Gen_UseStream {
		gen15_DialogSendMessageStream(user, text)
	} else {
		gen15_DialogSendMessage(user, text)
	}

}

func gen15_DialogSendMessage(user *UserInfo, text string) {

	var msgText string
	cs := gen15_Model.StartChat()
	cs.History = user.Gen_History

	resp, err := cs.SendMessage(gen_ctx, genai.Text(text))
	if err != nil {

		resp, msgText, err = gen_ProcessErrorsOfResponse(user, err, cs, text)

		// Отправляем сообщение и завершаем процедуру если ошибка осталась
		if err != nil {
			Logs <- NewLog(user, "gen15", Error, err.Error())
			SendMessage(user, msgText, nil, "")
			return
		}
	}

	// Получение результата из ответа
	result, err := gen_GetResultFromResponse(user, resp)
	if err != nil {
		Logs <- NewLog(user, "gen15", Error, err.Error())
		SendMessage(user, result, nil, "")
		return
	}

	gen_AddToHistory(user, text, result)

	SendMessage(user, result, nil, "")

	user.Usage.Gen10++
	Operation := SQL_NewOperation(user, "gen15", "dialog", "", text)

	SQL_AddOperation(Operation)

}

func gen15_DialogSendMessageStream(user *UserInfo, text string) {

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
				Logs <- NewLog(user, "gen15", Error, err.Error())
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
			Logs <- NewLog(user, "gen15", Error, err.Error())
			SendMessage(user, result, nil, "")
			return
		}

		resultFull = resultFull + result

		if (M == tgbotapi.Message{}) {
			M = SendMessage(user, resultFull, nil, "")
		} else {
			SendEditMessage(user, M.MessageID, resultFull)
		}

	}

	// Было прерывание потока
	if withoutStream {
		result, err := gen_GetResultFromResponse(user, resp)
		if err != nil {
			Logs <- NewLog(user, "gen15", Error, err.Error())
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
	Operation := SQL_NewOperation(user, "gen15", "dialog", "", text)
	SQL_AddOperation(Operation)

}
