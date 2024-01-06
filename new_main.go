package main

import (
	"slices"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// art39 : 403059287

func main() {

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

			if upd.Message.From.IsBot {
				return
			}

			// Получаем данные пользователя
			User, ok := ListOfUsers[upd.Message.Chat.ID]
			if !ok {
				User = NewUserInfo(upd.Message.From, upd.Message.Chat.ID)
				ListOfUsers[upd.Message.Chat.ID] = User
			}

			// Запишем panic если горутина завершилась с ошибкой
			defer FinishGorutine(User, upd.Message.Text, false)

			// Проверка подписки пользователя на канал
			if !AccessIsAllowed(upd, User) {
				return
			}

			// Пустой текст пропускаем только в случае загрузки картинок у gemini
			if User.Path != "gemini/type/image" && upd.Message.Text == "" {
				return
			}

			// Фиксируем пользователя и входящее сообщение
			Logs <- NewLog(User, "", Info, upd.Message.Text)

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

func HandleMessage(u *UserInfo, m *tgbotapi.Message) {

	// 1. Определяем команду
	cmd := MsgCommand(m)
	if cmd != "" {
		u.Path = cmd
		u.ClearUserData()
	}

	// 2. Формируем ответ
	switch u.Path {
	case "start":
		SendMessage(u, start(m.From.FirstName), buttons_start, "HTML")

	case "gemini":
		gen_start(u)

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

	case "chatgpt/dialog":
		gpt_dialog(u, m.Text, true)

	default:
		if slices.Contains(admins, u.Username) {
			switch cmd {
			case "info":
				SendMessage(u, GetInfo(), button_RemoveKeyboard, "")
			case "updconf":
				LoadConfig()
				SendMessage(u, "Config updated.", button_RemoveKeyboard, "")
			}
		} else {
			SendMessage(u, "Не выбрана нейросеть для обработки запроса.", buttons_start, "")
		}
	}

}
