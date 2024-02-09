package main

import tgbotapi "github.com/Syfaro/telegram-bot-api"

type Button int

const (
	btn_RemoveKeyboard Button = iota
	btn_Subscribe
	btn_Models
	btn_Languages

	btn_KandStyles
	btn_ImgNewgenFull
	btn_ImgNewgen

	btn_SDXLStyles

	btn_GenTypes
	btn_GenNewgen
	btn_GenEndDialog

	btn_GptTypes
	btn_GptVoices
	btn_GptClearContext
	btn_GptSpeechNewgen
	btn_GptSampleSpeech
	btn_GptImageNewgen
)

func GetButton(btn Button, lang string) (keyboard any) {

	switch btn {

	//COMMON
	case btn_RemoveKeyboard:
		keyboard = tgbotapi.NewRemoveKeyboard(false)

	case btn_Subscribe:
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(GetText(BtnText_Subscribe, lang), ChannelURL),
			))
	case btn_Models:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_Gemini, "")),
				tgbotapi.NewKeyboardButton(GetText(BtnText_ChatGPT, "")),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_SDXL, "")),
				tgbotapi.NewKeyboardButton(GetText(BtnText_Kandinsky, "")),
			),
		)
	case btn_Languages:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("English"),
				tgbotapi.NewKeyboardButton("Русский"),
			))

	//KANDINSKY
	case btn_KandStyles:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("No style"),
				tgbotapi.NewKeyboardButton("Art"),
				tgbotapi.NewKeyboardButton("4K"),
				tgbotapi.NewKeyboardButton("Anime"),
			))

		//SDXL
	case btn_SDXLStyles:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("No style"),
				tgbotapi.NewKeyboardButton("3D model"),
				tgbotapi.NewKeyboardButton("Photo"),
				tgbotapi.NewKeyboardButton("Neon punk"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Cinematic"),
				tgbotapi.NewKeyboardButton("Analog film"),
				tgbotapi.NewKeyboardButton("Fantasy art"),
				tgbotapi.NewKeyboardButton("Anime"),
			),
		)

	//IMAGE
	case btn_ImgNewgenFull:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_ChangeQuerryText, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_ChooseAnotherStyle, lang)),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_Upscale, lang)),
			))
	case btn_ImgNewgen:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_ChangeQuerryText, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_ChooseAnotherStyle, lang)),
			))

	//GEMINI
	case btn_GenTypes:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_StartDialog, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_SendPictureWithText, lang)),
			))
	case btn_GenNewgen:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_ChangeQuestionText, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_UploadNewImages, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_StartDialog, lang)),
			))
	case btn_GenEndDialog:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_EndDialog, lang))))

	//CHATGPT
	case btn_GptTypes:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_StartDialog, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_GenerateAudioFromText, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_SendPictureWithText, lang)),
			))
	case btn_GptVoices:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("onyx"),
				tgbotapi.NewKeyboardButton("nova"),
				tgbotapi.NewKeyboardButton("echo"),
				tgbotapi.NewKeyboardButton("fable"),
			))
	case btn_GptClearContext:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_ClearContext, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_EndDialog, lang)),
			))
	case btn_GptSpeechNewgen:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_ChangeText, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_ChooseAnotherVoice, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_StartDialog, lang)),
			))
	case btn_GptSampleSpeech:
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Audio samples", "gpt_speech_samples"),
			))
	case btn_GptImageNewgen:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GetText(BtnText_ChangeQuestionText, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_UploadNewImage, lang)),
				tgbotapi.NewKeyboardButton(GetText(BtnText_StartDialog, lang)),
			))
	}

	return keyboard

}
