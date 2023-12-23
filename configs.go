package main

import (
	"slices"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type Log struct {
	UserName string
	Text     string
	IsError  bool
}

type UserInfo struct {
	IsRunning   bool
	LastCommand string
}

type ResultOfRequest struct {
	Message     tgbotapi.Chattable
	Log_author  string
	Log_message string
}

func (u *UserInfo) SetIsRunning(v bool) {
	u.IsRunning = v
}

func MsgIsCommand(m *tgbotapi.Message) bool {

	if slices.Contains(arrayCMD, strings.ToLower(m.Text)) {
		return true
	}

	return m.IsCommand()

}

func MsgCommand(m *tgbotapi.Message) string {

	if slices.Contains(arrayCMD, strings.ToLower(m.Text)) {
		return strings.ToLower(m.Text)
	}

	return m.Command()

}
