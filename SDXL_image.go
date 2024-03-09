package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"slices"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"golang.org/x/image/draw"
)

// После отправки картинки пользователем
func sdxl_image(user *UserInfo, message *tgbotapi.Message) {

	if sdxl_DailyLimitOfRequestsIsOver(user, btn_Models) {
		return
	}

	// Проверяем наличие картинок в сообщении
	if message.Photo == nil && message.Document == nil {
		msgText := GetText(MsgText_UploadImage, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	// Файлы и фотографии (со сжатием) приводим к одному формату
	Photo := tgbotapi.PhotoSize{}
	if message.Photo != nil {
		photos := *message.Photo
		Photo = photos[len(photos)-1]
	} else {
		if slices.Contains(GetImageMimeTypes(), message.Document.MimeType) {
			Photo.FileID = message.Document.FileID
			Photo.FileSize = message.Document.FileSize
		} else {
			msgText := GetText(MsgText_AvailiableImageFormats, user.Language)
			SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
			return
		}
	}

	msgText := GetText(MsgText_LoadingImage, user.Language)
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

	// Сохраняем картинку в файловую систему
	name := fmt.Sprintf("img_%d_0", user.ChatID)
	filename, err := DownloadFile(Photo.FileID, name)
	if err != nil {
		Logs <- NewLog(user, "sdxl{img}", Error, err.Error())
		msgText := GetText(MsgText_FailedLoadImages, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	// Уменьшаем разрешение если требуется
	filename, err = AdaptImageResolution(&Photo, filename, user)
	if err != nil {
		Logs <- NewLog(user, "sdxl{img2}", Error, err.Error())
		msgText := GetText(MsgText_FailedImageUpscale, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	user.Options["image"] = filename
	if Photo.Height < Photo.Width {
		user.Options["minSide"] = "width"
	} else {
		user.Options["minSide"] = "height"
	}

	sdxl_upscale(user, btn_ImgNewgen)

	msgText = GetText(MsgText_SelectOption, user.Language)
	SendMessage(user, msgText, GetButton(btn_SDXLTypes, user.Language), "")

	user.Path = "sdxl/type"
}

// Приводит изображение к 1048576 пикселей, пропорционально сжимая
// Если исходное изображение <= 1048576 пикселей, возвращается ссылка на исходник
func AdaptImageResolution(ph *tgbotapi.PhotoSize, filename string, user *UserInfo) (newFilename string, err error) {

	// 4194304 (max out)
	maxPixels := float64(1048576)

	reader, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}

	// Если это файл, то ширина и высота будут пустые
	if ph.Width == 0 && ph.Height == 0 {
		ph.Width = img.Bounds().Max.X
		ph.Height = img.Bounds().Max.Y
	}

	// Проверяем общее количество пикселей
	if ph.Width*ph.Height <= int(maxPixels) {
		return filename, nil
	}

	// Вычисляем ширину и высоту для итогового изображения
	ratio := float64(ph.Width) / float64(ph.Height)
	h2 := math.Sqrt(maxPixels / ratio)
	w2 := maxPixels / h2

	// Создаём пустое изображение с необходимым размером
	newImg := image.NewRGBA(image.Rect(0, 0, int(w2), int(h2)))

	// Изменение размера
	draw.CatmullRom.Scale(newImg, newImg.Bounds(),
		img, img.Bounds(),
		draw.Over, nil)

	// Подтоговка пути к файлу
	outFile := getFilepathForImage(user.ChatID, "png")

	// Создание пустого файла
	file, err := os.Create(outFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Помещение в файл изображения
	if err := png.Encode(file, newImg); err != nil {
		return "", err
	}

	return outFile, nil

}
