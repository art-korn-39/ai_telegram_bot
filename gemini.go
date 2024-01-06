package main

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	//Gemini_APIKEY = "AIzaSyC0myz4bPIDyx6pPtW0PBZqmJW37A5VJ_k"
	URL = "https://generativelanguage.googleapis.com/v1beta3/models/text-bison-001:generateText"
)

var (
	ctx_Gemini    context.Context
	client_Gemini *genai.Client
	model_Gemini  *genai.GenerativeModel
)

// - FinishReasonSafety означает, что потенциальное содержимое было помечено по соображениям безопасности.
// - BlockReasonSafety означает, что промт был заблокирован по соображениям безопасности. Вы можете проверить
// `safety_ratings`, чтобы понять, какая категория безопасности заблокировала его.

func SendRequestToGemini(text string, user *UserInfo) string {

	<-delay_Gemini

	cs := model_Gemini.StartChat()
	cs.History = user.History_Gemini

	resp, err := cs.SendMessage(ctx_Gemini, genai.Text(text))
	if err != nil {
		errorString := err.Error()
		//		Logs <- Log{"Gemini", errorString, true}
		if errorString == "blocked: candidate: FinishReasonSafety" {
			NewConnectionGemini() // В случае данного вида ошибки - запускаем новый клиент соединения
			return "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса или очистить историю диалога командой /clearcontext."
		} else if errorString == "blocked: prompt: BlockReasonSafety" {
			return "Запрос был заблокирован по соображениям безопасности. Попробуйте изменить текст запроса."
		}
	}

	if resp.Candidates[0].Content == nil {
		//		Logs <- Log{"Gemini", "resp.Candidates[0].Content = nil", true}
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

func NewConnectionGemini() {
	ctx_Gemini = context.Background()
	client_Gemini, _ = genai.NewClient(ctx_Gemini, option.WithAPIKey(Cfg.GeminiKey))
	model_Gemini = client_Gemini.GenerativeModel("gemini-pro")
}
