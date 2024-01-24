package main

import (
	"testing"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func TestValidMessage(t *testing.T) {
	upd := &tgbotapi.Update{
		Message:       nil,
		CallbackQuery: nil,
	}
	result := ValidMessage(upd)
	if result != false {
		t.Errorf("ValidMessage(upd) = %t; expected false", result)
	}

	upd = &tgbotapi.Update{
		Message:       &tgbotapi.Message{From: &tgbotapi.User{IsBot: true}},
		CallbackQuery: nil,
	}
	result = ValidMessage(upd)
	if result == true {
		t.Errorf("ValidMessage(upd) = %t; expected false", result)
	}
}

func TestSQL(t *testing.T) {

	LoadConfig()
	SQL_Connect()
	defer db.Close()

	SQL_LoadUserStates()
	SQL_SaveUserStates()

	u := T_GetUser()
	Operation := SQL_NewOperation(u, "gemini", "img", "test operation")
	SQL_AddOperation(Operation)

	SQL_AddLog(NewLog(u, "", Warning, "test add log"))

	_, IsErr := SQL_CountOfUserOperations(u)
	if IsErr {
		t.Errorf("SQL_CountOfUserOperations(u)")
	}

}

func T_GetUser() *UserInfo {
	return &UserInfo{
		Username:           "test",
		ChatID:             00000001,
		IsRunning:          true,
		Language:           "en",
		Path:               "gemini/type/image/text",
		Tokens_used_gpt:    1000,
		Requests_today_gen: 12}
}
