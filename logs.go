package main

import (
	"fmt"
	"log"
	"os"
)

type Log struct {
	UserName string
	Text     string
	IsError  bool
}

func saveLogs() {

	for v := range Logs {
		log.Printf("[%s] %s", v.UserName, v.Text)
		if v.IsError { // дополнительно: ошибку записываем в файл
			err_file, _ := os.OpenFile("Errors.txt", os.O_APPEND|os.O_WRONLY, 0600)
			err_file.WriteString(fmt.Sprintf("[%s] %s\n", v.UserName, v.Text))
			err_file.Close()
		}
	}
}

func logPanic() {
	if r := recover(); r != nil {
		log.Println("Panic in gorutine. Error:\n", r)
	}
}
