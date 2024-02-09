package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type bl_response struct {
	Err    string   `json:"error"`
	Status string   `json:"status"`
	Files  []string `json:"output,omitempty"`
}

func bl_createImage() {

	api := "sk-2810764ef697391d86f15f0f9751c9acd56b791735345e4b9554c6f8ba1107ce"
	url := "https://www.basedlabs.ai/api/v1/create/image"

	body_struct := struct {
		ModelId string `json:"modelId"`
		Prompt  string `json:"prompt"`
	}{
		ModelId: "von-speed", //"bluematter/basedlabs-comfyui-v1", //von-speed
		Prompt:  "A scenic mountain landscape",
	}

	body_json, _ := json.Marshal(body_struct)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body_json))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api)

	// Execute the request & read all the bytes of the body
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()

	var dat map[string]string
	if err = json.NewDecoder(res.Body).Decode(&dat); err != nil {
		fmt.Println(err.Error())
		return
	}

}

func bl_getImage() {

	prompt_id := "gqutb23biia4glhdutlvu6i6ou"

	api := "sk-2810764ef697391d86f15f0f9751c9acd56b791735345e4b9554c6f8ba1107ce"
	url := "https://www.basedlabs.ai/api/v1/read/image/"

	body_struct := struct {
		Id string `json:"id"`
	}{
		Id: prompt_id,
	}

	body_json, _ := json.Marshal(body_struct)

	req, _ := http.NewRequest("POST", url+prompt_id, bytes.NewBuffer(body_json))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()

	var dat any
	if err = json.NewDecoder(res.Body).Decode(&dat); err != nil {
		fmt.Println(err.Error())
		return
	}

	var dat2 bl_response
	if err = json.NewDecoder(res.Body).Decode(&dat2); err != nil {
		fmt.Println(err.Error())
		return
	}

	// if dat.Status == "succeeded" {

	// }

	//"bluematter/basedlabs-comfyui-v1"

}
