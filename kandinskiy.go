package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func SendRequestToKandinsky(text string) (string, error) {

	_, callerFile, _, _ := runtime.Caller(0)
	dir := strings.ReplaceAll(filepath.Dir(callerFile), "\\", "/")
	scriptPath := dir + "/scripts/generate_image.py"
	dataFolder := dir + "/data"

	cmd := exec.Command(`python`,
		scriptPath,
		dataFolder,
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
