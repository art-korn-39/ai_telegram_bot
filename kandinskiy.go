package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	styles_knd    = map[string]string{"Без стиля": "DEFAULT", "Art": "KANDINSKY", "4K": "UHD", "Anime": "ANIME"}
	buttons_style = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Без стиля"),
			tgbotapi.NewKeyboardButton("Art"),
			tgbotapi.NewKeyboardButton("4K"),
			tgbotapi.NewKeyboardButton("Anime"),
		),
	)
	button_newGenerate = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Новая генерация по тексту"),
		),
	)
)

func ProcessInputText_Kandinsky(text string, user *UserInfo, upd tgbotapi.Update) ResultOfRequest {

	var result ResultOfRequest

	switch user.Stage {
	case "NewGeneration": // предложение о новой генерации картинки
		user.Stage = "text"

		msg_text := "Введите свой запрос:"
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

		result.Message = msg
		result.Log_author = "bot"
		result.Log_message = msg_text

	case "text": // после ввода текста запроса
		user.InputText = text
		user.Stage = "style"

		msg_text := "Выберите стиль, в котором генерировать изображение."
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)
		msg.ReplyMarkup = buttons_style

		result.Message = msg
		result.Log_author = "kandinsky"
		result.Log_message = msg_text

	case "style": // после указания стиля картинки

		style, ok := styles_knd[text]
		if ok {

			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Запущена генерация картинки, она может занять 1-2 минуты.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			Bot.Send(msg)

			res, isError := SendRequestToKandinsky(user.InputText, style, upd.Message.Chat.ID)
			if isError {
				msg = tgbotapi.NewMessage(upd.Message.Chat.ID, res)
				msg.ReplyMarkup = button_newGenerate
				result.Message = msg
			} else {
				msg := tgbotapi.NewPhotoUpload(upd.Message.Chat.ID, res)
				msg.Caption = fmt.Sprintf(`Результат генерации по запросу "%s", стиль: "%s"`, user.InputText, text)
				msg.ReplyMarkup = button_newGenerate
				result.Message = msg
			}

			result.Log_author = "Kandinsky"
			result.Log_message = res

			user.Stage = "NewGeneration"
			user.InputText = ""

		} else {
			msg_text := "Пожалуйста, введите стиль из предложенных вариантов"
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg_text)

			result.Message = msg
			result.Log_author = "kandinsky"
			result.Log_message = msg_text
		}

	}

	return result

}

func SendRequestToKandinsky(text string, style string, userid int64) (result string, isError bool) {

	<-delay_Kandinsky

	_, callerFile, _, _ := runtime.Caller(0)
	dir := strings.ReplaceAll(filepath.Dir(callerFile), "\\", "/")
	scriptPath := dir + "/scripts/generate_image.py"
	dataFolder := dir + "/data"

	cmd := exec.Command(`python`,
		scriptPath,
		dataFolder,
		text,
		style,
		strconv.Itoa(int(userid)))

	if cmd.Err != nil {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, cmd.Err.Error())
		Logs <- Log{"Kandinsky{cmd}", description, true}
		return "Не удалось сгенерировать изображение. Попробуйте позже.", true
	}

	// Получение результата команды
	res, err := cmd.Output()

	if err != nil {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, cmd.Err.Error())
		Logs <- Log{"Kandinsky{cmd.Output()}", description, true}
		return "Не удалось сгенерировать изображение. Попробуйте позже.", true
	}

	pathToImage := strings.TrimSpace(string(res[:]))

	if pathToImage == "" {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, cmd.Err.Error())
		Logs <- Log{"Kandinsky{API}", description, true}
		return "Не удалось сгенерировать изображение. Попробуйте позже.", true
	}

	return pathToImage, false

}
