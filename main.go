package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
)

const (
	Version       = "1.3"
	ChannelChatID = -1001997602646
	ChannelURL    = "https://t.me/+6ZMACWRgFdRkNGEy"
)

var (
	db              *sql.DB
	Bot             *tgbotapi.BotAPI
	Cfg             config
	Logs            chan Log
	ListOfUsers     = map[int64]*UserInfo{}
	arrayCMD        = []string{"gemini", "kandinsky", "chatgpt"}
	delay_ChatGPT   = time.Tick(time.Second * 12 / 11) // 55 запросов в минуту
	delay_Gemini    = time.Tick(time.Second * 12 / 11) // 55 запросов в минуту
	delay_Kandinsky = time.Tick(time.Second / 3)       // 20 запросов в минуту
)

//sql
//счетчик запросов
//последняя команда

//таблица записей:
// data_time | user id | username | chatgpt | gemini | kandinskiy | request

//ограничения ChatGPT в бесплатной версии – 60 запросов в минуту
//Gemini в бесплатном тарифе действует ограничение на 60 запросов в минуту.

//ID chat (art_korn_39) = 403059287
//ID chat (art_korneev) = 609614322
//ID chat (apolo39) = 6648171361

func main() {

	defer LogPanic("", true)

	// Загрузить файл конфигурации
	LoadConfig()

	// Запустить бота
	StartBot()

	// Установить соединение с базой данных
	SQL_Connect()

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, _ := Bot.GetUpdatesChan(u)

	// В отдельно горутине обрабатываем информацию по логам
	Logs = make(chan Log, 10)
	go SaveLogs()

	// Читаем входящие запросы из канала
	for update := range updates {

		go func(upd tgbotapi.Update) {

			// Запишем panic если горутина завершилась с ошибкой
			defer LogPanic(upd.Message.Text, false)

			if upd.Message == nil {
				return
			}

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
				User = &UserInfo{}
				ListOfUsers[upd.Message.Chat.ID] = User
			}

			// Фиксируем пользователя и входящее сообщение
			Logs <- Log{upd.Message.From.UserName, upd.Message.Text, false}

			// Если предыдущий запрос ещё выполняется, то новые команды не обрабатываем
			if User.CheckUserLock(upd) {
				return
			}
			defer User.SetIsRunning(false)

			// Обработка запроса от пользователя
			var result ResultOfRequest
			if MsgIsCommand(upd.Message) {
				cmd := MsgCommand(upd.Message)
				User.LastCommand = cmd
				result = ProcessCommand(cmd, upd)
			} else {
				result = ProcessText(upd.Message.Text, User.LastCommand, upd)
			}

			// Отправка сообщения
			Bot.Send(result.Message)

			// Общий лог, пишем сюда все запросы
			Logs <- Log{result.Log_author, result.Log_message, false}

		}(update)
	}
}

func LoadConfig() {

	log.Println("Version: " + Version)

	file, err := os.OpenFile("config.txt", os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(b, &Cfg)

	log.Println("Config download complete")

}

func StartBot() {

	// Используя токен создаем новый инстанс бота
	var err error
	Bot, err = tgbotapi.NewBotAPI(Cfg.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", Bot.Self.UserName)

}

func SQL_Connect() {

	return

	// Capture connection properties.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		Cfg.DB_host, Cfg.DB_port, Cfg.DB_user, Cfg.DB_password, Cfg.DB_name)

	// Get a database handle.
	var err error
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println("Unsuccessful connection to PostgreSQL!")
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Println("Unsuccessful connection to PostgreSQL!")
		log.Fatal(pingErr)
	}

	log.Println("Successful connection to PostgreSQL")

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

// Реализация через webhook
// if _, err := bot.SetWebhook(tgbotapi.NewWebhook(WebHook)); err != nil {
// 	log.Fatalf("setting webhook %v; error: %v", WebHook, err)
// }

//bot.ListenForWebhook("/")
