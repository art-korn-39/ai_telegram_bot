package main

import (
	"sync"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type UserInfo struct {
	IsRunning   bool
	LastCommand string
	Mutex       sync.Mutex
}

func (u *UserInfo) CheckUserLock(upd tgbotapi.Update) (isLocking bool) {

	// Устанавливаем блокировку по объекту UserInfo
	u.Mutex.Lock()
	// Разблокировку ставим через defer, чтобы она не стала вечной, если в методе произойдёт panic
	defer u.Mutex.Unlock()

	// Проверка на наличие уже обрабатываемого запроса от пользователя
	if u.AlreadyRunning(upd) {
		return true
	}
	// Блокируем новые запросы от пользователя
	u.SetIsRunning(true)
	return false
}

func (u *UserInfo) AlreadyRunning(upd tgbotapi.Update) bool {

	if u.IsRunning {
		warning := "Последняя операция ещё выполняется, дождитесь её завершения перед отправкой новых запросов."
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, warning)
		Bot.Send(msg)
		Logs <- Log{"bot", warning, false}
		return true
	}

	return false

}

func (u *UserInfo) SetIsRunning(v bool) {
	u.IsRunning = v
}
