package main

import "fmt"

func gen_DailyLimitOfRequestsIsOver(u *UserInfo) bool {

	if u.Requests_today_gen >= Cfg.RPD_gen {
		duration := GetDurationToNextDay()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) - hours*60
		msgText := fmt.Sprintf("Превышен дневной лимит запросов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.", hours, minutes)
		SendMessage(u, msgText, buttons_Models, "")
		u.Path = "start"
		return true
	}

	return false

}
