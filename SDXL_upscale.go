package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const SDXL_URL_UPSCALE = "https://api.stability.ai/v1/generation/esrgan-v1-x2plus/image-to-image/upscale"

func sdxl_Image2ImageUpscale(user *UserInfo, filepath string) (result string, err error) {

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	// Копируем файл в imageWriter
	imageWriter, _ := writer.CreateFormField("image")
	imageFile, err := os.Open(filepath)
	if err != nil {
		err = errors.New("{os.Open()} " + err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}
	io.Copy(imageWriter, imageFile)

	// Параметры запроса
	sdxl_SetResolution(writer, user)
	writer.Close()

	// Отправка запроса
	payload := bytes.NewReader(data.Bytes())
	req, _ := http.NewRequest("POST", SDXL_URL_UPSCALE, payload)

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "image/png")
	req.Header.Add("Authorization", "Bearer "+Cfg.Stability_Key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New("{http.DefaultClient.Do()} " + err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}
	defer res.Body.Close()

	// Обработка ответа с ошибкой
	if res.StatusCode != 200 {

		msgText := GetText(MsgText_FailedImageUpscale, user.Language)

		type res_body struct {
			Id      string `json:"id"`
			Name    string `json:"name"`
			Message string `json:"message"`
		}
		var body res_body

		if res.StatusCode == 522 || res.StatusCode == 520 {
			err = errors.New("522 status code =/")
			msgText = GetText(MsgText_APIdead, user.Language)
		} else if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
			err = errors.New("{Decode(&body) 1} " + err.Error())
		} else {
			err = fmt.Errorf("{%s} %s", body.Name, body.Message)
		}

		return msgText, err
	}

	// Проверка на отказ из-за цензуры
	if res.Header["Finish-Reason"][0] == "CONTENT_FILTERED" {
		err = errors.New("CONTENT_FILTERED")
		msgText := GetText(MsgText_FailedImageUpscale, user.Language)
		return msgText, err
	}

	// Создание файла
	filepath, err = sdxl_CreateImageFileFromBody(user, res)
	if err != nil {
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}

	return filepath, nil

}

func sdxl_SetResolution(writer *multipart.Writer, user *UserInfo) {

	side, ok := user.Options["minSide"]
	if ok {
		writer.WriteField(side, "2048")
	} else {
		writer.WriteField("width", "2048")
	}

}

func sdxl_SetOptionMinSide(user *UserInfo, photo *tgbotapi.PhotoSize) {
	if photo.Height < photo.Width {
		user.Options["minSide"] = "width"
	} else {
		user.Options["minSide"] = "height"
	}
}
