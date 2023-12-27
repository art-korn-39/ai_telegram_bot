package main

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	Gemini_APIKEY = "AIzaSyC0myz4bPIDyx6pPtW0PBZqmJW37A5VJ_k"
	URL           = "https://generativelanguage.googleapis.com/v1beta3/models/text-bison-001:generateText"
)

var (
	ctx_Gemini    context.Context
	client_Gemini *genai.Client
	model_Gemini  *genai.GenerativeModel
)

func init() {
	ctx_Gemini = context.Background()
	client_Gemini, _ = genai.NewClient(ctx_Gemini, option.WithAPIKey(Gemini_APIKEY))
	model_Gemini = client_Gemini.GenerativeModel("gemini-pro")
}

func SendRequestToGemini(text string, user *UserInfo) string {

	<-delay_Gemini

	cs := model_Gemini.StartChat()
	cs.History = user.History_Gemini

	resp, err := cs.SendMessage(ctx_Gemini, genai.Text(text))
	if err != nil {
		errorString := err.Error()
		if errorString == "blocked: candidate: FinishReasonSafety" {

			// В случае данного вида ошибка - запускаем новый клиент соединения
			ctx_Gemini = context.Background()
			client_Gemini, _ = genai.NewClient(ctx_Gemini, option.WithAPIKey(Gemini_APIKEY))
			model_Gemini = client_Gemini.GenerativeModel("gemini-pro")

		}
		Logs <- Log{"Gemini", errorString, true}
		return "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса."
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)

	history := append(user.History_Gemini,
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(text),
			},
			Role: "user",
		},
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(result),
			},
			Role: "model",
		},
	)

	user.History_Gemini = history

	return string(result)

}
