package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

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
		Date:   MskTimeNow(),
		ChatID: chatid,
		Author: author,
		Path:   path,
		Level:  level,
		Text:   text,
	}
}

func WriteIntoFile(values ...any) {

	Mutex.Lock()
	defer Mutex.Unlock()

	file, _ := os.OpenFile("errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	file.WriteString(fmt.Sprintf("(%s) %d [%s] %s\n", values[0], values[1], values[2], values[3]))

}

func FinishGorutine(u *UserInfo, text string, main bool) {

	timeNow := MskTimeNow().Format(time.DateTime)
	if r := recover(); r != nil {

		text := "text: " + text + "\n" + "Error: " + fmt.Sprint(r)

		fmt.Println(timeNow+" Panic in gorutine:", text)

		chatid := 0
		if u != nil {
			chatid = int(u.ChatID)
		}

		WriteIntoFile(timeNow, chatid, Ternary(main, "main", "gorutine"), text)

		SQL_AddLog(NewLog(u, "", FatalError, fmt.Sprint(r)))

	}

}
