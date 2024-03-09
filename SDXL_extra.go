package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
)

var (
	SDXL_Styles = map[string]string{
		"No style":    "",
		"3D model":    "3d-model",
		"Photo":       "photographic",
		"Neon punk":   "neon-punk",
		"Cinematic":   "cinematic",
		"Analog film": "analog-film",
		"Fantasy art": "fantasy-art",
		"Anime":       "anime",
	}
)

func sdxl_WelcomeTextMessage(u *UserInfo) string {

	duration := GetDurationToNextDay()
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - hours*60

	return fmt.Sprintf(GetText(MsgText_SDXLinfo, u.Language),
		max(Get_RPD_sdxl(u)-u.Requests_today_sdxl, 0),
		hours,
		minutes)
}

func sdxl_DailyLimitOfRequestsIsOver(u *UserInfo, btn Button) bool {

	if slices.Contains(Cfg.WhiteList, u.Username) {
		return false
	}

	if u.Requests_today_sdxl >= Get_RPD_sdxl(u) {
		duration := GetDurationToNextDay()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) - hours*60
		msgText := fmt.Sprintf(GetText(MsgText_DailyRequestLimitExceeded, u.Language), hours, minutes)
		if btn == 0 {
			SendMessage(u, msgText, nil, "")
		} else {
			SendMessage(u, msgText, GetButton(btn_Models, u.Language), "")
			u.Path = "start"
		}
		return true
	}

	return false

}

func sdxl_CreateImageFileFromBase64(i int, image TextToImage, user *UserInfo) (string, error) {

	outFile := getFilepathForImage(user.ChatID, "png")
	file, err := os.Create(outFile)
	if err != nil {
		err = errors.New("{os.Create()} " + err.Error())
		return "", err
	}

	imageBytes, err := base64.StdEncoding.DecodeString(image.Base64)
	if err != nil {
		err = errors.New("{DecodeString()} " + err.Error())
		return "", err
	}

	if _, err := file.Write(imageBytes); err != nil {
		err = errors.New("{file.Write()} " + err.Error())
		return "", err
	}

	if err := file.Close(); err != nil {
		err = errors.New("{file.Close()} " + err.Error())
		return "", err
	}

	return outFile, nil

}

func sdxl_CreateImageFileFromBody(user *UserInfo, res *http.Response) (string, error) {

	// Write the response to a file
	outFile := getFilepathForImage(user.ChatID, "png")
	file, err := os.Create(outFile)
	if err != nil {
		err = errors.New("{os.Create()} " + err.Error())
		return "", err
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		err = errors.New("{io.Copy} " + err.Error())
		return "", err
	}

	return outFile, nil

}
