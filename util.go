package main

import (
	"slices"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type ResultOfRequest struct {
	Message     tgbotapi.Chattable
	Log_author  string
	Log_message string
}

func MsgIsCommand(m *tgbotapi.Message) bool {

	if slices.Contains(arrayCMD, strings.ToLower(m.Text)) {
		return true
	}

	return m.IsCommand()

}

func MsgCommand(m *tgbotapi.Message) string {

	if slices.Contains(arrayCMD, strings.ToLower(m.Text)) {
		return strings.ToLower(m.Text)
	}

	return m.Command()

}

func start(user string) string {

	t := `
Привет, %user! 👋

Я бот для работы с нейросетями.
С моей помощью ты можешь использовать следующие модели:

<b>ChatGPT</b> - используется для генерации текста.
<b>Gemini</b> - аналог ChatGPT от компании Google.
<b>Kandinsky</b> - используется для создания изображений по текстовому описанию.

Чтобы начать - просто выбери подходящую нейросеть и задай ей вопрос (или попроси сделать картинку), удачи!🔥`

	return strings.ReplaceAll(t, "%user", user)

}
