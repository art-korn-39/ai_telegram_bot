package main

import (
	"strconv"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// https://replicate.com/yan-ops/face_swap?input=http
// https://icons8.ru/swapper ($99 в год)

// После команды "/faceswap" или при вводе текста = "faceswap"
func fs_start(user *UserInfo) {

	if fs_DailyLimitOfRequestsIsOver(user) {
		return
	}

	msgText := fs_WelcomeTextMessage(user)
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "HTML")

	msgText = GetText(MsgText_FSimage1, user.Language)
	SendMessage(user, msgText, nil, "")

	user.Path = "faceswap/image1"

}

// После отправки картинок пользователем
func fs_image(user *UserInfo, message *tgbotapi.Message, i int) {

	if fs_DailyLimitOfRequestsIsOver(user) {
		return
	}

	// Проверяем наличие картинок в сообщении
	if message.Photo == nil {
		msgText := GetText(MsgText_UploadImage, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	SendMessage(user, GetText(MsgText_LoadingImage, user.Language), nil, "")

	photos := *message.Photo

	// Получаем URL картинки
	fileURL, err := Bot.GetFileDirectURL(photos[len(photos)-1].FileID)
	if err != nil {
		Logs <- NewLog(user, "fs{img1}", Error, err.Error())
		msgText := GetText(MsgText_FailedLoadImages, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	user.Options["image"+strconv.Itoa(i)] = fileURL

	if i == 1 {
		// Просим написать запрос к картинке
		msgText := GetText(MsgText_FSimage2, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

		user.Path = "faceswap/image2"
	} else {

		msgText := GetText(MsgText_ImageGenerationStarted2, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

		// Отправка запроса
		result, msg, err := fs_SendRequest(user)
		if err != nil {
			Logs <- NewLog(user, "faceswap", Error, err.Error())
			SendMessage(user, GetText(msg, user.Language), GetButton(btn_FSNewgen, user.Language), "")
		} else {
			// Отправка полученного изображения пользователю
			err := SendPhotoMessage(user, result, "", GetButton(btn_FSNewgenFull, user.Language))
			if err != nil {
				Logs <- NewLog(user, "faceswap", Error, "{ImgSend} "+err.Error())
				SendMessage(user, GetText(MsgText_ErrorWhileSendingPicture, user.Language), GetButton(btn_FSNewgen, user.Language), "")
			} else {
				user.Requests_today_fs++
				user.Usage.FS++
				Operation := SQL_NewOperation(user, "faceswap", "", "")
				SQL_AddOperation(Operation)
				user.Options["image"] = result
			}
		}

		user.Path = "faceswap/newgen"
	}

}

// После получения результата генерации
func fs_newgen(user *UserInfo, text string) {

	if fs_DailyLimitOfRequestsIsOver(user) {
		return
	}

	switch text {

	// NEW IMAGES
	case GetText(BtnText_UploadNewImages, user.Language):
		SendMessage(user, GetText(MsgText_FSimage1, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "faceswap/image1"

	// UPSCALE
	case GetText(BtnText_Upscale, user.Language):
		err := fs_PrepareImageToUpscale(user)
		if err != nil {
			SendMessage(user, GetText(MsgText_UnexpectedError, user.Language), nil, "")
		} else {
			sdxl_upscale(user, btn_FSNewgen)
		}
		user.Path = "faceswap/newgen"

	default:
		SendMessage(user, GetText(MsgText_UnknownCommand, user.Language), nil, "")
	}

}
