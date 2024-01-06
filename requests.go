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
	Message         tgbotapi.Chattable
	Log_author      string
	Log_message     string
	UserInfoChanged bool
}

func (r *ResultOfRequest) addUsernameIntoLog(username string) {
	r.Log_author = r.Log_author + "->" + username
}

func ProcessCommand(cmd string, upd tgbotapi.Update, user *UserInfo) ResultOfRequest {

	var result ResultOfRequest
	result.Log_author = "bot"
	result.UserInfoChanged = true

	Message := tgbotapi.NewMessage(user.ChatID, "")
	Message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	var msg_text string

	switch cmd {
	case "start":
		//		user.Model = ""
		msg_text = start(upd.Message.Chat.FirstName)
		Message.ParseMode = "HTML"
		Message.ReplyMarkup = buttons_start

	case "chatgpt":
		//		user.Model = "chatgpt"
		msg_text = "Напишите свой вопрос:"

	case "gemini":
		//		user.Model = "gemini"
		msg_text = "Напишите свой вопрос:"

	case "kandinsky":
		//		user.Model = "kandinsky"
		//		user.Stage = "text"
		msg_text = "Введите свой запрос:"

	case "clearcontext":
		user.History_Gemini = []*genai.Content{}
		msg_text = "История диалога с Gemini очищена."

	case "stop":
		if user.Username == "Art_Korn_39" {
			os.Exit(1)
		}
	default:
		result.UserInfoChanged = false
		if slices.Contains(admins, user.Username) {
			switch cmd {
			case "updconf":
				LoadConfig()
				msg_text = "Config updated."
			case "info":
				msg_text = GetInfo()
			}
		}
	}

	Message.Text = msg_text
	result.Message = Message
	result.Log_message = msg_text
	switch cmd {
	case "start":
		result.Log_message = "/start for " + user.Username
	case "info":
		result.Log_message = "/info for " + user.Username
	}

	return result

}

func ProcessText(text string, user *UserInfo, upd tgbotapi.Update) ResultOfRequest {

	var result ResultOfRequest
	//	result.Log_author = user.Model
	result.UserInfoChanged = true

	Message := tgbotapi.NewMessage(user.ChatID, "")
	var msg_text string

	//	switch user.Model {
	switch "123" {
	case "chatgpt":
		//Operation := SQL_NewOperation(user, text)
		//SQL_AddOperation(Operation)

		msg_text = SendRequestToChatGPT(upd.Message.Text, user, true)

	case "gemini":
		//Operation := SQL_NewOperation(user, text)
		//SQL_AddOperation(Operation)

		msg_text = SendRequestToGemini(upd.Message.Text, user)

	case "kandinsky": // пользователь ввел текст по модели kand

		return ProcessInputText_Kandinsky(text, user)

	default:
		msg_text = "Не выбрана нейросеть для обработки запроса."
		Message.ReplyMarkup = buttons_start
		result.Log_author = "bot"
		result.UserInfoChanged = false
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
	December30 := time.Date(2023, 12, 30, 0, 0, 0, 0, time.Local)
	December25 := time.Date(2023, 12, 25, 0, 0, 0, 0, time.Local)

	result_dec25, err0 := SQL_GetInfoOnDate(December25)
	if err0 != "" {
		return err0
	}

	result_dec30, err1 := SQL_GetInfoOnDate(December30)
	if err1 != "" {
		return err1
	}

	result_24h, err2 := SQL_GetInfoOnDate(Yesterday)
	if err2 != "" {
		return err2
	}

	result_Today, err3 := SQL_GetInfoOnDate(StartDay)
	if err3 != "" {
		return err3
	}

	return fmt.Sprintf(
		`All time
Gemini: %d | ChatGPT: %d | Kandinsky: %d

From 30.12.23
Users: %d
Gemini: %d | ChatGPT: %d | Kandinsky: %d

Last 24h
Users: %d
Gemini: %d | ChatGPT: %d | Kandinsky: %d

Today
Users: %d
Gemini: %d | ChatGPT: %d | Kandinsky: %d`,
		result_dec25["gemini"], result_dec25["chatgpt"], result_dec25["kandinsky"],
		result_dec30["users"], result_dec30["gemini"], result_dec30["chatgpt"], result_dec30["kandinsky"],
		result_24h["users"], result_24h["gemini"], result_24h["chatgpt"], result_24h["kandinsky"],
		result_Today["users"], result_Today["gemini"], result_Today["chatgpt"], result_Today["kandinsky"])

}
