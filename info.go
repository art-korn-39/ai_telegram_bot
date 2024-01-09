package main

import (
	"fmt"
	"time"
)

func GetInfo() string {

	Now := time.Now().UTC().Add(3 * time.Hour)
	Yesterday := Now.AddDate(0, 0, -1)
	StartDay := time.Date(Now.Year(), Now.Month(), Now.Day(), 0, 0, 0, 0, Now.Location())
	December30 := time.Date(2023, 12, 30, 0, 0, 0, 0, time.Local)
	December25 := time.Date(2023, 12, 25, 0, 0, 0, 0, time.Local)

	result_dec25, err0 := SQL_GetInfoOnDate(December25)
	if err0 != "" {
		return err0
	}

	result_dec30, err1 := SQL_GetInfoOnDate(December30)
	if err1 != "" {
		return err1
	}

	result_24h, err2 := SQL_GetInfoOnDate(Yesterday)
	if err2 != "" {
		return err2
	}

	result_Today, err3 := SQL_GetInfoOnDate(StartDay)
	if err3 != "" {
		return err3
	}

	return fmt.Sprintf(
		`All time
Gemini: %d | ChatGPT: %d | Kandinsky: %d

From 30.12.23
Users: %d
Gemini: %d | ChatGPT: %d | Kandinsky: %d

Last 24h
Users: %d
Gemini: %d | ChatGPT: %d | Kandinsky: %d

Today
Users: %d
Gemini: %d | ChatGPT: %d | Kandinsky: %d`,
		result_dec25["gemini"], result_dec25["chatgpt"], result_dec25["kandinsky"],
		result_dec30["users"], result_dec30["gemini"], result_dec30["chatgpt"], result_dec30["kandinsky"],
		result_24h["users"], result_24h["gemini"], result_24h["chatgpt"], result_24h["kandinsky"],
		result_Today["users"], result_Today["gemini"], result_Today["chatgpt"], result_Today["kandinsky"])

}
