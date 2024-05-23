package main

import (
	"testing"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
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

func TestAccessIsAllowed(t *testing.T) {

	LoadConfig()
	StartBot()

	upd := tgbotapi.Update{
		Message:       &tgbotapi.Message{From: &tgbotapi.User{IsBot: true}},
		CallbackQuery: nil,
	}
	u := T_GetUser()

	if AccessIsAllowed(upd, u) && Cfg.CheckSubscription {
		t.Errorf("AccessIsAllowed() is true for chat_id = %d", u.ChatID)
	}

}

func TestGeminiBag(t *testing.T) {

	LoadConfig()
	NewConnectionGemini()

	text := "اگر خالص حقوق پرداختی توسط پیمانکار بعد از کسر ده درصد مالیات،بیمه و یازده درصد پیش پرداخت مبلغ چهارده میلیون و چهارصد هزار ریال باشد،مبلغ بدهکار حساب پیلم در جریان چقدر میشه؟"

	cs := gen_TextModel.StartChat()
	resp, err := cs.SendMessage(gen_ctx, genai.Text(text))
	Unused(resp, err)

}

func TestSQL(t *testing.T) {

	LoadConfig()
	SQL_Connect()
	defer db.Close()

	SQL_LoadUserStates()
	SQL_SaveUserStates()

	u := T_GetUser()

	Operation := SQL_NewOperation(u, "gemini", "img", "", "test operation")
	SQL_AddOperation(Operation)

	SQL_AddLog(NewLog(u, "", Warning, "test add log"))

	_, IsErr := SQL_CountOfUserOperations(u)
	if IsErr {
		t.Errorf("SQL_CountOfUserOperations(u)")
	}

	_, IsErr = SQL_GetFirstDate(u)
	if IsErr {
		t.Errorf("SQL_GetFirstDate(u)")
	}

	_, IsErr = SQL_UserDayStreak(u)
	if IsErr {
		t.Errorf("SQL_UserDayStreak(u)")
	}

}

func T_GetUser() *UserInfo {
	return &UserInfo{
		Username:  "test",
		ChatID:    00000001,
		IsRunning: true,
		Language:  "en",
		Path:      "gemini/type/image/text",
		Usage:     Usage{Gen10: 12, SDXL: 1, GPT: 1000},
		Level:     Basic,
	}
}

func T_GetUser2() *UserInfo {
	return &UserInfo{
		Username:  "test",
		ChatID:    104868607, // 403059287
		IsRunning: true,
		Language:  "en",
		Path:      "gemini/type/image/text",
		Usage:     Usage{Gen10: 12, SDXL: 1, GPT: 1000},
		Level:     Advanced,
	}
}
