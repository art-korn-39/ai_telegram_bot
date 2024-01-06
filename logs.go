package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// писать логи в БД
// date | chat_id | author (username/bot/ai) | path | level (info/warning/error) | text

//echo "Write ChatID:"
//read chat_id
//psql
//select distinct username from operations where chat_id = $chat_id;

type LevelOfLog int32

const (
	FatalError LevelOfLog = 0
	Error      LevelOfLog = 1
	Warning    LevelOfLog = 2
	Info       LevelOfLog = 4
)

var (
	Mutex sync.Mutex
)

type Log struct {
	Date   time.Time
	ChatID int64
	Author string
	Path   string
	Level  LevelOfLog
	Text   string
	//	IsError bool //del
}

func NewLog(u *UserInfo, name string, level LevelOfLog, text string) Log {

	var chatid int64
	var author, path string

	if u != nil {
		chatid = u.ChatID
		author = u.Username
		path = u.Path
	}

	if name != "" {
		author = name
	}

	return Log{
		Date:   time.Now().UTC().Add(3 * time.Hour),
		ChatID: chatid,
		Author: author,
		Path:   path,
		Level:  level,
		Text:   text,
	}
}

func SaveLogs() {

	for v := range Logs {
		if v.Text != "" {

			// время по Мск час. поясу
			timeNow := v.Date.Format(time.DateTime)

			// date - username - text
			fmt.Printf("(%s) [%s] %s\n", timeNow, v.Author, v.Text)

			// если это ошибка, то записываем в отдельный файл
			if v.Level == FatalError || v.Level == Error {
				WriteIntoFile(timeNow, v.Author, v.Text)
			}

			SQL_AddLog(v)
		}
	}
}

func WriteIntoFile(values ...any) {

	Mutex.Lock()
	defer Mutex.Unlock()

	file, _ := os.OpenFile("errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	file.WriteString(fmt.Sprintf("(%s) [%s] %s\n", values[0], values[1], values[2]))

}

func FinishGorutine(u *UserInfo, text string, main bool) {

	timeNow := time.Now().UTC().Add(3 * time.Hour).Format(time.DateTime)
	if r := recover(); r != nil {

		text := "text: " + text + "\n" + "Error: " + fmt.Sprint(r)

		fmt.Println(timeNow+" Panic in gorutine:", text)

		WriteIntoFile(timeNow, Ternary(main, "main", "gorutine"), text)

		SQL_AddLog(NewLog(u, "", FatalError, fmt.Sprint(r)))

	}

}
