package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

var (
	kandinskiy_API_Key    string
	kandinskiy_Secret_Key string
)

func init() {
	kandinskiy_API_Key = "1B189E2CFA69FFD2130FC56294B96DA9"
	kandinskiy_Secret_Key = "7C0B7C7FE4FFA4F6EE9DF0CFA257167C"
}

func kandinskiy_get_model() float64 {

	url := "https://api-key.fusionbrain.ai/key/api/v1/models"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("X-Key", "Key "+kandinskiy_API_Key)
	req.Header.Add("X-Secret", "Secret "+kandinskiy_Secret_Key)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	resBytes, _ := io.ReadAll(res.Body)
	var dat []map[string]any
	json.Unmarshal(resBytes, &dat)

	return dat[0]["id"].(float64)

}

func kandinskiy_generate_1(prompt string, model float64, style int) {

	url := "https://api-key.fusionbrain.ai/key/api/v1/text2image/run"

	params := map[string]interface{}{
		"type":      "GENERATE",
		"numImages": 1,
		"width":     1024,
		"height":    1024,
		"generateParams": map[string]string{
			"query": prompt,
		},
	}

	paramsJSON, _ := json.Marshal(params)

	body := &bytes.Buffer{}
	body.WriteString(fmt.Sprintf("model_id=%s&params=%s", model, paramsJSON))

	res, _ := http.Post(url, "application/x-www-form-urlencoded", body)
	defer res.Body.Close()

	resBytes, _ := io.ReadAll(res.Body)
	var dat any
	json.Unmarshal(resBytes, &dat)

}

func kandinskiy_generate_2(prompt string, model float64, style int) {

	url := "https://api-key.fusionbrain.ai/key/api/v1/text2image/run"

	params := map[string]interface{}{
		"type":      "GENERATE",
		"numImages": 1,
		"width":     1024,
		"height":    1024,
		"generateParams": map[string]string{
			"query": prompt,
		},
	}

	paramsJSON, _ := json.Marshal(params)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	modelField, _ := writer.CreateFormField("model_id")
	modelField.Write([]byte("4"))

	paramsField, _ := writer.CreateFormField("params")
	paramsField.Write(paramsJSON)

	writer.Close()

	req, _ := http.NewRequest("POST", url, body)

	// Устанавливаем заголовоки в теле запроса
	req.Header.Set("X-Key", "Key "+kandinskiy_API_Key)
	req.Header.Set("X-Secret", "Secret "+kandinskiy_Secret_Key)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполняем запрос
	res, _ := http.DefaultClient.Do(req)

	resBytes, _ := io.ReadAll(res.Body)
	var dat any
	json.Unmarshal(resBytes, &dat)

}

func kandinskiy_generate_3(prompt string, model float64, style int) {

	url := "https://api-key.fusionbrain.ai/key/api/v1/text2image/run"

	params := map[string]interface{}{
		"type":      "GENERATE",
		"numImages": 1,
		"width":     1024,
		"height":    1024,
		"generateParams": map[string]string{
			"query": prompt,
		},
	}

	paramsJSON, _ := json.Marshal(params)
	params_body := [3]any{nil, paramsJSON, "application/json"}

	body := &bytes.Buffer{}
	body.WriteString(fmt.Sprintf("model_id=%s&params=%s", model, params_body))
	//writer := multipart.NewWriter(body)

	json_new :=
		`{"type":      "GENERATE",
"numImages": 1,
"width":     1024,
"height":    1024,
"generateParams": {
"query": prompt,
},}`

	//var _nil *string
	params_new := myParameters{MyNil: nil, MyParams: json_new, MyType: "application/json"}
	model_new := myModel{nil, 4}
	body_new := myBody{Params: params_new, Model: model_new}

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.

	// Encode (send) the value.
	err := enc.Encode(body_new)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(network.Bytes()))

	// Устанавливаем заголовоки в теле запроса
	req.Header.Set("X-Key", "Key "+kandinskiy_API_Key)
	req.Header.Set("X-Secret", "Secret "+kandinskiy_Secret_Key)
	//req.Header.Set("Content-Type", "multipart/form-data")
	//req.Header.Set("Content-Type", writer.FormDataContentType())
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Выполняем запрос
	res, _ := http.DefaultClient.Do(req)

	resBytes, _ := io.ReadAll(res.Body)
	var dat any
	json.Unmarshal(resBytes, &dat)
}

type myBody struct {
	Model  myModel
	Params myParameters
}

type myModel struct {
	MyNil *string
	Id    int
}

type myParameters struct {
	MyNil    *string
	MyParams string
	MyType   string
}

