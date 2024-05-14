package main

import (
	"fmt"
	"time"
)

func GetInfo(isRO bool) string {

	Now := time.Now().UTC().Add(3 * time.Hour) // 03.02.24 14:50

	// Если это регламентная операция, то берем предыдущий день
	if isRO {
		Now = Now.AddDate(0, 0, -1) // 02.02.24 14:50
	}

	StartDay := time.Date(Now.Year(), Now.Month(), Now.Day(), 0, 0, 0, 0, Now.Location()) // 03.02.24 00:00
	December29 := time.Date(2023, 12, 29, 0, 0, 0, 0, time.Local)
	December25 := time.Date(2023, 12, 25, 0, 0, 0, 0, time.Local)

	result_dec25, err0 := SQL_GetOperationsFromDate(December25)
	if err0 != "" {
		return err0
	}

	result_dec29, err1 := SQL_GetOperationsFromDate(December29)
	if err1 != "" {
		return err1
	}

	result_Today, err3 := SQL_GetOperationsFromDate(StartDay)
	if err3 != "" {
		return err3
	}

	newUsersToday, err4 := SQL_GetNewUsersForDay(StartDay)
	if err4 != "" {
		return err4
	}

	errors_Today, err5 := SQL_GetErrorsForDay(StartDay)
	if err5 != "" {
		return err5
	}

	alltimeOP := result_dec25["gemini"] + result_dec25["chatgpt"] + result_dec25["kandinsky"] + result_dec25["sdxl"] + result_dec25["faceswap"]
	todayOP := result_Today["gemini"] + result_Today["chatgpt"] + result_Today["kandinsky"] + result_Today["sdxl"] + result_Today["faceswap"]

	return fmt.Sprintf(`
All time [%d]
Users: %d	
Gen: %d | GPT: %d | Kand: %d | SDXL: %d | FS: %d

Today [%d]
Users: %d (new: %d)
Gen: %d (%d) | GPT: %d (%d%%) | Kand: %d (%d) | SDXL: %d (%d) | FS: %d (%d%%)`,
		alltimeOP,
		result_dec29["users"],
		result_dec25["gemini"], result_dec25["chatgpt"], result_dec25["kandinsky"], result_dec25["sdxl"], result_dec25["faceswap"],
		todayOP,
		result_Today["users"], newUsersToday,
		result_Today["gemini"], GetPartOfErrors("gemini", result_Today, errors_Today),
		result_Today["chatgpt"], GetPartOfErrors("chatgpt", result_Today, errors_Today),
		result_Today["kandinsky"], GetPartOfErrors("kandinsky", result_Today, errors_Today),
		result_Today["sdxl"], GetPartOfErrors("sdxl", result_Today, errors_Today),
		result_Today["faceswap"], GetPartOfErrors("faceswap", result_Today, errors_Today),
	)

}

func GetPartOfErrors(model string, operations, errors map[string]int) int {

	sum := operations[model] + errors[model]
	if sum == 0 {
		return 0
	} else {
		return operations[model] * 100 / sum
	}

}
