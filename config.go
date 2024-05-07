package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type config struct {
	TelegramBotToken string
	OpenAIToken      string
	GeminiKey        string
	Kandinsky_Key    string
	Kandinsky_Secret string
	Stability_Key    string
	Replicate_Key    string
	Faceswap_Version string
	Faceswap_id      string

	TPD_gpt           int
	RPD_gen           int
	RPD_sdxl          int
	RPD_fs            int
	TPD_advanced_gpt  int
	RPD_advanced_sdxl int
	RPD_advanced_fs   int

	Gen_UseStream bool
	Gen_Rip       bool

	DB_name     string
	DB_host     string
	DB_port     int
	DB_user     string
	DB_password string

	DaysForAdvancedStatus         int
	CheckSubscription             bool
	OperationsWithoutSubscription int
	Debug                         bool // для функции recovery()
	WhiteList                     []string
	Admins                        []string
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

// Получить доступный суточный лимит использования ChatGPT по пользователю
func Get_TPD_gpt(u *UserInfo) (res int) {
	if u.Level == Basic {
		res = Cfg.TPD_gpt
	} else if u.Level == Advanced {
		res = Cfg.TPD_advanced_gpt
	}
	return
}

// Получить доступный суточный лимит использования SDXL по пользователю
func Get_RPD_sdxl(u *UserInfo) (res int) {
	if u.Level == Basic {
		res = Cfg.RPD_sdxl
	} else if u.Level == Advanced {
		res = Cfg.RPD_advanced_sdxl
	}
	return
}

// Получить доступный суточный лимит использования faceswap по пользователю
func Get_RPD_fs(u *UserInfo) (res int) {
	if u.Level == Basic {
		res = Cfg.RPD_fs
	} else if u.Level == Advanced {
		res = Cfg.RPD_advanced_fs
	}
	return
}
