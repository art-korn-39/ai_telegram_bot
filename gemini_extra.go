package main

import (
	"fmt"
)

func gen_DailyLimitOfRequestsIsOver(u *UserInfo) bool {

	if u.Requests_today_gen >= Cfg.RPD_gen {
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
