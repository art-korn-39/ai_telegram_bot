package main

import (
	"os"
	"slices"
	"sync"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
	openai "github.com/sashabaranov/go-openai"
)

type UserInfo struct {
	Username         string
	ChatID           int64
	IsRunning        bool
	Path             string
	Options          map[string]string
	Messages_ChatGPT []openai.ChatCompletionMessage
	Messages_Gemini  []*genai.Content
	Images_Gemini    map[int]string // удалять не забыть
	Tokens_used_gpt  int
	Mutex            sync.Mutex
	WG               sync.WaitGroup
}

func NewUserInfo(u *tgbotapi.User, id int64) *UserInfo {
	username := u.UserName
	if username == "" {
		username = u.FirstName + "_" + u.LastName
	}
	return &UserInfo{Username: username, ChatID: id}
}

func (u *UserInfo) ClearUserData() {
	u.Options = map[string]string{}
	u.Messages_ChatGPT = []openai.ChatCompletionMessage{}
	u.Messages_Gemini = []*genai.Content{}
	u.DeleteImages()
}

func (u *UserInfo) DeleteImages() {
	for _, v := range u.Images_Gemini {
		os.Remove(v)
	}
	u.Images_Gemini = map[int]string{}
}

func AccessIsAllowed(upd tgbotapi.Update, u *UserInfo) bool {

	if !Cfg.CheckSubscription {
		return true
	}

	if slices.Contains(Cfg.WhiteList, u.Username) {
		return true
	}

	result := true

	conf := tgbotapi.ChatConfigWithUser{ChatID: ChannelChatID, UserID: int(u.ChatID)}
	chatMember, err := Bot.GetChatMember(conf)
	if err != nil {
		Logs <- NewLog(u, "bot", Error, "{GetChatMember} "+err.Error())
		msgText := "Произошла непредвиденная ошибка, попробуйте позже."
		SendMessage(u, msgText, nil, "")
		result = false
	}

	if !(chatMember.IsCreator() ||
		chatMember.IsAdministrator() ||
		chatMember.IsMember()) && upd.Message.Text != "/start" {
		msgText := "Для использования бота необходимо подписаться на канал👇"
		SendMessage(u, msgText, buttons_Subscribe, "")
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
	if u.IsRunning && !u.ImagesLoading(upd) {
		msgText := "Последняя операция ещё выполняется, дождитесь её завершения перед отправкой новых запросов."
		SendMessage(u, msgText, nil, "")
		Logs <- NewLog(u, "bot", Info, msgText)
		return true
	}

	// Блокируем новые запросы от пользователя
	u.SetIsRunning(true)
	return false
}

func (u *UserInfo) SetIsRunning(v bool) {
	u.IsRunning = v
}

func (u *UserInfo) ImagesLoading(upd tgbotapi.Update) bool {
	if u.Path == "gemini/type/image" && upd.Message.Photo != nil {
		return true
	}
	return false
}
