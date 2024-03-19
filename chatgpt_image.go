package main

import (
	"context"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
)

// После отправки картинок пользователем
func gpt_image(user *UserInfo, message *tgbotapi.Message) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	// Проверяем наличие картинок в сообщении
	if message.Photo == nil {
		msgText := GetText(MsgText_UploadImage, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	SendMessage(user, GetText(MsgText_LoadingImage, user.Language), nil, "")

	photos := *message.Photo

	// Получаем URL картинки
	fileURL, err := Bot.GetFileDirectURL(photos[len(photos)-1].FileID)
	if err != nil {
		Logs <- NewLog(user, "chatgpt{img1}", Error, err.Error())
		msgText := GetText(MsgText_FailedLoadImages, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	user.Options["fileURL"] = fileURL

	// Просим написать запрос к картинке
	msgText := GetText(MsgText_PhotoUploadedWriteQuestion, user.Language)
	SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")

	user.Path = "chatgpt/type/image/text"

}

// После ввода вопроса пользователем
func gpt_imgtext(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	// Проверяем наличие текста в сообщении
	if text == "" {
		msgText := GetText(MsgText_WriteQuestionToImage, user.Language)
		SendMessage(user, msgText, GetButton(btn_RemoveKeyboard, ""), "")
		return
	}

	<-delay_ChatGPT

	Operation := SQL_NewOperation(user, "chatgpt", "img", text)
	SQL_AddOperation(Operation)

	resp, err := clientOpenAI.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4VisionPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{Type: "text",
							Text: text,
						},
						{Type: "image_url",
							ImageURL: &openai.ChatMessageImageURL{
								URL:    user.Options["fileURL"],
								Detail: openai.ImageURLDetailLow,
							},
						},
					},
				},
			},
			MaxTokens: 200,
		},
	)

	if err != nil {
		Logs <- NewLog(user, "chatgpt{img2}", Error, err.Error())
		msgText := GetText(MsgText_BadRequest1, user.Language)
		SendMessage(user, msgText, GetButton(btn_GptImageNewgen, user.Language), "")
		user.Path = "chatgpt/type/image/text/newgen"
		return
	}

	// умножаем на 20, хотя в реале соотношение цены токенов gpt3.5 к gpt4 ~ 26-27
	user.Tokens_used_gpt = user.Tokens_used_gpt + resp.Usage.TotalTokens*20
	user.Usage.GPT = user.Usage.GPT + resp.Usage.TotalTokens*20

	content := resp.Choices[0].Message.Content

	SendMessage(user, content, GetButton(btn_GptImageNewgen, user.Language), "")

	user.Path = "chatgpt/type/image/text/newgen"

}

// После ответа пользователя на результат по вопросу и картинкам
func gpt_imgtext_newgen(user *UserInfo, text string) {

	if gpt_DailyLimitOfTokensIsOver(user) {
		return
	}

	switch text {
	// ИЗМЕНИТЬ ЗАПРОС
	case GetText(BtnText_ChangeQuerryText, user.Language):
		SendMessage(user, GetText(MsgText_WriteQuestionToImage, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "chatgpt/type/image/text"

	// ЗАГРУЗИТЬ НОВЫЕ ФОТО
	case GetText(BtnText_UploadNewImage, user.Language):
		user.DeleteImages() // на всякий почистим, если что-то осталось
		SendMessage(user, GetText(MsgText_UploadImage, user.Language), GetButton(btn_RemoveKeyboard, ""), "")
		user.Path = "chatgpt/type/image"

	// НАЧАТЬ ДИАЛОГ
	case GetText(BtnText_StartDialog, user.Language):
		user.DeleteImages() // на всякий почистим, если что-то осталось
		SendMessage(user, GetText(MsgText_HelloCanIHelpYou, user.Language), GetButton(btn_GptClearContext, user.Language), "")
		user.Path = "chatgpt/type/dialog"

	default:
		// Предполагаем, что там новый вопрос к загруженным картинкам
		gpt_imgtext(user, text)
	}

}
