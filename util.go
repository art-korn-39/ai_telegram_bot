package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var WorkDir string //C:/DEV/GO/ai_telegram_bot

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

func MapToJSON(m map[string]string) string {

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
	now := MskTimeNow()

	// добавили сутки
	tomorrow := now.Add(time.Hour * 24)

	// округлили до начала дня
	startDay := tomorrow.Truncate(time.Hour * 24)

	// сколько времени до начала дня
	// т.к. берется текущее время в UTC, то вычитаем 3 часа
	duration := time.Until(startDay) - (time.Hour * 3)

	return duration

}

func CreateFile(filename string, data io.Reader) error {

	outFile, _ := os.Create(filename)
	defer outFile.Close()
	_, err := io.Copy(outFile, data) //res = io.ReadClose
	if err != nil {
		return err
	}

	return nil

}
