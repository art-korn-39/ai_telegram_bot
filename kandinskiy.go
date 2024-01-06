package main

import (
	"fmt"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// del
func ProcessInputText_Kandinsky(text string, user *UserInfo) ResultOfRequest {

	var result ResultOfRequest

	//	switch user.Stage {
	switch "123" {
	case "NewGeneration": // предложение о новой генерации картинки
		//		user.Stage = "text"

		msg_text := "Введите свой запрос:"
		msg := tgbotapi.NewMessage(user.ChatID, msg_text)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

		result.Message = msg
		result.Log_author = "bot"
		result.Log_message = msg_text

	case "text": // после ввода текста запроса
		//		user.InputText = text
		//		user.Stage = "style"

		msg_text := "Выберите стиль, в котором генерировать изображение."
		msg := tgbotapi.NewMessage(user.ChatID, msg_text)
		msg.ReplyMarkup = buttons_style

		result.Message = msg
		result.Log_author = "kandinsky"
		result.Log_message = msg_text

	case "style": // после указания стиля картинки

		style, ok := styles_knd[text]
		if ok {

			msg := tgbotapi.NewMessage(user.ChatID, "Запущена генерация картинки, она может занять 1-2 минуты.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			Bot.Send(msg)

			//			Operation := SQL_NewOperation(user, "["+text+"] "+user.InputText)
			//			SQL_AddOperation(Operation)

			//			res, isError := SendRequestToKandinsky(user.InputText, style, user.ChatID)
			res, isError := SendRequestToKandinsky("123", style, user.ChatID)
			if isError {
				msg = tgbotapi.NewMessage(user.ChatID, res)
				msg.ReplyMarkup = button_newGenerate
				result.Message = msg
			} else {
				msg := tgbotapi.NewPhotoUpload(user.ChatID, res)
				//msg.Caption = fmt.Sprintf(`Результат генерации по запросу "%s", стиль: "%s"`, user.InputText, text)
				msg.Caption = fmt.Sprintf(`Результат генерации по запросу "%s", стиль: "%s"`, "123", text)
				msg.ReplyMarkup = button_newGenerate
				result.Message = msg
			}

			result.Log_author = "Kandinsky"
			result.Log_message = res

			//			user.Stage = "NewGeneration"
			//			user.InputText = ""

		} else {
			msg_text := "Пожалуйста, выберите стиль из предложенных вариантов"
			msg := tgbotapi.NewMessage(user.ChatID, msg_text)

			result.Message = msg
			result.Log_author = "kandinsky"
			result.Log_message = msg_text
		}

	}

	return result

}
