package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
)

// /start - start msg
// /gemini - "введите вопрос"
//    text - result
// /chatgpt - "введите вопрос"
//    text - result
// /kandinsky - "введите запрос"
//    text - "выберите стиль изображения"
//	     style - result

var buttons_start = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Gemini"),
		tgbotapi.NewKeyboardButton("Kandinsky"),
		tgbotapi.NewKeyboardButton("ChatGPT"),
	),
)

type ResultOfRequest struct {
	Message     tgbotapi.Chattable
	Log_author  string
	Log_message string
}

func (r *ResultOfRequest) addUsernameIntoLog(username string) {
	r.Log_author = r.Log_author + "->" + username
}

func ProcessCommand(cmd string, upd tgbotapi.Update, user *UserInfo) ResultOfRequest {

	var result ResultOfRequest
	result.Log_author = "bot"

	Message := tgbotapi.NewMessage(user.ChatID, "")
	Message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	var msg_text string

	switch cmd {
	case "start":
		user.Model = ""
		msg_text = start(upd.Message.Chat.FirstName)
		Message.ParseMode = "HTML"
		Message.ReplyMarkup = buttons_start

	case "chatgpt":
		user.Model = "chatgpt"
		msg_text = "Напишите свой вопрос:"

	case "gemini":
		user.Model = "gemini"
		msg_text = "Напишите свой вопрос:"

	case "kandinsky":
		user.Model = "kandinsky"
		user.Stage = "text"
		msg_text = "Введите свой запрос:"

	case "clearcontext":
		user.History_Gemini = []*genai.Content{}
		msg_text = "История диалога с Gemini очищена."

	default:
		if slices.Contains(admins, user.Username) {
			switch cmd {
			case "stop":
				os.Exit(1)
			case "updconf":
				LoadConfig()
				msg_text = "Config updated."
			case "info":
				//msg_text = fmt.Sprintf("Gemini: %d\nChatGPT: %d\nKandinsky: %d",
				//	counter_gemini, counter_chatgpt, counter_kandinsky)
				msg_text = GetInfo()
			}
		}
	}

	Message.Text = msg_text
	result.Message = Message
	result.Log_message = msg_text
	if cmd == "start" {
		result.Log_message = "/start for " + user.Username
	}

	return result

}

func ProcessText(text string, user *UserInfo, upd tgbotapi.Update) ResultOfRequest {

	var result ResultOfRequest
	result.Log_author = user.Model

	Message := tgbotapi.NewMessage(user.ChatID, "")
	var msg_text string

	switch user.Model {
	case "chatgpt":
		Operation := SQL_NewOperation(user, text)
		SQL_AddOperation(Operation)

		msg_text = SendRequestToChatGPT(upd.Message.Text, user, true)

	case "gemini":
		Operation := SQL_NewOperation(user, text)
		SQL_AddOperation(Operation)

		msg_text = SendRequestToGemini(upd.Message.Text, user)

	case "kandinsky": // пользователь ввел текст по модели kand

		return ProcessInputText_Kandinsky(text, user)

	default:
		msg_text = "Не выбрана нейросеть для обработки запроса."
		Message.ReplyMarkup = buttons_start
		result.Log_author = "bot"
	}

	Message.Text = msg_text
	result.Message = Message
	result.Log_message = msg_text

	return result

}

func GetInfo() string {

	Now := time.Now().UTC().Add(3 * time.Hour)
	Yesterday := Now.AddDate(0, 0, -1)
	StartDay := time.Date(Now.Year(), Now.Month(), Now.Day(), 0, 0, 0, 0, Now.Location())

	result_24h, err1 := SQL_GetInfoOnDate(Yesterday)
	if err1 != "" {
		return err1
	}

	result_Today, err2 := SQL_GetInfoOnDate(StartDay)
	if err2 != "" {
		return err2
	}

	return fmt.Sprintf(
		`Last 24h
Users: %d
Gemini: %d | ChatGPT: %d | Kandinsky: %d

This day
Users: %d
Gemini: %d | ChatGPT: %d | Kandinsky: %d`,
		result_24h["users"], result_24h["gemini"], result_24h["chatgpt"], result_24h["kandinsky"],
		result_Today["users"], result_Today["gemini"], result_Today["chatgpt"], result_Today["kandinsky"])

}
