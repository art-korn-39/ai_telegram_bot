package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func SendRequestToKandinsky(text string) (string, error) {

	dir := filepath.Dir(os.Args[0])

	log.Println("путь к баз. каталогу", dir)

	cmd := exec.Command(`python3`,
		dir+`\scripts\generate_image.py`,
		strings.ReplaceAll(dir+`\data`, "\\", "/"),
		text)

	if cmd.Err != nil {
		return "", cmd.Err
	}

	res, err := cmd.Output()

	if err != nil {
		return "", err
	}

	pathToImage := strings.TrimSpace(string(res[:]))

	return pathToImage, nil

}
