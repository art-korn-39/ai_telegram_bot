package main

import (
	"fmt"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
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

func ProcessCommand(cmd string, upd tgbotapi.Update, user *UserInfo) ResultOfRequest {

	var result ResultOfRequest
	result.Log_author = "bot"

	switch cmd {
	case "start":
		user.Model = ""
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, start(upd.Message.Chat.FirstName))
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = buttons_start

		result.Message = msg
		result.Log_message = "/start for " + upd.Message.Chat.UserName
	case "chatgpt":
		user.Model = "chatgpt"
		msg_text := "Напишите свой вопрос:"
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		result.Message = msg
		result.Log_message = msg_text
	case "gemini":
		user.Model = "gemini"
		msg_text := "Напишите свой вопрос:"
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		result.Message = msg
		result.Log_message = msg_text
	case "kandinsky":
		user.Model = "kandinsky"
		user.Stage = "text"
		msg_text := "Введите свой запрос:"
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		result.Message = msg
		result.Log_message = msg_text
	default:
		if upd.Message.From.UserName == "Art_Korn_39" {
			switch cmd {
			case "stop":
				os.Exit(1)
			case "updconf":
				LoadConfig()
				result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, "Config updated.")
			case "info":
				msg_text := fmt.Sprintf("Gemini: %d\nChatGPT: %d\nKandinsky: %d",
					counter_gemini, counter_chatgpt, counter_kandinsky)
				result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
			}
		}
	}

	return result

}

func ProcessText(text string, user *UserInfo, upd tgbotapi.Update) ResultOfRequest {

	var result ResultOfRequest

	switch user.Model {
	case "chatgpt":
		Operation := NewSQLOperation(user, upd, text)
		SQL_AddOperation(Operation)
		counter_chatgpt++

		msg_text := SendRequestToChatGPT(upd.Message.Text, user)
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "ChatGPT"
		result.Log_message = msg_text

	case "gemini":
		Operation := NewSQLOperation(user, upd, text)
		SQL_AddOperation(Operation)
		counter_gemini++

		msg_text := SendRequestToGemini(upd.Message.Text, user)
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "Gemini"
		result.Log_message = msg_text

	case "kandinsky": // пользователь ввел текст по модели kand

		return ProcessInputText_Kandinsky(text, user, upd)

	default:
		msg_text := "Не выбрана нейросеть для обработки запроса."
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "bot"
		result.Log_message = msg_text
	}

	return result

}
