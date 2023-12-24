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
)

func init() {

	ctx_Gemini = context.Background()
	client_Gemini, _ = genai.NewClient(ctx_Gemini, option.WithAPIKey(Gemini_APIKEY))
}

func SendRequestToGemini(text string) string {

	// For text-only input, use the gemini-pro model
	model := client_Gemini.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx_Gemini, genai.Text(text))
	if err != nil {
		errorString := err.Error()
		if errorString == "blocked: candidate: FinishReasonSafety" {
			client_Gemini, _ = genai.NewClient(ctx_Gemini, option.WithAPIKey(Gemini_APIKEY))
		}
		Logs <- Log{"Gemini", err.Error(), true}
		return "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса."
	}

	result := resp.Candidates[0].Content.Parts[0].(genai.Text)
	return string(result)

}
