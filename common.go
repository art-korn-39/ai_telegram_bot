package main

import (
	"fmt"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func start(user *UserInfo, message *tgbotapi.Message) {

	name := message.From.FirstName
	if name == "" {
		name = user.Username
	}

	msgtxt := fmt.Sprintf(GetText(MsgText_Start, user.Language), name)
	SendMessage(user, msgtxt, GetButton(btn_Models, ""), "HTML")

	user.Path = ""

}

func account(user *UserInfo) {

	// пример +++
	sample :=
		`
👤 ID Пользователя: <b>%d</b>
⭐️ Уровень: <b>%s</b>
✌️ Посещений подряд (дней): <b>%d</b>
✅ Дата первого использования: <b>%s</b>
----------------------------------------------
Дневные лимиты:
🃏 Gemini 1.5 запросы: <b>%d</b> (осталось <b>%d</b>)
🚀 Gemini запросы: <b>%d</b> (осталось <b>%d</b>)
🤖 ChatGPT токены: <b>%d</b> (осталось <b>%d</b>)
🗿 Kandinsky: <b>без ограничений</b>
🏔 Stable Diffusion: <b>%d</b> (осталось <b>%d</b>)
🎭 Face Swap: <b>%d</b> (осталось <b>%d</b>)
----------------------------------------------                
		
<i>Лимиты обновятся через : %d ч. %d мин.</i>
			
Регулярные пользователи бота (%d дней подряд и более) получают <b>%s</b> уровень, на котором доступно 
по <b>%d</b> генераций в Stable Diffusion и Face Swap + <b>%d</b> запросов Gemini 1.5 в сутки 🔥`

	// пример ---

	sample = GetText(MsgText_Account, user.Language)

	duration := GetDurationToNextDay()
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - hours*60
	mapWithStreak, _ := SQL_UserDayStreak(user)
	DayStreak := mapWithStreak[user.ChatID]
	FirstDate, _ := SQL_GetFirstDate(user)

	msgText := fmt.Sprintf(sample,
		user.ChatID,                                  // ID Пользователя
		GetLevelName(user.Level, user.Language),      // Уровень
		DayStreak,                                    // Посещений подряд (дней)
		FirstDate.Format(time.DateOnly),              // Дата первого использования
		Get_RPD_gen15(user),                          // Gemini 1.5 на день у пользователя
		max(Get_RPD_gen15(user)-user.Usage.Gen15, 0), // Gemini 1.5 остаток
		Cfg.RPD_gen10,                                // Gemini 1.0 на день у пользователя
		max(Cfg.RPD_gen10-user.Usage.Gen10, 0),       // Gemini 1.0 остаток
		Get_TPD_gpt(user),                            // ChatGPT на день у пользователя
		max(Get_TPD_gpt(user)-user.Usage.GPT, 0),     // ChatGPT остаток
		Get_RPD_sdxl(user),                           // Stable Diffusion на день у пользователя
		max(Get_RPD_sdxl(user)-user.Usage.SDXL, 0),   // Stable Diffusion остаток
		Get_RPD_fs(user),                             // Face Swap на день у пользователя
		max(Get_RPD_fs(user)-user.Usage.FS, 0),       // Face Swap остаток
		hours, minutes,                               // time to refresh
		Cfg.DaysForAdvancedStatus,             // дней для продвинутого уровня
		GetLevelName(Advanced, user.Language), // Уровень строкой
		Cfg.RPD_advanced_sdxl,                 // Stable Diffusion продвинутый
		//Cfg.TPD_advanced_gpt,                  // ChatGPT продвинутый
		Cfg.RPD_advanced_gen15, // Gen 1.5 продвинутый
	)

	SendMessage(user, msgText, GetButton(btn_Models, ""), "HTML")

	user.Path = ""

}

func language_start(user *UserInfo) {

	msgText := "Select language / Выберите язык"
	SendMessage(user, msgText, GetButton(btn_Languages, ""), "")

	user.Path = "language/type"

}

// После выбора языка
func language_type(user *UserInfo, text string) {

	switch text {
	case "English":
		user.Language = "en"
	case "Русский":
		user.Language = "ru"
	default:
		msgText := GetText(MsgText_SelectOption, user.Language)
		SendMessage(user, msgText, nil, "")
		return
	}

	msgText := GetText(MsgText_LanguageChanged, user.Language)
	SendMessage(user, msgText, GetButton(btn_Models, user.Language), "")

	user.Path = "start"

}
