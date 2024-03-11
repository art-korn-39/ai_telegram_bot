package main

import (
	"fmt"
	"slices"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func fs_WelcomeTextMessage(u *UserInfo) string {

	duration := GetDurationToNextDay()
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - hours*60

	return fmt.Sprintf(GetText(MsgText_FSinfo, u.Language),
		max(Get_RPD_fs(u)-u.Requests_today_fs, 0),
		hours,
		minutes)
}

func fs_DailyLimitOfRequestsIsOver(u *UserInfo) bool {

	if slices.Contains(Cfg.WhiteList, u.Username) {
		return false
	}

	if u.Requests_today_fs >= Get_RPD_fs(u) {
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

func fs_PrepareImageToUpscale(user *UserInfo) (err error) {

	filename := user.Options["image"]

	photo := tgbotapi.PhotoSize{}
	filename, err = AdaptImageResolution(&photo, filename, user)
	if err == nil {
		user.Options["image"] = filename
		sdxl_SetOptionMinSide(user, &photo)
	}

	return

}
