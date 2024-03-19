package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

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
	case "clearTokens":
		ClearTokens(u)
	default:
		SendMessage(u, GetText(MsgText_UnknownCommand, u.Language), GetButton(btn_Models, ""), "")
	}

}

// 2466 за 9m 32s -> 250/мин (прод)
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
	file, err := os.OpenFile(folder+"/text_ru.str", os.O_RDONLY, 0600)
	if err != nil {
		SendMessage(u, "Не удалось открыть файл text_ru.str", nil, "")
		SendMessage(u, err.Error(), nil, "")
		return
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		SendMessage(u, "Не удалось прочитать данные из text_ru.str", nil, "")
		SendMessage(u, err.Error(), nil, "")
		return
	}

	text_ru := string(b)

	// text_en
	file, err = os.OpenFile(folder+"/text_en.str", os.O_RDONLY, 0600)
	if err != nil {
		SendMessage(u, "Не удалось открыть файл text_en.str", nil, "")
		SendMessage(u, err.Error(), nil, "")
		return
	}
	defer file.Close()

	b, err = io.ReadAll(file)
	if err != nil {
		SendMessage(u, "Не удалось прочитать данные из text_en.str", nil, "")
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

	if toAllUsers {
		os.Remove(filepathTODO)
	} else {
		users = []*UserInfo{
			{ChatID: 403059287, Language: "ru"},  // art_korn_39
			{ChatID: 6648171361, Language: "en"}, // apolo39
			{ChatID: 609614322},                  // art_korneev
			{ChatID: 0},
		}
	}

	SendMessage(u, fmt.Sprintf("Найдено %d чатов, начинаю отправку сообщений.", len(users)), nil, "")

	M := tgbotapi.NewPhotoUpload(u.ChatID, filepath)
	M.Caption = "Отправка картинки для получения ID"
	res, err := Bot.Send(M)

	if err != nil {
		log.Println(err)
		return
	}

	photos := *res.Photo
	FileID := photos[len(photos)-1].FileID

	delay_msg := time.Tick(20 * time.Millisecond) // limit: 30 req per sec

	buttons := GetButton(btn_Models, "")

	counter_done := 0
	counter_fail := 0
	for _, user := range users {

		<-delay_msg

		var caption string
		if user.Language == "ru" || user.Language == "uk" {
			caption = text_ru
		} else {
			caption = text_en
		}

		Message := tgbotapi.NewPhotoShare(user.ChatID, FileID)
		Message.Caption = caption
		Message.ReplyMarkup = buttons
		Message.ParseMode = "HTML"
		_, err := Bot.Send(Message)

		if err != nil {
			Logs <- NewLog(user, "SendMessageToAllUsers()", Error, err.Error())
			counter_fail++
		} else {
			counter_done++
		}

		sum := counter_done + counter_fail
		if sum%100 == 0 {
			SendMessage(u, fmt.Sprintf("Done: %d | Fail: %d", counter_done, counter_fail), nil, "")
		}
	}

	SendMessage(u, "Отправка сообщений завершена.", nil, "")
	SendMessage(u, fmt.Sprintf("%d пользователей получили сообщение.", counter_done), nil, "")
	SendMessage(u, fmt.Sprintf("%d заблокировано.", counter_fail), nil, "")

}

func ClearTokens(u *UserInfo) {

	for _, u := range ListOfUsers {
		u.Mutex.Lock()
		u.ClearTokens()
		u.Mutex.Unlock()
	}

	SendMessage(u, "Лимиты очищены у всех пользователей.", nil, "")

}
