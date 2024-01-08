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
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var WorkDir string //C:/DEV/GO/ai_telegram_bot

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

var (
	delay_upd            = time.Tick(time.Millisecond * 10)
	delay_ChatGPT        = time.Tick(time.Second * 15 / 10) // 40 RPM
	delay_Gemini         = time.Tick(time.Second * 12 / 11) // 55 RPM
	delay_Kandinsky      = time.Tick(time.Second * 3)       // 20 RPM
	delay_SaveUserStates = time.Tick(time.Minute * 1)       // 1 RPM
)

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

	if slices.Contains(Models, strings.ToLower(m.Text)) {
		return true
	}

	return m.IsCommand()

}

func MsgCommand(m *tgbotapi.Message) string {

	if slices.Contains(Models, strings.ToLower(m.Text)) {
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

<u>Последние обновления:</u>
<i>06.01.23 - 🏞 добавлена обработка картинок с вопросами в AI Gemini.</i>
<i>09.01.23 - 🎧 добавлена генерация аудио из текста в ChatGPT.</i>

Чтобы начать - просто выбери подходящую нейросеть и задай ей вопрос (или попроси сделать картинку), удачи 🔥`,
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

func GetDurationToNextDay() time.Duration {

	// тек. время по Мск
	now := time.Now().UTC().Add(3 * time.Hour)

	// добавили сутки
	tomorrow := now.Add(time.Hour * 24)

	// округлили до начала дня
	startDay := tomorrow.Truncate(time.Hour * 24)

	// сколько времени до начала дня
	// т.к. берется текущее время в UTC, то вычитаем 3 часа
	duration := time.Until(startDay) - (time.Hour * 3)

	return duration

}
