package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type config struct {
	TelegramBotToken  string
	OpenAIToken       string
	GeminiKey         string
	DailyLimitTokens  int
	DB_name           string
	DB_host           string
	DB_port           int
	DB_user           string
	DB_password       string
	CheckSubscription bool
	WhiteList         []string
}

func LoadConfig() {

	log.Println("Version: " + Version)

	file, err := os.OpenFile("config.txt", os.O_RDONLY, 0600)
	if err != nil {
		log.Println("Не удалось открыть файл config.txt")
		log.Println(err.Error())
		return
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Println("Не удалось прочитать данные из config.txt")
		log.Println(err.Error())
		return
	}

	err = json.Unmarshal(b, &Cfg)
	if err != nil {
		log.Println("Не удалось преобразовать в JSON файл config.txt")
		log.Println(err.Error())
		return
	}

	log.Println("Config download complete")

}
