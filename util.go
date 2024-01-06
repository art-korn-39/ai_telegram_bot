package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var WorkDir string

type config struct {
	TelegramBotToken  string
	OpenAIToken       string
	GeminiKey         string
	DB_name           string
	DB_host           string
	DB_port           int
	DB_user           string
	DB_password       string
	CheckSubscription bool
	WhiteList         []string
}

func init() {
	_, callerFile, _, _ := runtime.Caller(0)
	WorkDir = strings.ReplaceAll(filepath.Dir(callerFile), "\\", "/")
}

func LoadConfig() {

	log.Println("Version: " + Version)

	file, err := os.OpenFile("config.txt", os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(b, &Cfg)

	log.Println("Config download complete")

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

	return fmt.Sprintf(`
Привет, %s! 👋

Я бот для работы с нейросетями (v%s).
С моей помощью ты можешь использовать следующие модели:

<b>ChatGPT</b> - используется для генерации текста.
<b>Gemini</b> - аналог ChatGPT от компании Google.
<b>Kandinsky</b> - используется для создания изображений по текстовому описанию.

<i>Начиная с версии бота "2.0.1" добавлена возможность отправки картинок с вопросами в AI Gemini.</i>

Чтобы начать - просто выбери подходящую нейросеть и задай ей вопрос (или попроси сделать картинку), удачи!🔥`,
		user, Version)

}

func mapToJSON(m map[string]string) string {

	result, _ := json.Marshal(m)
	return string(result)

}

func JSONtoMap(JSON string) map[string]string {

	result := map[string]string{}

	resBytes := []byte(JSON)
	//resBytes := io.ReadAll(res.Body)
	json.Unmarshal(resBytes, &result)

	return result

}
