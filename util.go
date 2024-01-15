package main

import (
	"encoding/json"
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
	delay_upd       = time.Tick(time.Millisecond * 10)
	delay_ChatGPT   = time.Tick(time.Second * 15 / 10) // 40 RPM
	delay_Gemini    = time.Tick(time.Second * 12 / 11) // 55 RPM
	delay_Kandinsky = time.Tick(time.Second * 3)       // 20 RPM
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
