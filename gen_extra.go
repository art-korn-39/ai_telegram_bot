package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/google/generative-ai-go/genai"
)

func gen_DailyLimitOfRequestsIsOver(u *UserInfo, version string) bool {

	if slices.Contains(Cfg.WhiteList, u.Username) {
		return false
	}

	limitIsOver := false
	if version == gen10 {
		limitIsOver = u.Usage.Gen10 >= Cfg.RPD_gen10
	} else {
		limitIsOver = u.Usage.Gen15 >= Cfg.RPD_gen15
	}

	if limitIsOver {
		duration := GetDurationToNextDay()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) - hours*60
		msgText := fmt.Sprintf(GetText(MsgText_DailyRequestLimitExceeded, u.Language), hours, minutes)
		SendMessage(u, msgText, GetButton(btn_Models, u.Language), "")
		u.Path = "start"
		return true
	}

	return false

}

func gen15_MIME_TYPES() map[string]string {

	m := map[string]string{}

	//TEXT
	m[".txt"] = "text/plain"
	m[".html"] = "text/html"
	m[".css"] = "text/css"
	m[".xml"] = "text/xml"
	m[".js"] = "text/javascript"
	m[".py"] = "text/x-python" // "application/x-python-code"
	m[".json"] = "application/json"

	//AUDIO
	m[".mp3"] = "audio/mp3"
	m[".ogg"] = "audio/ogg"
	m[".oga"] = "audio/oga"
	m[".wav"] = "audio/wav"

	//IMAGES
	m[".jpeg"] = "image/jpeg"
	m[".png"] = "image/png"
	m[".webp"] = "image/webp"
	m[".heic"] = "image/heic"
	m[".heif"] = "image/heif"

	//VIDEO
	m[".mp4"] = "video/mp4"
	m[".mpeg"] = "video/mpeg"
	m[".mov"] = "video/mov"
	m[".avi"] = "video/avi"
	m[".mpg"] = "video/mpg"
	m[".wmv"] = "video/wmv"
	m[".3gp"] = "video/3gpp"

	return m
}

func gen15_GetFileID(message *tgbotapi.Message) (fileID string) {

	if message.Photo != nil {
		photos := *message.Photo
		fileID = photos[len(photos)-1].FileID
	} else if message.Video != nil {
		fileID = message.Video.FileID
	} else if message.Audio != nil {
		fileID = message.Audio.FileID
	} else if message.Voice != nil {
		fileID = message.Voice.FileID
	} else if message.Document != nil {
		fileID = message.Document.FileID
	}

	return fileID

}

func gen15_GetUploadFileOptions(filename string) *genai.UploadFileOptions {

	ext := filepath.Ext(filename)
	MIME, ok := gen15_MIME_TYPES()[ext]

	if ok {
		return &genai.UploadFileOptions{MIMEType: MIME}
	} else {
		return nil
	}

}

func gen15_ProcessErrorsOfResponse(user *UserInfo, errIn error) (msgText string) {

	errString := errIn.Error()

	if strings.Contains(errString, "googleapi: Error 400: Unsupported MIME type:") {
		msgText = GetText(MsgText_BadRequest6, user.Language)
	} else {
		msgText = GetText(MsgText_BadRequest5, user.Language)
	}

	return

}

func gen15_UploadFileToCloudStorage(filename string) (*genai.File, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file, err := gen_client.UploadFile(gen_ctx, "", f, gen15_GetUploadFileOptions(filename))
	if err != nil {
		return nil, err
	}

	return file, nil

}
