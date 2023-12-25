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

func saveLogs() {

	for v := range Logs {
		log.Printf("[%s] %s", v.UserName, v.Text)
		if v.IsError { // дополнительно: ошибку записываем в файл
			writeIntoFile(v.UserName, v.Text)
		}
	}
}

func logPanic(main bool) {
	if r := recover(); r != nil {
		log.Println("Panic in gorutine. Error:\n", r)
		writeIntoFile(Ternary(main, "main", "gorutine"), r)
	}
}

func writeIntoFile(values ...any) {

	Mutex.Lock()
	defer Mutex.Unlock()

	file, _ := os.OpenFile("Errors.txt", os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	file.WriteString(fmt.Sprintf("(%s) [%s] %s\n", time.Now().Format(time.DateTime), values[0], values[1]))

}
