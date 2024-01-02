package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var Mutex sync.Mutex

type Log struct {
	UserName string
	Text     string
	IsError  bool
}

func SaveLogs() {

	for v := range Logs {
		if v.Text != "" {
			timeNow := time.Now().UTC().Add(3 * time.Hour).Format(time.DateTime)
			fmt.Printf("%s [%s] %s\n", timeNow, v.UserName, v.Text)
			if v.IsError { // дополнительно: ошибку записываем в файл
				WriteIntoFile(timeNow, v.UserName, v.Text)
			}
		}
	}
}

func WriteIntoFile(values ...any) {

	Mutex.Lock()
	defer Mutex.Unlock()

	file, _ := os.OpenFile("errors.txt", os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	file.WriteString(fmt.Sprintf("(%s) [%s] %s\n", values[0], values[1], values[2]))

}
