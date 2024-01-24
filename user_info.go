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
	Username           string
	ChatID             int64
	Language           string
	IsRunning          bool
	Path               string
	Options            map[string]string
	Messages_ChatGPT   []openai.ChatCompletionMessage
	Messages_Gemini    []*genai.Content
	Images_Gemini      map[int]string // удалять не забыть
	Tokens_used_gpt    int
	Requests_today_gen int
	Mutex              sync.Mutex
	WG                 sync.WaitGroup
}

func NewUserInfo(m *tgbotapi.Message) *UserInfo {

	u := m.From
	id := m.Chat.ID

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

	if upd.Message.Text == "/start" ||
		upd.Message.Text == "/language" ||
		u.Path == "language/type" {
		return true
	}

	conf := tgbotapi.ChatConfigWithUser{ChatID: ChannelChatID, UserID: int(u.ChatID)}
	chatMember, err := Bot.GetChatMember(conf)
	if err != nil {
		Logs <- NewLog(u, "bot", Error, "{GetChatMember} "+err.Error())
		msgText := GetText(MsgText_UnexpectedError, u.Language)
		SendMessage(u, msgText, nil, "")
		return false
	}

	if chatMember.IsCreator() ||
		chatMember.IsAdministrator() ||
		chatMember.IsMember() {
		return true
	}

	// Если пользователь сделал больше 2 операций, то без подписки не даем продолжить
	cnt, isErr := SQL_CountOfUserOperations(u)
	if isErr {
		msgText := GetText(MsgText_UnexpectedError, u.Language)
		SendMessage(u, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return false
	} else if cnt >= 2 {
		msgText := GetText(MsgText_SubscribeForUsing, u.Language)
		SendMessage(u, msgText, GetButton(btn_Subscribe, u.Language), "")
		return false
	} else { // меньше 2 операций
		return true
	}

}

func (u *UserInfo) CheckUserLock(upd tgbotapi.Update) (isLocking bool) {

	// Устанавливаем блокировку по объекту UserInfo на период проверки запуска и блокировки
	u.Mutex.Lock()
	// Разблокировку ставим через defer, чтобы она не стала вечной, если в методе произойдёт panic
	defer u.Mutex.Unlock()

	// Проверка на наличие уже обрабатываемого запроса от пользователя
	if u.IsRunning && !u.ImagesLoading(upd) {
		msgText := GetText(MsgText_LastOperationInProgress, u.Language)
		SendMessage(u, msgText, nil, "")
		return true
	}

	// Блокируем новые запросы от пользователя
	u.SetIsRunning(true)
	return false
}

func (u *UserInfo) SetIsRunning(v bool) {
	u.IsRunning = v
}

func (u *UserInfo) FillLanguage(lang string) {
	if u.Language == "" {
		u.Language = lang
	}
}

func (u *UserInfo) ImagesLoading(upd tgbotapi.Update) bool {
	if u.Path == "gemini/type/image" && upd.Message.Photo != nil {
		return true
	}
	return false
}
