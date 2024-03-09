package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// starting
// processing
// succeed
// succeeded
// failed
// canceled

type fs_body_request struct {
	Id     string   `json:"id"`
	Logs   string   `json:"logs"`
	Error  error    `json:"error"`
	Status string   `json:"status"`
	URLS   []string `json:"urls"`
}

func fs_SendRequest(user *UserInfo) (result string, msg Text, err error) {

	body_struct := struct {
		Version string         `json:"version"`
		Input   map[string]any `json:"input"`
	}{
		Version: Cfg.Faceswap_Version,
		Input: map[string]any{
			"request_id":   Cfg.Faceswap_id,
			"local_source": user.Options["image1"], // лицо
			"local_target": user.Options["image2"], // фон
			"weight":       0.5,
			"cache_days":   1,
			"det_thresh":   0.1,
		},
	}

	body_json, _ := json.Marshal(body_struct)

	req, _ := http.NewRequest("POST", "https://api.replicate.com/v1/predictions", bytes.NewBuffer(body_json))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+Cfg.Replicate_Key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	resBytes, _ := io.ReadAll(res.Body)

	var dat fs_body_request
	json.Unmarshal(resBytes, &dat)

	if dat.Status == "starting" {
		url, msg, err := fs_GetURLofResult(dat.Id)
		if err != nil {
			return "", msg, err
		} else {
			outFile := getFilepathForImage(user.ChatID, "jpg")
			err = DownloadFileByURL(outFile, url)
			if err != nil {
				return "", MsgText_UnexpectedError, err
			} else {
				return outFile, MsgText_nil, nil
			}
		}
	}

	return

}

type fs_body_response struct {
	Id   string `json:"id"`
	Logs string `json:"logs"`

	Output struct {
		Code   float64 `json:"code"`
		Image  string  `json:"image"`
		Msg    string  `json:"msg"`
		Status string  `json:"status"`
	} `json:"output"`

	Created_at   time.Time `json:"created_at"`
	Started_at   time.Time `json:"started_at"`
	Completed_at time.Time `json:"completed_at"`
	Status       string    `json:"status"`
	Error        error     `json:"error"`

	Metrics struct {
		Predict_time float64 `json:"predict_time"`
	} `json:"metrics"`
}

func fs_GetURLofResult(id string) (image string, msg Text, err error) {

	delay := time.Tick(1500 * time.Millisecond)

	var i int
	for {

		i++
		<-delay

		req, _ := http.NewRequest("GET", "https://api.replicate.com/v1/predictions/"+id, nil)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Token "+Cfg.Replicate_Key)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()

		resBytes, _ := io.ReadAll(res.Body)

		var dat fs_body_response
		json.Unmarshal(resBytes, &dat)

		switch dat.Status {
		case "succeeded":

			switch dat.Output.Status {
			case "succeed":
				return dat.Output.Image, MsgText_nil, nil
			case "failed":
				err := fmt.Errorf("code:%f, msg:%s", dat.Output.Code, dat.Output.Msg)
				if dat.Output.Msg == "no face" {
					return "", MsgText_NoFaceFound, err
				} else {
					return "", MsgText_UnexpectedError, err
				}
			default:
				err := fmt.Errorf("code:%f, msg:%s, status:%s", dat.Output.Code, dat.Output.Msg, dat.Output.Status)
				return "", MsgText_UnexpectedError, err
			}
		case "failed":
			return "", MsgText_UnexpectedError, dat.Error
		default:
			// Для подстраховки, чтобы вечный цикл не получился, удалить потом
			if i > 20 {
				return "", MsgText_UnexpectedError, errors.New("Прервано. Ожидание больше 30 секунд, статус:" + dat.Status)
			}
			continue
		}
	}
}
