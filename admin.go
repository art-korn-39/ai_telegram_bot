package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func HandleAdminCommand(u *UserInfo, cmd string) {

	switch cmd {
	case "info":
		SendMessage(u, GetInfo(false), GetButton(btn_RemoveKeyboard, ""), "")
	case "updconf":
		LoadConfig()
		SendMessage(u, "Config updated.", GetButton(btn_RemoveKeyboard, ""), "")
	case "sendMessageToAllUsers":
		SendMessageToAllUsers(u)
	default:
		SendMessage(u, GetText(MsgText_UnknownCommand, u.Language), GetButton(btn_Models, ""), "")
	}

}

func SendMessageToAllUsers(u *UserInfo) {

	folder := WorkDir + "/messageToAll"

	var toAllUsers bool

	// если есть файл TODO.txt, то отправляем всем (и файл удаляем), иначе только в тестовые акки
	filepathTODO := folder + "/TODO.txt"
	if _, err := os.Stat(filepathTODO); errors.Is(err, os.ErrNotExist) {
		toAllUsers = false
	} else {
		toAllUsers = true
	}

	// text_ru
	file, err := os.OpenFile(folder+"/text_ru.txt", os.O_RDONLY, 0600)
	if err != nil {
		SendMessage(u, "Не удалось открыть файл text_ru.txt", nil, "")
		SendMessage(u, err.Error(), nil, "")
		return
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		SendMessage(u, "Не удалось прочитать данные из text_ru.txt", nil, "")
		SendMessage(u, err.Error(), nil, "")
		return
	}

	text_ru := string(b)

	// text_en
	file, err = os.OpenFile(folder+"/text_en.txt", os.O_RDONLY, 0600)
	if err != nil {
		SendMessage(u, "Не удалось открыть файл text_en.txt", nil, "")
		SendMessage(u, err.Error(), nil, "")
		return
	}
	defer file.Close()

	b, err = io.ReadAll(file)
	if err != nil {
		SendMessage(u, "Не удалось прочитать данные из text_en.txt", nil, "")
		SendMessage(u, err.Error(), nil, "")
		return
	}

	text_en := string(b)

	// image.png
	filepath := folder + "/image.png"
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		SendMessage(u, "файл картинки не найден", nil, "")
		return
	}

	users, isErr := SQL_GetAllUsers()
	if isErr {
		SendMessage(u, "Ошибка при получении всех пользователей", nil, "")
		return
	}
	SendMessage(u, fmt.Sprintf("Найдено %d чатов, начинаю отправку сообщений.", len(users)), nil, "")

	if toAllUsers {
		os.Remove(filepathTODO)
	}

	counter := 0
	for _, user := range users {

		if toAllUsers || (user.ChatID == 403059287 ||
			user.ChatID == 6648171361 ||
			user.ChatID == 609614322 ||
			user.ChatID == 0) {

			var caption string
			if user.Language == "ru" || user.Language == "uk" {
				caption = text_ru
			} else {
				caption = text_en
			}

			Message := tgbotapi.NewPhotoUpload(user.ChatID, filepath)
			Message.Caption = caption
			Message.ReplyMarkup = GetButton(btn_Models, "")
			Message.ParseMode = "HTML"
			_, err := Bot.Send(Message)

			if err != nil {
				SendMessage(u, fmt.Sprintf("{id:%d} %s", user.ChatID, err.Error()), nil, "")
			} else {
				counter++
			}

		}
	}

	SendMessage(u, "Отправка сообщений завершена.", nil, "")
	SendMessage(u, fmt.Sprintf("%d пользователей получили сообщение.", counter), nil, "")

}
