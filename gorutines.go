package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func SaveLogs() {

	for v := range Logs {
		if v.Text != "" {

			// время по Мск час. поясу
			timeNow := v.Date.Format(time.DateTime)

			// date - username - text
			fmt.Printf("(%s) [%s] %s\n", timeNow, v.Author, v.Text)

			// если это ошибка, то записываем в отдельный файл
			if v.Level == FatalError || v.Level == Error {
				WriteIntoFile(timeNow, v.ChatID, v.Author, v.Text)
			}

			SQL_AddLog(v)
		}
	}
}

func SaveUserStates() {

	delay := time.Tick(time.Minute * 1) // 1 RPM

	for {
		<-delay
		if UserInfoChanged {
			UserInfoChanged = false
			SQL_SaveUserStates()
		}
	}

}

func ClearTokensEveryDay() {

	for {

		duration := GetDurationToNextDay()

		fmt.Println("до след. итерации часов:", duration.Hours())

		// Ожидание до начала след. дня (Мск 00:00)
		<-time.After(duration)

		// Очистка токенов у пользователей
		for _, u := range ListOfUsers {
			u.Tokens_used_gpt = 0
		}

		SQL_SaveUserStates()
		Logs <- NewLog(nil, "System", Info, "tokens = 0")

	}

}

func Kandinsky_CheckModelID() {

	delay := time.Tick(time.Minute * 30)

	for {

		url := "https://api-key.fusionbrain.ai/key/api/v1/models"
		req, _ := http.NewRequest(http.MethodGet, url, nil)

		req.Header.Add("X-Key", "Key "+Cfg.Kandinsky_Key)
		req.Header.Add("X-Secret", "Secret "+Cfg.Kandinsky_Secret)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			Logs <- NewLog(nil, "kandinsky", 1, "Не удалось получить model_id")
			return
		}
		defer res.Body.Close()

		resBytes, _ := io.ReadAll(res.Body)
		var dat []map[string]any
		json.Unmarshal(resBytes, &dat)

		kand_Model_id = strconv.Itoa(int(dat[0]["id"].(float64)))

		Logs <- NewLog(nil, "kandinsky", 1, "Значение model_id обновлено {"+kand_Model_id+"}")

		<-delay

	}

}
