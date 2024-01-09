package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	kand_Styles = map[string]string{"Без стиля": "DEFAULT", "Art": "KANDINSKY", "4K": "UHD", "Anime": "ANIME"}
)

// После команды "/kandinsky" или при вводе текста = "kandinsky"
func kand_start(user *UserInfo) {

	msgText := `Введите свой запрос:`
	SendMessage(user, msgText, button_RemoveKeyboard, "")

	user.Path = "kandinsky/text"

}

// После ввода запроса пользователем
func kand_text(user *UserInfo, text string) {

	user.Options["text"] = text

	msgText := `Выберите стиль, в котором генерировать изображение.`
	SendMessage(user, msgText, buttons_kandStyles, "")

	user.Path = "kandinsky/text/style"

}

// После выбора стиля пользователем
func kand_style(user *UserInfo, text string) {

	style, ok := kand_Styles[text]
	if !ok {
		msgText := "Выберите стиль из предложенных вариантов."
		SendMessage(user, msgText, buttons_kandStyles, "")
		return
	}

	user.Options["style"] = style
	inputText := user.Options["text"]

	msgText := "Запущена генерация картинки, она может занять 1-2 минуты."
	SendMessage(user, msgText, button_RemoveKeyboard, "")

	Operation := SQL_NewOperation(user, "kandinsky", text, inputText)
	SQL_AddOperation(Operation)

	res, isError := SendRequestToKandinsky(inputText, style, user.ChatID)
	if isError {
		// в errors уже записали в самой функции "SendRequestToKandinsky()"
		SendMessage(user, res, button_kandNewgen, "")
	} else {
		Message := tgbotapi.NewPhotoUpload(user.ChatID, res)
		Message.Caption = fmt.Sprintf(`Результат генерации по запросу "%s", стиль: "%s"`, inputText, text)
		Message.ReplyMarkup = button_kandNewgen
		Bot.Send(Message)
		Logs <- NewLog(user, "kandinsky", Info, res)
	}

	user.Path = "kandinsky/text/style/newgen"

}

// После получения результата генерации
func kand_newgen(user *UserInfo, text string) {

	switch text {
	case "Изменить текст запроса":
		SendMessage(user, "Введите свой запрос:", button_RemoveKeyboard, "")
		user.Path = "kandinsky/text"
	case "Выбрать другой стиль":
		SendMessage(user, "Выберите стиль, в котором генерировать изображение.", buttons_kandStyles, "")
		user.Path = "kandinsky/text/style"
	default:
		// Предполагаем, что там новый вопрос к загруженным картинкам
		kand_text(user, text)
	}

}

func SendRequestToKandinsky(text string, style string, chatid int64) (result string, isError bool) {

	<-delay_Kandinsky

	scriptPath := WorkDir + "/scripts/generate_image.py"
	dataFolder := WorkDir + "/data"

	cmd := exec.Command(`python`,
		scriptPath,
		dataFolder,
		text,
		style,
		strconv.Itoa(int(chatid)))

	if cmd.Err != nil {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, cmd.Err.Error())
		Logs <- NewLog(nil, "Kandinsky{cmd}", Error, description)
		return "Не удалось сгенерировать изображение. Попробуйте позже.", true
	}

	// Получение результата команды
	res, err := cmd.Output()

	if err != nil {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, err.Error())
		Logs <- NewLog(nil, "Kandinsky{cmd.Output()}", Error, description)
		return "Не удалось сгенерировать изображение. Попробуйте изменить текст описания картинки.", true
	}

	pathToImage := strings.TrimSpace(string(res[:]))

	if pathToImage == "" {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, "скрипт вернул пустой путь до картинки")
		Logs <- NewLog(nil, "Kandinsky{API}", Error, description)
		return "Не удалось сгенерировать изображение. Попробуйте позже.", true
	}

	return pathToImage, false

}
