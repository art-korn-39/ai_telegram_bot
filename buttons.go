package main

import tgbotapi "github.com/Syfaro/telegram-bot-api"

var (
	//COMMON
	buttons_Subscribe = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("✅Подписаться", ChannelURL),
		),
	)

	buttons_Models = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Gemini"),
			tgbotapi.NewKeyboardButton("Kandinsky"),
			tgbotapi.NewKeyboardButton("ChatGPT"),
		),
	)

	//KANDINSKY
	buttons_kandStyles = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Без стиля"),
			tgbotapi.NewKeyboardButton("Art"),
			tgbotapi.NewKeyboardButton("4K"),
			tgbotapi.NewKeyboardButton("Anime"),
		),
	)
	button_kandNewgen = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить текст запроса"),
			tgbotapi.NewKeyboardButton("Выбрать другой стиль"),
		),
	)

	//GEMINI
	buttons_genTypes = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Начать диалог"),
			tgbotapi.NewKeyboardButton("Отправить картинку с текстом"),
		),
	)
	buttons_genNewgen = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить текст вопроса"),
			tgbotapi.NewKeyboardButton("Загрузить новые изображения"),
			tgbotapi.NewKeyboardButton("Начать диалог"),
		),
	)
	buttons_genEndDialog = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Завершить диалог")),
	)

	//CHAT GPT
	buttons_gptTypes = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Начать диалог"),
			tgbotapi.NewKeyboardButton("Сгенерировать аудио из текста"),
		),
	)
	buttons_gptVoices = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("onyx"),
			tgbotapi.NewKeyboardButton("nova"),
			tgbotapi.NewKeyboardButton("echo"),
			tgbotapi.NewKeyboardButton("fable"),
		),
	)
	buttons_gptClearContext = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Очистить контекст"),
			tgbotapi.NewKeyboardButton("Завершить диалог"),
		),
	)
	buttons_gptSpeechNewgen = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить текст"),
			tgbotapi.NewKeyboardButton("Выбрать другой голос"),
			tgbotapi.NewKeyboardButton("Начать диалог"),
		),
	)
	buttons_gptSampleSpeech = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Audio samples", "gpt_speech_samples"),
		),
	)
)
