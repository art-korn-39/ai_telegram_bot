package main

import (
	"database/sql"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
)

const (
	Version       = "2.0.1"
	ChannelChatID = -1001997602646
	ChannelURL    = "https://t.me/+6ZMACWRgFdRkNGEy"
)

var (
	db              *sql.DB
	Bot             *tgbotapi.BotAPI
	Cfg             config
	Logs            = make(chan Log, 10)
	ListOfUsers     = map[int64]*UserInfo{}
	arrayCMD        = []string{"gemini", "kandinsky", "chatgpt"}
	admins          = []string{"Art_Korn_39", "Nik_05_04", "MnNik0"}
	UserInfoChanged = false

	delay_upd            = time.Tick(time.Millisecond * 10)
	delay_ChatGPT        = time.Tick(time.Second * 15 / 10) // 40 RPM
	delay_Gemini         = time.Tick(time.Second * 12 / 11) // 55 RPM
	delay_Kandinsky      = time.Tick(time.Second * 3)       // 20 RPM
	delay_SaveUserStates = time.Tick(time.Minute * 1)       // 1 RPM
)

func main2() {

	// // лучше использовать location = UTC для корректной работы Truncate()
	// StartOfDay := time.Date(2023, 12, 31, 2, 45, 0, 0, time.UTC).Truncate(time.Hour * 24)
	// DateString := StartOfDay.Format(time.DateTime)

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

	// Читаем входящие запросы из канала
	for update := range updates {

		<-delay_upd

		go func(upd tgbotapi.Update) {

			if upd.Message == nil {
				return
			}

			if upd.Message.Text == "" || upd.Message.From.IsBot {
				return
			}

			// Запишем panic если горутина завершилась с ошибкой
			defer FinishGorutine(nil, upd.Message.Text, false)

			// Если сообщение было больше 10 минут назад, то пропускаем
			if time.Since(upd.Message.Time()).Seconds() > 600 {
				//				Logs <- Log{upd.Message.From.UserName, "(timeout) " + upd.Message.Text, false}
				return
			}

			// Проверка подписки пользователя на канал
			if !AccessIsAllowed(upd, nil) {
				return
			}

			// Получаем данные пользователя
			User, ok := ListOfUsers[upd.Message.Chat.ID]
			if !ok {
				User = NewUserInfo(upd.Message.From, upd.Message.Chat.ID)
				ListOfUsers[upd.Message.Chat.ID] = User
			}

			// Фиксируем пользователя и входящее сообщение
			//			Logs <- Log{User.Username, upd.Message.Text, false}

			// Если предыдущий запрос ещё выполняется, то новые команды не обрабатываем
			if User.CheckUserLock(upd) {
				return
			}
			defer User.SetIsRunning(false)

			// Обработка запроса от пользователя
			var result ResultOfRequest
			// cmd это всё что начинается с "/" и 3 модели строкой
			if MsgIsCommand(upd.Message) {
				cmd := MsgCommand(upd.Message)
				//				User.LastCommand = cmd
				result = ProcessCommand(cmd, upd, User)
			} else {
				result = ProcessText(upd.Message.Text, User, upd)
			}

			// Фиксируем факт поступивших изменений
			if result.UserInfoChanged {
				UserInfoChanged = true
			}

			result.addUsernameIntoLog(User.Username)

			// Отправка сообщения
			Bot.Send(result.Message)

			// Общий лог, пишем сюда все ответы
			//			Logs <- Log{result.Log_author, result.Log_message, false}

		}(update)
	}
}
