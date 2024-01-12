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

	defer FinishGorutine(nil, "panic in SaveLogs()", false)

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

	defer FinishGorutine(nil, "panic in SaveUserStates()", false)

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

	defer FinishGorutine(nil, "panic in ClearTokensEveryDay()", false)

	for {

		duration := GetDurationToNextDay()

		fmt.Printf("Следующая операция по очистке токенов через: %f ч.\n", duration.Hours())

		// Ожидание до начала след. дня (Мск 00:00)
		<-time.After(duration)

		// Очистка токенов у пользователей
		for _, u := range ListOfUsers {
			u.Tokens_used_gpt = 0
			u.Requests_today_gen = 0
		}

		SQL_SaveUserStates()
		Logs <- NewLog(nil, "System", Info, "Счетчик использованных токенов очищен")

	}

}

func Kandinsky_CheckModelID() {

	defer FinishGorutine(nil, "panic in Kandinsky_CheckModelID()", false)

	delay := time.Tick(time.Minute * 30)

	for {

		url := "https://api-key.fusionbrain.ai/key/api/v1/models"
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			Logs <- NewLog(nil, "kandinsky", Error, "Не удалось получить model_id {1}")
			Logs <- NewLog(nil, "kandinsky", Error, err.Error())
			return
		}

		req.Header.Add("X-Key", "Key "+Cfg.Kandinsky_Key)
		req.Header.Add("X-Secret", "Secret "+Cfg.Kandinsky_Secret)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			Logs <- NewLog(nil, "kandinsky", Error, "Не удалось получить model_id {2}")
			Logs <- NewLog(nil, "kandinsky", Error, err.Error())
			return
		}
		defer res.Body.Close()

		resBytes, err := io.ReadAll(res.Body)
		if err != nil {
			Logs <- NewLog(nil, "kandinsky", Error, "Не удалось получить model_id {3}")
			Logs <- NewLog(nil, "kandinsky", Error, err.Error())
			return
		}

		var dat []map[string]any
		err = json.Unmarshal(resBytes, &dat)
		if err != nil {
			Logs <- NewLog(nil, "kandinsky", Error, "Не удалось получить model_id {4}")
			Logs <- NewLog(nil, "kandinsky", Error, err.Error())
			return
		}

		kand_Model_id = strconv.Itoa(int(dat[0]["id"].(float64)))

		Logs <- NewLog(nil, "kandinsky", Info, "Значение model_id обновлено ["+kand_Model_id+"]")

		<-delay

	}

}
