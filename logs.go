package main

import (
	"fmt"
	"log"
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
			log.Printf("[%s] %s", v.UserName, v.Text)
			if v.IsError { // дополнительно: ошибку записываем в файл
				WriteIntoFile(v.UserName, v.Text)
			}
		}
	}
}

func LogPanic(inputtext string, main bool) {

	if r := recover(); r != nil {
		text := "Inputtext: " + inputtext + "\n" + "Error: " + fmt.Sprint(r)
		log.Println("Panic in gorutine.\n", text)
		WriteIntoFile(Ternary(main, "main", "gorutine"), text)
	}

}

func WriteIntoFile(values ...any) {

	Mutex.Lock()
	defer Mutex.Unlock()

	file, _ := os.OpenFile("Errors.txt", os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	file.WriteString(fmt.Sprintf("(%s) [%s] %s\n", time.Now().Format(time.DateTime), values[0], values[1]))

}
