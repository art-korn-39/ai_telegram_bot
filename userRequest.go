package main

import (
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type ResultOfRequest struct {
	Message     tgbotapi.Chattable
	Log_author  string
	Log_message string
}

func ProcessCommand(cmd string, upd tgbotapi.Update) ResultOfRequest {

	var result ResultOfRequest
	result.Log_author = "bot"

	switch cmd {
	case "start":
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, start(upd.Message.Chat.FirstName))
		msg.ParseMode = "HTML"

		var buttons = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Gemini"),
				tgbotapi.NewKeyboardButton("Kandinsky"),
				tgbotapi.NewKeyboardButton("ChatGPT"),
			),
		)
		msg.ReplyMarkup = buttons

		result.Message = msg
		result.Log_message = "/start for " + upd.Message.Chat.UserName
	case "chatgpt":
		msg_text := "Напишите свой вопрос:"
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_message = msg_text
	case "gemini":
		msg_text := "Напишите свой вопрос:"
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_message = msg_text
	case "kandinsky":
		msg_text := "Введите свой запрос:"
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_message = msg_text
	default:
		if upd.Message.From.UserName == "Art_Korn_39" {
			switch cmd {
			case "stop":
				os.Exit(1)
			case "updconf":
				LoadConfig()
				result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, "Config updated.")
			}
		}
	}

	return result

}

func ProcessText(text string, cmd string, upd tgbotapi.Update) ResultOfRequest {

	var result ResultOfRequest

	switch cmd {
	case "chatgpt":
		msg_text := SendRequestToChatGPT(upd.Message.Text)
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "ChatGPT"
		result.Log_message = msg_text

	case "gemini":
		msg_text := SendRequestToGemini(upd.Message.Text)
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "Gemini"
		result.Log_message = msg_text

	case "kandinsky":
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Запущена генерация картинки, она может занять 1-2 минуты.")
		Bot.Send(msg)

		res, isError := SendRequestToKandinsky(upd.Message.Text, upd.Message.Chat.ID)
		if isError {
			result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, res)
		} else {
			result.Message = tgbotapi.NewPhotoUpload(upd.Message.Chat.ID, res)
		}
		result.Log_author = "Kandinsky"
		result.Log_message = res

	case "":
		msg_text := "Не выбрана нейросеть для обработки запроса."
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "bot"
		result.Log_message = msg_text

	case "start":
		msg_text := "Не выбрана нейросеть для обработки запроса."
		result.Message = tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		result.Log_author = "bot"
		result.Log_message = msg_text
	}

	return result

}
