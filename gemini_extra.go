package main

import (
	"fmt"
	"slices"
)

func gen_DailyLimitOfRequestsIsOver(u *UserInfo) bool {

	if slices.Contains(Cfg.WhiteList, u.Username) {
		return false
	}

	if u.Usage.Gen >= Cfg.RPD_gen {
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
