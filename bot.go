package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
)

func SendMessage(user *UserInfo, text string, ReplyMarkup any, ParseMod string) {

	// Если длина сообщения превышает лимиты телеграма
	lenght := utf8.RuneCountInString(text)
	if lenght > 4000 {
		text1 := SubString(text, 0, lenght/2)
		SendMessage(user, text1, ReplyMarkup, ParseMod)

		time.Sleep(time.Microsecond * 10)

		text2 := SubString(text, lenght/2, lenght)
		SendMessage(user, text2, ReplyMarkup, ParseMod)
		return
	}

	Message := tgbotapi.NewMessage(user.ChatID, text)
	Message.ReplyMarkup = ReplyMarkup
	Message.ParseMode = ParseMod
	Bot.Send(Message)

	// Общий лог, пишем сюда все ответы пользователю
	// Логи с ошибками пишем в месте их возникновения
	if user.Path == "start" {
		text = "/start for " + user.Username
	} else if user.Path == "info" {
		text = "/info for " + user.Username
	} else if user.Path == "account" {
		text = "/account for " + user.Username
	}
	Logs <- NewLog(user, "bot", Info, text)

}

func SendAudioMessage(user *UserInfo, filename string, caption string, ReplyMarkup any) error {

	Message := tgbotapi.NewAudioUpload(user.ChatID, filename)
	Message.Caption = SubString(caption, 0, 1000)
	Message.ReplyMarkup = ReplyMarkup
	_, err := Bot.Send(Message)

	Logs <- NewLog(user, "bot", Info, filename)

	return err

}

func SendPhotoMessage(user *UserInfo, filename string, caption string, ReplyMarkup any) error {

	Message := tgbotapi.NewPhotoUpload(user.ChatID, filename)
	Message.Caption = SubString(caption, 0, 1000)
	Message.ReplyMarkup = ReplyMarkup
	_, err := Bot.Send(Message)

	Logs <- NewLog(user, "bot", Info, filename)

	return err

}

func SendFileMessage(user *UserInfo, filename string, caption string, ReplyMarkup any) error {

	Message := tgbotapi.NewDocumentUpload(user.ChatID, filename)
	Message.Caption = SubString(caption, 0, 1000)
	Message.ReplyMarkup = ReplyMarkup
	_, err := Bot.Send(Message)

	Logs <- NewLog(user, "bot", Info, filename)

	return err

}

func DownloadFile(FileID, name string) (string, error) {

	dataFolder := WorkDir + "/data"

	fileURL, err := Bot.GetFileDirectURL(FileID)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(fileURL)

	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := os.Create(dataFolder + "/" + name + ext)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Downloading complete")
	return file.Name(), nil

}

func StartBot() {

	var err error
	Bot, err = tgbotapi.NewBotAPI(Cfg.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	clientOpenAI = openai.NewClient(Cfg.OpenAIToken)
	NewConnectionGemini()

}
