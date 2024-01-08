package main

import (
	"fmt"
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
	for {
		<-delay_SaveUserStates
		if UserInfoChanged {
			UserInfoChanged = false
			SQL_SaveUserStates()
		}
	}
}

func ClearTokensEveryDay() {

	for {

		// // тек. время по Мск
		// now := time.Now().UTC().Add(3 * time.Hour)

		// // добавили сутки
		// tomorrow := now.Add(time.Hour * 24)

		// // округлили до начала дня
		// startDay := tomorrow.Truncate(time.Hour * 24)

		// // сколько времени до след. итерации
		// // т.к. берется текущее время в UTC, то вычитаем 3 часа
		// duration := time.Until(startDay) - (time.Hour * 3)

		duration := GetDurationToNextDay()

		fmt.Println("до след. итерации часов:", duration.Hours())

		// Ожидание до начала след. дня (Мск 00:00)
		<-time.After(duration)

		// Очистка токенов у пользователей
		for _, u := range ListOfUsers {
			u.TokensUsed_ChatGPT = 0
		}

		SQL_SaveUserStates()
		Logs <- NewLog(nil, "System", Info, "tokens = 0")

	}

}
