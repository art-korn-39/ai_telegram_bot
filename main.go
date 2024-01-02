package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
)

const (
	Version       = "1.8"
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
	delay_ChatGPT        = time.Tick(time.Second * 5)       // 12 RPM
	delay_Gemini         = time.Tick(time.Second * 12 / 11) // 55 RPM
	delay_Kandinsky      = time.Tick(time.Second * 3)       // 20 RPM
	delay_SaveUserStates = time.Tick(time.Second * 30)      // 1 RPM
)

func main() {

	defer FinishGorutine("", true)

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
			defer FinishGorutine(upd.Message.Text, false)

			// Если сообщение было больше 10 минут назад, то пропускаем
			if time.Since(upd.Message.Time()).Seconds() > 600 {
				Logs <- Log{upd.Message.From.UserName, "(timeout) " + upd.Message.Text, false}
				return
			}

			// Проверка подписки пользователя на канал
			if !AccessIsAllowed(upd) {
				return
			}

			// Получаем данные пользователя
			User, ok := ListOfUsers[upd.Message.Chat.ID]
			if !ok {
				User = NewUserInfo(upd.Message.From, upd.Message.Chat.ID)
				ListOfUsers[upd.Message.Chat.ID] = User
			}

			// Фиксируем пользователя и входящее сообщение
			Logs <- Log{User.Username, upd.Message.Text, false}

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
				User.LastCommand = cmd
				result = ProcessCommand(cmd, upd, User)
			} else {
				result = ProcessText(upd.Message.Text, User, upd)
			}

			UserInfoChanged = true // фиксируем факт поступивших изменений

			result.addUsernameIntoLog(User.Username)

			// Отправка сообщения
			Bot.Send(result.Message)

			// Общий лог, пишем сюда все запросы
			Logs <- Log{result.Log_author, result.Log_message, false}

		}(update)
	}
}

func StartBot() {

	var err error
	Bot, err = tgbotapi.NewBotAPI(Cfg.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", Bot.Self.UserName)

}

func FinishGorutine(inputtext string, main bool) {

	timeNow := time.Now().UTC().Add(3 * time.Hour).Format(time.DateTime)
	if r := recover(); r != nil {
		text := "Inputtext: " + inputtext + "\n" + "Error: " + fmt.Sprint(r)
		fmt.Println(timeNow+" Panic in gorutine:", text)
		WriteIntoFile(timeNow, Ternary(main, "main", "gorutine"), text)
	}

}
