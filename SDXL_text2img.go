package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const SDXL_URL_TEXT2IMG = "https://api.stability.ai/v1/generation/stable-diffusion-xl-1024-v1-0/text-to-image"

type TextToImage struct {
	Base64       string `json:"base64"`
	Seed         uint32 `json:"seed"`
	FinishReason string `json:"finishReason"`
}

type TextToImageResponse struct {
	Images []TextToImage `json:"artifacts"`
}

func sdxl_Text2image(user *UserInfo) (result string, err error) {

	inputText := user.Options["text"]
	style := user.Options["style"]

	type Text_prompt struct {
		Text   string  `json:"text"`
		Weight float64 `json:"weight"`
	}

	body_struct := struct {
		Text_prompts []Text_prompt `json:"text_prompts"`
		Cfg_scale    int           `json:"cfg_scale"`
		Height       int           `json:"height"`
		Width        int           `json:"width"`
		Samples      int           `json:"samples"`
		Steps        int           `json:"steps"`
		Style_preset string        `json:"style_preset,omitempty"`
	}{
		Text_prompts: []Text_prompt{
			{
				Text:   inputText,
				Weight: 0.5,
			},
		},
		Cfg_scale:    7,
		Height:       1024,
		Width:        1024,
		Samples:      1,
		Steps:        30,
		Style_preset: style,
	}

	body_json, _ := json.Marshal(body_struct)

	req, _ := http.NewRequest("POST", SDXL_URL_TEXT2IMG, bytes.NewBuffer(body_json))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+Cfg.Stability_Key)

	// Execute the request & read all the bytes of the body
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New("{http.DefaultClient.Do()} " + err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {

		msgText := GetText(MsgText_FailedGenerateImage2, user.Language)

		type res_body struct {
			Id      string `json:"id"`
			Name    string `json:"name"`
			Message string `json:"message"`
		}
		var body res_body

		if res.StatusCode/100 == 5 { //522, 520, 503
			err = errors.New("5XX status code =/")
			msgText = GetText(MsgText_APIdead, user.Language)
		} else if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
			err = errors.New("{Decode(&body) 1} " + err.Error())
		} else {
			err = fmt.Errorf("{%s} %s", body.Name, body.Message)
		}

		return msgText, err
	}

	// Decode the JSON body
	var body TextToImageResponse
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		err = errors.New("{Decode(&body) 2} " + err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}

	// Write the images to disk
	if len(body.Images) > 0 {
		image := body.Images[0]
		result, err = sdxl_CreateImageFileFromBase64(0, image, user)
		if err != nil {
			msgText := GetText(MsgText_UnexpectedError, user.Language)
			return msgText, err
		}
	} else {
		err = errors.New("len(body.Images) = 0")
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}

	return result, nil

}
