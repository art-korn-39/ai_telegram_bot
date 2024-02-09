package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	SDXL_Styles = map[string]string{
		"No style":    "",
		"3D model":    "3d-model",
		"Photo":       "photographic",
		"Neon punk":   "neon-punk",
		"Cinematic":   "cinematic",
		"Analog film": "analog-film",
		"Fantasy art": "fantasy-art",
		"Anime":       "anime",
	}

	SDXL_URL_TEXT2IMG = "https://api.stability.ai/v1/generation/stable-diffusion-xl-1024-v1-0/text-to-image"
	SDXL_URL_UPSCALE  = "https://api.stability.ai/v1/generation/esrgan-v1-x2plus/image-to-image/upscale"
)

func sdxl_WelcomeTextMessage(u *UserInfo) string {

	duration := GetDurationToNextDay()
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - hours*60

	return fmt.Sprintf(GetText(MsgText_SDXLinfo, u.Language),
		max(Get_RPD_sdxl(u)-u.Requests_today_sdxl, 0),
		hours,
		minutes)
}

func sdxl_DailyLimitOfRequestsIsOver(u *UserInfo, btn Button) bool {

	if u.Requests_today_sdxl >= Get_RPD_sdxl(u) {
		duration := GetDurationToNextDay()
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) - hours*60
		msgText := fmt.Sprintf(GetText(MsgText_DailyRequestLimitExceeded, u.Language), hours, minutes)
		if btn == 0 {
			SendMessage(u, msgText, nil, "")
		} else {
			SendMessage(u, msgText, GetButton(btn_Models, u.Language), "")
			u.Path = "start"
		}
		return true
	}

	return false

}

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
		var body map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
			err = errors.New("{Decode(&body)} " + err.Error())
		} else {
			err = errors.New("{StatusCode != 200} " + fmt.Sprintf("%s", body))
		}
		msgText := GetText(MsgText_FailedGenerateImage2, user.Language)
		return msgText, err
	}

	// Decode the JSON body
	var body TextToImageResponse
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		err = errors.New("{Decode(&body)} " + err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}

	// Write the images to disk
	if len(body.Images) > 0 {
		image := body.Images[0]
		result, err = sdxl_CreateImageFile(0, image, user)
		if err != nil {
			msgText := GetText(MsgText_UnexpectedError, user.Language)
			return msgText, err
		}
	}

	return result, nil

}

func sdxl_Image2ImageUpscale(user *UserInfo, filepath string) (result string, err error) {

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	// Write the init image to the request
	imageWriter, _ := writer.CreateFormField("image")
	imageFile, err := os.Open(filepath)
	if err != nil {
		err = errors.New("{os.Open()} " + err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}
	io.Copy(imageWriter, imageFile)

	// Write the options to the request
	writer.WriteField("width", "2048")

	// +++
	//writer.WriteField("text_prompts[0][text]", "highly detailed")
	// ---

	writer.Close()

	// Execute the request
	payload := bytes.NewReader(data.Bytes())
	req, _ := http.NewRequest("POST", SDXL_URL_UPSCALE, payload)

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "image/png")
	req.Header.Add("Authorization", "Bearer "+Cfg.Stability_Key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New("{http.DefaultClient.Do()} " + err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var body map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
			err = errors.New("{Decode(&body)} " + err.Error())
		} else {
			err = errors.New("{StatusCode != 200} " + fmt.Sprintf("%s", body))
		}
		msgText := GetText(MsgText_FailedGenerateImage2, user.Language)
		return msgText, err
	}

	if res.Header["Finish-Reason"][0] == "CONTENT_FILTERED" {
		err = errors.New("CONTENT_FILTERED")
		msgText := GetText(MsgText_FailedImageUpscale, user.Language)
		return msgText, err
	}

	// Write the response to a file
	filepath = fmt.Sprintf(WorkDir+"/data/img_%d_0.png", user.ChatID)
	file, err := os.Create(filepath)
	if err != nil {
		err = errors.New("{os.Create()} " + err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		err = errors.New("{io.Copy} " + err.Error())
		msgText := GetText(MsgText_UnexpectedError, user.Language)
		return msgText, err
	}

	return filepath, nil

}

func sdxl_CreateImageFile(i int, image TextToImage, user *UserInfo) (string, error) {

	outFile := fmt.Sprintf(WorkDir+"/data/img_%d_0.png", user.ChatID)
	file, err := os.Create(outFile)
	if err != nil {
		err = errors.New("{os.Create()} " + err.Error())
		return "", err
	}

	imageBytes, err := base64.StdEncoding.DecodeString(image.Base64)
	if err != nil {
		err = errors.New("{DecodeString()} " + err.Error())
		return "", err
	}

	if _, err := file.Write(imageBytes); err != nil {
		err = errors.New("{file.Write()} " + err.Error())
		return "", err
	}

	if err := file.Close(); err != nil {
		err = errors.New("{file.Close()} " + err.Error())
		return "", err
	}

	return outFile, nil

}
