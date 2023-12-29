package main

import (
	"fmt"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
)

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

	Message := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	Message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	var msg_text string

	switch cmd {
	case "start":
		user.Model = ""
		msg_text = start(upd.Message.Chat.FirstName)
		//msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		Message.ParseMode = "HTML"
		Message.ReplyMarkup = buttons_start
		//result.Message = msg
		//result.Log_message = "/start for " + upd.Message.Chat.UserName

	case "chatgpt":
		user.Model = "chatgpt"
		msg_text = "Напишите свой вопрос:"
		// msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		// msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		// result.Message = msg
		// result.Log_message = msg_text

	case "gemini":
		user.Model = "gemini"
		msg_text = "Напишите свой вопрос:"
		// msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		// msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		// result.Message = msg
		// result.Log_message = msg_text

	case "kandinsky":
		user.Model = "kandinsky"
		user.Stage = "text"
		msg_text = "Введите свой запрос:"
		// msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		// msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		// result.Message = msg
		// result.Log_message = msg_text

	case "clearcontext":
		user.History_Gemini = []*genai.Content{}
		msg_text = "История диалога с Gemini очищена."
		// msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		// msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		// result.Message = msg
		// result.Log_message = msg_text

	default:
		if upd.Message.From.UserName == "Art_Korn_39" {
			switch cmd {
			case "stop":
				os.Exit(1)
			case "updconf":
				LoadConfig()
				msg_text = "Config updated."
				//result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
			case "info":
				msg_text = fmt.Sprintf("Gemini: %d\nChatGPT: %d\nKandinsky: %d",
					counter_gemini, counter_chatgpt, counter_kandinsky)
				//result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
			}
		}
	}

	Message.Text = msg_text
	result.Message = Message
	result.Log_message = msg_text
	if cmd == "start" {
		result.Log_message = "/start for " + upd.Message.Chat.UserName
	}

	return result

}

func ProcessText(text string, user *UserInfo, upd tgbotapi.Update) ResultOfRequest {

	var result ResultOfRequest
	result.Log_author = user.Model

	Message := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	var msg_text string

	switch user.Model {
	case "chatgpt":
		Operation := NewSQLOperation(user, upd, text)
		SQL_AddOperation(Operation)
		counter_chatgpt++

		msg_text = SendRequestToChatGPT(upd.Message.Text, user, true)
		//result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		//result.Log_author = "ChatGPT"
		//result.Log_message = msg_text

	case "gemini":
		Operation := NewSQLOperation(user, upd, text)
		SQL_AddOperation(Operation)
		counter_gemini++

		msg_text = SendRequestToGemini(upd.Message.Text, user)
		//result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		//result.Log_author = "Gemini"
		//result.Log_message = msg_text

	case "kandinsky": // пользователь ввел текст по модели kand

		return ProcessInputText_Kandinsky(text, user, upd)

	default:
		msg_text = "Не выбрана нейросеть для обработки запроса."
		//result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "bot"
		//result.Log_message = msg_text
	}

	Message.Text = msg_text
	result.Message = Message
	result.Log_message = msg_text

	return result

}
