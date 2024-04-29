package main

import (
	"slices"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

// art39 : 403059287
// apolo39 : 6648171361
// https://elevenlabs.io/voice-lab

const (
	Version       = "2.5.6"
	ChannelChatID = -1001997602646
	ChannelURL    = "https://t.me/+6ZMACWRgFdRkNGEy"
)

var (
	db              *sqlx.DB
	Bot             *tgbotapi.BotAPI
	Cfg             config
	Logs            = make(chan Log, 10)
	ListOfUsers     = map[int64]*UserInfo{}
	recoveryChatID  = []int64{}
	UserInfoChanged = false
)

func main() {

	defer FinishGorutine(nil, "", true)

	// Загрузить файл конфигурации
	LoadConfig()

	// Запустить бота
	StartBot()

	// Установить соединение с базой данных
	SQL_Connect()

	// Загрузить текущие состояния по пользователям
	SQL_LoadUserStates()

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, _ := Bot.GetUpdatesChan(u)

	// В отдельно горутине обрабатываем информацию по логам
	go SaveLogs()

	// При наличии изменений - регулярно обновляем инфо по юзерам в БД
	go SaveUserStates()

	// Каждый день в 00:00 по Мск
	go EveryDayAt2400()

	// Обновление model_id каждые 20 минут
	go Kandinsky_CheckModelID()

	// Читаем входящие запросы из канала
	for update := range updates {

		<-delay_upd

		go func(upd tgbotapi.Update) {

			// Пустые сообщения и сообщения от ботов - пропускаем
			// Если это CallbackQuery, то помещаем его в Message
			if !ValidMessage(&upd) {
				return
			}

			// Получаем данные пользователя
			User, ok := ListOfUsers[upd.Message.Chat.ID]
			if !ok {
				User = NewUserInfo(upd.Message)
				ListOfUsers[upd.Message.Chat.ID] = User
			}

			// Запишем panic если горутина завершилась с ошибкой
			defer FinishGorutine(User, upd.Message.Text, false)

			// Если язык пустой, то заполним из данных сообщения
			User.FillLanguage(upd.Message.From.LanguageCode)

			// В случае восстановления после простоя - старые сообщения не обрабатываем
			if IsRecovery(upd, User) {
				return
			}

			// Фиксируем пользователя и входящее сообщение
			Logs <- NewLog(User, "", Info, upd.Message.Text)

			// Обновим уровень пользователя
			User.EditLevelManualy()

			// Проверка подписки пользователя на канал
			if !AccessIsAllowed(upd, User) {
				return
			}

			// Пустой текст пропускаем только в случае загрузки картинок
			if !MessageWithData(upd.Message.Text, User.Path) {
				SendMessage(User, GetText(MsgText_WrongDataType, User.Language), nil, "")
				return
			}

			// Если предыдущий запрос ещё выполняется, то новые команды не обрабатываем
			if User.CheckUserLock(upd) {
				return
			}
			defer User.SetIsRunning(false)

			// Обработка запроса от пользователя
			HandleMessage(User, upd.Message)

			// Пока что всегда ставим флаг
			UserInfoChanged = true

		}(update)
	}
}

func ValidMessage(upd *tgbotapi.Update) bool {

	if upd.Message == nil && upd.CallbackQuery == nil {
		return false
	}

	//if upd.Message != nil {
	//if upd.Message.From.IsBot {
	//return false
	//}
	//}

	// Переносим данные из CallbackQuery в Message
	if upd.Message == nil && upd.CallbackQuery != nil {
		Message := tgbotapi.Message{
			From: upd.CallbackQuery.From,
			Text: upd.CallbackQuery.Data,
			Chat: upd.CallbackQuery.Message.Chat,
			Date: int(time.Now().Unix()),
		}
		upd.Message = &Message
	}

	return true

}

func MessageWithData(text, path string) bool {

	if text != "" {
		return true
	} else if path == "gemini/type/image" {
		return true
	} else if path == "chatgpt/type/image" {
		return true
	} else if path == "sdxl/type/image" {
		return true
	} else if path == "faceswap/image1" {
		return true
	} else if path == "faceswap/image2" {
		return true

	}

	return false

}

func IsRecovery(upd tgbotapi.Update, user *UserInfo) bool {

	// Если сообщение было больше 100 секунд назад, то один раз отвечаем
	if time.Since(upd.Message.Time()).Seconds() > 100 {
		if !slices.Contains(recoveryChatID, user.ChatID) {
			recoveryChatID = append(recoveryChatID, user.ChatID)
			if Cfg.Debug {
				SendMessage(user, GetText(MsgText_AfterRecoveryDebug, user.Language), nil, "")
			} else {
				SendMessage(user, GetText(MsgText_AfterRecoveryProd, user.Language), nil, "")
			}
		}
		return true
	}

	return false

}

func HandleMessage(u *UserInfo, m *tgbotapi.Message) {

	// 1. Определяем команду
	cmd := MsgCommand(m)

	// 2. Если команда, то устанавливаем базовый путь и очищаем временные данные пользователя
	if cmd != "" {
		u.Path = cmd
		u.ClearUserData()
	}

	// 3. Формируем ответ
	switch u.Path {
	case "start":
		start(u, m)

	case "account":
		account(u)

	case "language":
		language_start(u)

	case "language/type":
		language_type(u, m.Text)

	case "gemini":

		if slices.Contains(Cfg.Admins, u.Username) {
			gen_start(u)
		} else {
			gen_rip(u)
		}

	case "gemini/type":
		gen_type(u, m.Text)

	case "gemini/type/dialog":
		gen_dialog(u, m.Text)

	case "gemini/type/image":
		gen_image(u, m)

	case "gemini/type/image/text":
		gen_imgtext(u, m.Text)

	case "gemini/type/image/text/newgen":
		gen_imgtext_newgen(u, m.Text)

	case "kandinsky":
		kand_start(u)

	case "kandinsky/text":
		kand_text(u, m.Text)

	case "kandinsky/text/style":
		kand_style(u, m.Text)

	case "kandinsky/text/style/newgen":
		kand_newgen(u, m.Text)

	case "chatgpt":
		gpt_start(u)

	case "chatgpt/type":
		gpt_type(u, m.Text)

	case "chatgpt/type/dialog":
		gpt_dialog(u, m.Text)

	case "chatgpt/type/speech_text":
		gpt_speech_text(u, m.Text)

	case "chatgpt/type/speech_text/voice":
		gpt_speech_voice(u, m.Text)

	case "chatgpt/type/speech_text/voice/newgen":
		gpt_speech_newgen(u, m.Text)

	case "chatgpt/type/image":
		gpt_image(u, m)

	case "chatgpt/type/image/text":
		gpt_imgtext(u, m.Text)

	case "chatgpt/type/image/text/newgen":
		gpt_imgtext_newgen(u, m.Text)

	case "sdxl":
		sdxl_start(u)

	case "sdxl/type":
		sdxl_type(u, m.Text)

	case "sdxl/type/text":
		sdxl_text(u, m.Text)

	case "sdxl/type/text/style":
		sdxl_style(u, m.Text)

	case "sdxl/type/text/style/newgen":
		sdxl_newgen(u, m.Text)

	case "sdxl/type/image":
		sdxl_image(u, m)

	case "faceswap":
		fs_start(u)

	case "faceswap/image1":
		fs_image(u, m, 1)

	case "faceswap/image2":
		fs_image(u, m, 2)

	case "faceswap/newgen":
		fs_newgen(u, m.Text)

	default:
		if slices.Contains(Cfg.Admins, u.Username) {
			HandleAdminCommand(u, cmd)
		} else {
			SendMessage(u, GetText(MsgText_UnknownCommand, u.Language), GetButton(btn_Models, ""), "")
		}
	}

}
