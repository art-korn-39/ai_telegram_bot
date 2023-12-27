package main

import (
	"slices"
	"sync"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
	openai "github.com/sashabaranov/go-openai"
)

type UserInfo struct {
	IsRunning        bool
	Model            string
	LastCommand      string
	InputText        string
	Stage            string
	Messages_ChatGPT []openai.ChatCompletionMessage
	History_Gemini   []*genai.Content
	Mutex            sync.Mutex
}

func AccessIsAllowed(upd tgbotapi.Update) bool {

	if !Cfg.CheckSubscription {
		return true
	}

	if slices.Contains(Cfg.WhiteList, upd.Message.Chat.UserName) {
		return true
	}

	result := true

	conf := tgbotapi.ChatConfigWithUser{ChatID: ChannelChatID, UserID: int(upd.Message.Chat.ID)}
	chatMember, err := Bot.GetChatMember(conf)
	if err != nil {
		Logs <- Log{"bot{GetChatMember}", err.Error(), true}
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Произошла непредвиденная ошибка, попробуйте позже.")
		Bot.Send(msg)
		result = false
	}

	if !(chatMember.IsCreator() ||
		chatMember.IsAdministrator() ||
		chatMember.IsMember()) && upd.Message.Text != "/start" {

		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Для использования бота необходимо подписаться на канал👇")
		var button = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("✅Подписаться", ChannelURL),
			),
		)
		msg.ReplyMarkup = button
		Bot.Send(msg)
		result = false
	}

	return result

}

func (u *UserInfo) CheckUserLock(upd tgbotapi.Update) (isLocking bool) {

	// Устанавливаем блокировку по объекту UserInfo на период проверки запуска и блокировки
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
