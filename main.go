package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
)

const (
	Version       = "1.6"
	ChannelChatID = -1001997602646
	ChannelURL    = "https://t.me/+6ZMACWRgFdRkNGEy"
)

var (
	db          *sql.DB
	Bot         *tgbotapi.BotAPI
	Cfg         config
	Logs        chan Log
	ListOfUsers = map[int64]*UserInfo{}
	arrayCMD    = []string{"gemini", "kandinsky", "chatgpt"}

	delay_upd       = time.Tick(time.Millisecond * 10)
	delay_ChatGPT   = time.Tick(time.Second * 3)       // 20 RPM
	delay_Gemini    = time.Tick(time.Second * 12 / 11) // 55 RPM
	delay_Kandinsky = time.Tick(time.Second * 3)       // 20 RPM
	delay_stat      = time.Tick(time.Minute * 10)

	counter_chatgpt   = 0
	counter_gemini    = 0
	counter_kandinsky = 0
)

// по юзерам с пустым username записать детализацию

//ограничения ChatGPT в бесплатной версии – 3 запросов в минуту (200 в день?)
//зарегистрировать несколько API?
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

	// Пишем счетчики в логи, на всякий случай + пока БД не готова
	go SaveStatistics()

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
			defer LogPanic(upd.Message.Text, false)

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
			// cmd это всё что начинается с "/" и 3 модели строкой
			if MsgIsCommand(upd.Message) {
				cmd := MsgCommand(upd.Message)
				User.LastCommand = cmd
				result = ProcessCommand(cmd, upd, User)
			} else {
				result = ProcessText(upd.Message.Text, User, upd)
			}

			result.addUsernameIntoLog(upd.Message.From.UserName)

			// Отправка сообщения
			Bot.Send(result.Message)

			// Общий лог, пишем сюда все запросы
			Logs <- Log{result.Log_author, result.Log_message, false}

		}(update)
	}
}

func SaveStatistics() {
	for {
		<-delay_stat
		text := fmt.Sprintf("Gemini: %d ChatGPT: %d Kandinsky: %d",
			counter_gemini, counter_chatgpt, counter_kandinsky)
		log.Println(text)
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

// Реализация через webhook
// if _, err := bot.SetWebhook(tgbotapi.NewWebhook(WebHook)); err != nil {
// 	log.Fatalf("setting webhook %v; error: %v", WebHook, err)
// }

//bot.ListenForWebhook("/")

// /start - start msg
// /gemini - "введите вопрос"
//    text - result
// /chatgpt - "введите вопрос"
//    text - result
// /kandinsky - "введите запрос"
//    text - "выберите стиль изображения"
//	     style - result
