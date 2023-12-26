package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func SendRequestToKandinsky(text string, userid int64) (result string, isError bool) {

	<-delay_Kandinsky

	_, callerFile, _, _ := runtime.Caller(0)
	dir := strings.ReplaceAll(filepath.Dir(callerFile), "\\", "/")
	scriptPath := dir + "/scripts/generate_image.py"
	dataFolder := dir + "/data"

	cmd := exec.Command(`python`,
		scriptPath,
		dataFolder,
		text,
		strconv.Itoa(int(userid)))

	if cmd.Err != nil {
		Logs <- Log{"Kandinsky", "request: " + text + "\nerror: " + cmd.Err.Error(), true}
		return "Не удалось сгенерировать изображение. Попробуйте позже.", true
	}

	// Получение результата команды
	res, err := cmd.Output()

	if err != nil {
		Logs <- Log{"Kandinsky", "request: " + text + "\nerror: " + err.Error(), true}
		return "Не удалось сгенерировать изображение. Попробуйте позже.", true
	}

	pathToImage := strings.TrimSpace(string(res[:]))

	return pathToImage, false

}
