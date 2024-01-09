package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	kand_Styles   = map[string]string{"Без стиля": "DEFAULT", "Art": "KANDINSKY", "4K": "UHD", "Anime": "ANIME"}
	kand_Model_id = "4"
)

// После команды "/kandinsky" или при вводе текста = "kandinsky"
func kand_start(user *UserInfo) {

	msgText := `Введите свой запрос:`
	SendMessage(user, msgText, button_RemoveKeyboard, "")

	user.Path = "kandinsky/text"

}

// После ввода запроса пользователем
func kand_text(user *UserInfo, text string) {

	if utf8.RuneCountInString(text) >= 900 {
		SendMessage(user, "Текст описания картинки не должен превышать 1000 символов.", nil, "")
		return
	}

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

	msgText := "Запущена генерация картинки, среднее время выполнения 25-30 секунд."
	SendMessage(user, msgText, button_RemoveKeyboard, "")

	Operation := SQL_NewOperation(user, "kandinsky", text, inputText)
	SQL_AddOperation(Operation)

	res, err := SendRequestToKandinsky(inputText, style, user)
	if err != nil {
		Logs <- NewLog(user, "kandinsky", Error, err.Error())
		SendMessage(user, res, button_kandNewgen, "")
	} else {
		caption := fmt.Sprintf(`Результат генерации по запросу "%s", стиль: "%s"`, inputText, text)
		err := SendPhotoMessage(user, res, caption, button_kandNewgen)
		if err != nil {
			Logs <- NewLog(user, "kandinsky", Error, "{ImgSend} "+err.Error())
			SendMessage(user, "При отправке картинки возникла ошибка, попробуйте ещё раз позже.", button_kandNewgen, "")
		}
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

// Максимальный размер текстового описания - 1000 символов.
// "censored": true (возвращая ошибку, когда запрос или изображение не соответствует)
func SendRequestToKandinsky(text string, style string, user *UserInfo) (result string, err error) {

	<-delay_Kandinsky

	scriptPath := WorkDir + "/scripts/generate_image.py"
	dataFolder := WorkDir + "/data"

	cmd := exec.Command(`python`,
		scriptPath,
		dataFolder,
		text,
		style,
		strconv.Itoa(int(user.ChatID)),
		kand_Model_id,
		Cfg.Kandinsky_Key,
		Cfg.Kandinsky_Secret,
	)

	if cmd.Err != nil {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, cmd.Err.Error())
		err = errors.New("{cmd} " + description)
		//Logs <- NewLog(user, "Kandinsky{cmd}", Error, description)
		return "Не удалось сгенерировать изображение. Попробуйте позже.", err
	}

	// Получение результата команды
	res, err2 := cmd.Output()

	if err2 != nil {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, err2.Error())
		err = errors.New("{cmd.Output()} " + description)
		//Logs <- NewLog(user, "Kandinsky{cmd.Output()}", Error, description)
		return "Не удалось сгенерировать изображение. Попробуйте изменить текст описания картинки.", err
	}

	pathToImage := strings.TrimSpace(string(res[:]))

	if pathToImage == "" {
		description := fmt.Sprintf("text: %s [%s]\nerror: %s", text, style, "скрипт вернул пустой путь до картинки (не успел получит ответ по api)")
		err = errors.New("{API} " + description)
		//Logs <- NewLog(user, "Kandinsky{API}", Error, description)
		return "Не удалось сгенерировать изображение. Попробуйте позже.", err
	}

	return pathToImage, nil

}
