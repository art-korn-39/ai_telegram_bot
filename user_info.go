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
	Username            string
	ChatID              int64  `db:"chat_id"`
	Language            string `db:"language"`
	IsRunning           bool
	Path                string
	Options             map[string]string
	Messages_ChatGPT    []openai.ChatCompletionMessage
	Messages_Gemini     []*genai.Content
	Images_Gemini       map[int]string // Удалять не забыть
	Tokens_used_gpt     int
	Requests_today_gen  int
	Requests_today_sdxl int
	Level               UserLevel
	LevelChecked        bool // Если false, то выполняем EditLevel()
	Mutex               sync.Mutex
	WG                  sync.WaitGroup
}

type UserLevel int

const (
	Basic UserLevel = iota
	Advanced
)

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
		upd.Message.Text == "/account" ||
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

	// Если пользователь сделал больше 3 операций, то без подписки не даем продолжить
	cnt, isErr := SQL_CountOfUserOperations(u)
	if isErr {
		msgText := GetText(MsgText_UnexpectedError, u.Language)
		SendMessage(u, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return false
	} else if cnt >= 3 {
		msgText := GetText(MsgText_SubscribeForUsing, u.Language)
		SendMessage(u, msgText, GetButton(btn_Subscribe, u.Language), "")
		return false
	} else { // меньше 3 операций
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

func (u *UserInfo) SetPath(v string) {
	u.Path = v
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

func (u *UserInfo) ClearTokens() {
	u.Tokens_used_gpt = 0
	u.Requests_today_gen = 0
	u.Requests_today_sdxl = 0
}

// Выполняет обновление уровня пользователя
// Если выполняется регламентно (с флагом customOperation = false),
// то после установки уровня - флаг выполненной проверки не ставится
func (u *UserInfo) EditLevel(customOperation bool) {

	if customOperation {
		if u.LevelChecked {
			return
		} else {
			// Для подстраховки сразу запишем пустой лог пользователя, чтобы сегодняшний день попал в серию
			SQL_AddLog(NewLog(u, "", Info, "first log today by user"))
		}
	}

	days, _ := SQL_UserDayStreak(u)

	//fmt.Println(days, u.Username)
	if days >= Cfg.DaysForAdvancedStatus {
		u.Level = Advanced
	} else {
		u.Level = Basic
	}

	if customOperation {
		u.LevelChecked = true
	}

}