func kandinskiy_generate_4(prompt string, model float64, style int) {

	url_str := "https://api-key.fusionbrain.ai/key/api/v1/text2image/run"

	// params := params{
	// 	Type:           "GENERATE",
	// 	NumImages:      1,
	// 	Width:          1024,
	// 	Height:         1024,
	// 	GenerateParams: map[string]string{"query": prompt},
	// }

	// JSON, _ := json.Marshal(params)

	// body := &bytes.Buffer{}
	// body.WriteString(fmt.Sprintf("model_id='%s';params='%s'", "4", JSON))

	json_string := "{\"type\":\"GENERATE\",\"numImages\":1,\"width\":1024,\"height\":1024,\"generateParams\":{\"query\":\"красный порше\"}}"
	body_2 := url.Values{
		"model_id": {"4"},
		"params":   {json_string, "application/json"},
	}

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	enc.Encode(body_2)

	// Request
	req, _ := http.NewRequest("POST", url_str, bytes.NewBuffer(network.Bytes()))

	// Устанавливаем заголовоки в теле запроса
	req.Header.Set("X-Key", "Key "+kandinskiy_API_Key)
	req.Header.Set("X-Secret", "Secret "+kandinskiy_Secret_Key)

	// Выполняем запрос
	res, _ := http.DefaultClient.Do(req)

	resBytes, _ := io.ReadAll(res.Body)
	var dat any
	json.Unmarshal(resBytes, &dat)
}

type params struct {
	Type           string            `json:"type"`
	NumImages      int               `json:"numImages"`
	Width          int               `json:"width"`
	Height         int               `json:"height"`
	GenerateParams map[string]string `json:"generateParams"`
}

func kandinskiy_generate_5(prompt string, model float64, style int) {

	url_str := "https://api-key.fusionbrain.ai/key/api/v1/text2image/run"

	body_str := `
	-F 'params="{
		\"type\":\"GENERATE\",
		\"generateParams\": {
			\"query\":\"красный порше\"
		}
	}";type=application/json' \
	--form 'model_id="1"'
	`

	body := []byte(body_str)

	// Request
	req, _ := http.NewRequest("POST", url_str, bytes.NewBuffer(body))

	// Устанавливаем заголовки в запрос
	req.Header.Set("X-Key", "Key "+kandinskiy_API_Key)
	req.Header.Set("X-Secret", "Secret "+kandinskiy_Secret_Key)
	req.Header.Set("Content-Type", "multipart/form-data")

	// Выполняем запрос
	res, _ := http.DefaultClient.Do(req)

	resBytes, _ := io.ReadAll(res.Body)
	var dat any
	json.Unmarshal(resBytes, &dat)
}

func kandinskiy_generate_6(prompt string, model float64, style int) {

	url_str := "https://api-key.fusionbrain.ai/key/api/v1/text2image/run"

	body_str := `{"model_id": [null, 4], "params": [null, "{"type": "GENERATE", "generateParams": {"query": "Sun in sky"}}", "application/json"]}`
	body_str = `{"model_id": [null, 4], "params": [null, "{\"type\": \"GENERATE\", \"generateParams\": {\"query\": \"Sun in sky\"}}", "application/json"]}`

	body := strings.NewReader(body_str)

	// Request
	req, _ := http.NewRequest("POST", url_str, body)

	// Устанавливаем заголовки в запрос
	req.Header.Set("X-Key", "Key "+kandinskiy_API_Key)
	req.Header.Set("X-Secret", "Secret "+kandinskiy_Secret_Key)
	req.Header.Set("Content-Type", "multipart/form-data")

	// Выполняем запрос
	res, _ := http.DefaultClient.Do(req)

	resBytes, _ := io.ReadAll(res.Body)
	var dat any
	json.Unmarshal(resBytes, &dat)
}

func kandinskiy_generate_7(prompt string, model float64, style int) {

	url_str := "https://api-key.fusionbrain.ai/key/api/v1/text2image/run"

	body_str := `{'params': '{"type": "GENERATE", "generateParams": {"query": "Sun in sky"}}', 'type': application/json}`

	body := bytes.NewBuffer([]byte(body_str))
	//body := strings.NewReader(body_str)
	//body := bytes.NewBufferString(body_str)

	// Request
	req, _ := http.NewRequest("POST", url_str, body)

	req.ParseForm()

	// Устанавливаем заголовки в запрос
	req.Header.Set("X-Key", "Key "+kandinskiy_API_Key)
	req.Header.Set("X-Secret", "Secret "+kandinskiy_Secret_Key)
	req.Form.Set("model_id", "4")

	//req.Header.Set("Content-Type", "multipart/form-data")

	// Выполняем запрос
	res, _ := http.DefaultClient.Do(req)

	resBytes, _ := io.ReadAll(res.Body)
	var dat any
	json.Unmarshal(resBytes, &dat)
}
