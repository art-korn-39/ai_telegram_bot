package main

import (
	"os/exec"
	"strings"
)

func SendRequestToKandinsky(text string) (string, error) {

	res, err := exec.Command(`python`,
		`C:\DEV\GO\telegram_bot_1\scripts\generate_image.py`,
		//		filepath.Dir(os.Args[0])+`\scripts\generate_image.py`,
		//`scripts\generate_image.py`,
		text).
		Output()

	if err != nil {
		return "", err
	}

	pathToImage := strings.TrimSpace(string(res[:]))

	return pathToImage, nil

}
