package main

import "strconv"

type Text int

type MultiText struct {
	ru, en string
}

var dictionary = map[Text]MultiText{}

const (

	// MESSAGES

	// COMMON

	MsgText_Start                   Text = iota
	MsgText_LastOperationInProgress      // Последняя операция ещё выполняется, дождитесь её завершения перед отправкой новых запросов.
	MsgText_SubscribeForUsing            // Для продолжения использования бота необходимо подписаться на канал👇
	MsgText_UnexpectedError              // Произошла непредвиденная ошибка, попробуйте позже.
	MsgText_AiNotSelected                // Не выбрана нейросеть для обработки запроса.
	MsgText_AfterRecoveryProd            // Функциональность бота восстановлена. Приносим извинения за неудобства.
	MsgText_AfterRecoveryDebug           // Этот бот предназначен для тестирования и отладки, полностью рабочий и бесплатный находится здесь: @AI_free_chat_bot
	MsgText_HelloCanIHelpYou             // Привет! Чем могу помочь?
	MsgText_SelectOption                 // Выберите один из предложенных вариантов:
	MsgText_UnknownCommand               // Неизвестная команда
	MsgText_EndDialog                    // Завершить диалог
	MsgText_LanguageChanged              // Язык успешено изменён!

	// GEMINI

	MsgText_GeminiHello                 // Вас приветствует Gemini Pro от компании Google 🚀
	MsgText_DailyRequestLimitExceeded   // Превышен дневной лимит запросов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.
	MsgText_WriteQuestionToImages       // Напишите свой вопрос к загруженным изображениям
	MsgText_UploadImages                // Загрузите одну или несколько картинок
	MsgText_PhotosUploadedWriteQuestion // Загружено фотографий: %d\nНапишите свой вопрос.\nНапример:\n\"Кто на фотографии?\"\n\"Чем отличаются эти картинки?\"
	MsgText_LoadingImages               // Выполняется загрузка изображений...
	MsgText_FailedLoadImages            // Не удалось загрузить изображение, попробуйте ещё раз.

	// CHATGPT

	MsgText_ChatGPTHello                          // Вас приветствует ChatGPT 3.5 Turbo 🤖\n\nТекущий остаток токенов: <b>%d</b> <i>(обновится через: %d ч. %d мин.)</i>
	MsgText_LimitOf4097TokensReached              // Достигнут лимит в 4097 токенов, контекст диалога очищен.
	MsgText_SelectVoice                           // Выберите голос для озвучивания текста:
	MsgText_EnterTextForAudio                     // Введите текст для аудио:
	MsgText_WriteTextForVoicing                   // Напишите текст для озвучивания:
	MsgText_ErrorSendingAudioFile                 // При отправке аудиофайла возникла ошибка, попробуйте ещё раз позже.
	MsgText_ResultAudioGeneration                 // Результат генерации по тексту "%s", голос: "%s"
	MsgText_AudioFileCreationStarted              // Запущено создание аудиофайла...
	MsgText_VoiceExamples                         // Примеры звучания голосов👇
	MsgText_SelectVoiceFromOptions                // Выберите голос из предложенных вариантов.
	MsgText_NotEnoughTokensWriteShorterTextLength // Недостаточно токенов, укажите текст меньшей длины.
	MsgText_ChatGPTDialogStarted                  // Запущен диалог с СhatGPT, чтобы очистить контекст от предыдущих сообщений - нажмите кнопку "Очистить контекст". Это позволяет сократить расход токенов.
	MsgText_DialogContextCleared                  // Контекст диалога очищен
	MsgText_GenerateAudioFromText                 // Сгенерировать аудио из текста
	MsgText_DailyTokenLimitExceeded               // Превышен дневной лимит токенов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.
	MsgText_ErrorWhileProcessingRequest           // Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже.
	MsgText_WriteQuestionToImage                  // Напишите свой вопрос к загруженному изображению.
	MsgText_UploadImage                           // Загрузите картинку
	MsgText_PhotoUploadedWriteQuestion            // Напишите свой запрос.\nНапример:\n\"Реши тест на картинке\"\n\"Как называется это блюдо?\"
	MsgText_LoadingImage                          // Выполняется загрузка изображения...

	// KANDINSKY

	MsgText_EnterYourRequest                // Введите свой запрос:
	MsgText_DescriptionTextNotExceed900Char // Текст описания картинки не должен превышать 900 символов.
	MsgText_SelectStyleForImage             // Выберите стиль, в котором генерировать изображение.
	MsgText_SelectStyleFromOptions          // Выберите стиль из предложенных вариантов.
	MsgText_ImageGenerationStarted          // Запущена генерация картинки, среднее время выполнения 30-40 секунд.
	MsgText_ResultImageGeneration           // Результат генерации по запросу "%s", стиль: "%s"
	MsgText_ErrorWhileSendingPicture        // При отправке картинки возникла ошибка, попробуйте ещё раз позже.
	MsgText_FailedGenerateImage1            // Не удалось сгенерировать изображение, попробуйте позже.
	MsgText_FailedGenerateImage2            // Не удалось сгенерировать изображение. Попробуйте изменить текст описания картинки.

	// BAD REQUEST

	MsgText_BadRequest1 // Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса или использовать другие изображения.
	MsgText_BadRequest2 // Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса.
	MsgText_BadRequest3 // Не удалось получить ответ от сервиса. Попробуйте изменить текст вопроса или начать новый диалог.
	MsgText_BadRequest4 // Запрос был заблокирован по соображениям безопасности. Попробуйте изменить текст запроса.

	// BUTTONS

	BtnText_Subscribe             // ✅Подписаться
	BtnText_SendPictureWithText   // 🖼 AI Vision
	BtnText_ChangeQuestionText    // Изменить вопрос
	BtnText_ChooseAnotherVoice    // Изменить голос
	BtnText_ChangeQuerryText      // Изменить запрос
	BtnText_ChooseAnotherStyle    // Изменить стиль
	BtnText_ChangeText            // Изменить текст
	BtnText_UploadNewImages       // Загрузить новые фото
	BtnText_UploadNewImage        // Загрузить новое фото
	BtnText_EndDialog             // Завершить диалог
	BtnText_StartDialog           // Начать диалог
	BtnText_GenerateAudioFromText // Озвучить текст
	BtnText_ClearContext          // Очистить контекст
)

func init() {

	// common
	dictionary[MsgText_Start] = textForStarting()

	dictionary[MsgText_ChatGPTHello] = MultiText{
		ru: "Вас приветствует ChatGPT 3.5 Turbo 🤖\n\nТекущий остаток токенов: <b>%d</b> <i>(обновится через: %d ч. %d мин.)</i>",
		en: "Welcome to ChatGPT 3.5 Turbo 🤖\n\nCurrent balance of tokens: <b>%d</b> <i>(updated in: %d hours %d min.)</i>"}
	dictionary[MsgText_GeminiHello] = MultiText{
		ru: "Вас приветствует Gemini Pro от компании Google 🚀",
		en: "Welcome to Gemini Pro from Google 🚀"}
	dictionary[MsgText_EnterYourRequest] = MultiText{
		ru: "Введите свой запрос:",
		en: "Enter your request:"}
	dictionary[MsgText_EnterTextForAudio] = MultiText{
		ru: "Введите текст для аудио:",
		en: "Enter text for audio:"}
	dictionary[MsgText_ErrorWhileProcessingRequest] = MultiText{
		ru: "Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже.",
		en: "An error occurred while processing the request. Please try again later."}
	dictionary[MsgText_SelectVoice] = MultiText{
		ru: "Выберите голос для озвучивания текста:",
		en: "Select a voice to read the text:"}
	dictionary[MsgText_SelectVoiceFromOptions] = MultiText{
		ru: "Выберите голос из предложенных вариантов.",
		en: "Select a voice from the options provided."}
	dictionary[MsgText_SelectOption] = MultiText{
		ru: "Выберите один из предложенных вариантов:",
		en: "Select one of the following options:"}
	dictionary[MsgText_SelectStyleForImage] = MultiText{
		ru: "Выберите стиль, в котором генерировать изображение.",
		en: "Select the style in which to generate the image."}
	dictionary[MsgText_SelectStyleFromOptions] = MultiText{
		ru: "Выберите стиль из предложенных вариантов.",
		en: "Select a style from the options provided."}
	dictionary[MsgText_LoadingImages] = MultiText{
		ru: "Выполняется загрузка изображений...",
		en: "Loading images..."}
	dictionary[MsgText_LoadingImage] = MultiText{
		ru: "Выполняется загрузка изображения...",
		en: "Loading image..."}
	dictionary[MsgText_SubscribeForUsing] = MultiText{
		ru: "Для продолжения использования бота необходимо подписаться на канал👇",
		en: "To continue using the bot you must subscribe to the channel👇"}
	dictionary[MsgText_LimitOf4097TokensReached] = MultiText{
		ru: "Достигнут лимит в 4097 токенов, контекст диалога очищен.",
		en: "The limit of 4097 tokens has been reached, the dialog context has been cleared."}
	dictionary[MsgText_EndDialog] = MultiText{
		ru: "Завершить диалог",
		en: "End dialog"}
	dictionary[MsgText_PhotosUploadedWriteQuestion] = MultiText{
		ru: "Загружено фотографий: %d\nНапишите свой вопрос.\nНапример:\n\"Кто на фотографии?\"\n\"Чем отличаются эти картинки?\"",
		en: "Photos uploaded: %d\nWrite your question.\nFor example:\n\"Who is in the photo?\"\n\"What is the difference between these pictures?\""}
	dictionary[MsgText_PhotoUploadedWriteQuestion] = MultiText{
		ru: "Напишите свой запрос.\nНапример:\n\"Реши тест на картинке\"\n\"Как называется это блюдо?\"",
		en: "Write your request.\nFor example:\n\"Solve the test in the picture\"\n\"What is the name of this dish?\""}
	dictionary[MsgText_UploadImages] = MultiText{
		ru: "Загрузите одну или несколько картинок",
		en: "Upload one or more images"}
	dictionary[MsgText_UploadImage] = MultiText{
		ru: "Загрузите одну картинку",
		en: "Upload image"}
	dictionary[MsgText_BadRequest4] = MultiText{
		ru: "Запрос был заблокирован по соображениям безопасности. Попробуйте изменить текст запроса.",
		en: "The request was blocked for security reasons. Try changing the request text."}
	dictionary[MsgText_ChatGPTDialogStarted] = MultiText{
		ru: `Запущен диалог с СhatGPT, чтобы очистить контекст от предыдущих сообщений - нажмите кнопку "Очистить контекст". Это позволяет сократить расход токенов.`,
		en: `A dialog has started with ChatGPT, to clear the context from previous messages - click the "Clear context" button. This allows you to reduce the consumption of tokens.`}
	dictionary[MsgText_ImageGenerationStarted] = MultiText{
		ru: "Запущена генерация картинки, среднее время выполнения 30-40 секунд.",
		en: "The generation of the image has started, the average execution time is 30-40 seconds."}
	dictionary[MsgText_AudioFileCreationStarted] = MultiText{
		ru: "Запущено создание аудиофайла...",
		en: "Audio file creation started..."}
	dictionary[MsgText_DialogContextCleared] = MultiText{
		ru: "Контекст диалога очищен",
		en: "The dialog context has been cleared"}
	dictionary[MsgText_WriteQuestionToImages] = MultiText{
		ru: "Напишите свой вопрос к загруженным изображениям",
		en: "Write your question to the uploaded images"}
	dictionary[MsgText_WriteQuestionToImage] = MultiText{
		ru: "Напишите свой вопрос к загруженному изображению.",
		en: "Write your question to the uploaded image."}
	dictionary[MsgText_WriteTextForVoicing] = MultiText{
		ru: "Напишите текст для озвучивания:",
		en: "Write the text for voicing:"}
	dictionary[MsgText_AiNotSelected] = MultiText{
		ru: "Не выбрана нейросеть для обработки запроса.",
		en: "The neural network for processing requests has not been selected."}
	dictionary[MsgText_FailedLoadImages] = MultiText{
		ru: "Не удалось загрузить изображение, попробуйте ещё раз.",
		en: "Failed to load image, try again. Failed to load image, try again."}
	dictionary[MsgText_BadRequest1] = MultiText{
		ru: "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса или использовать другие изображения.",
		en: "Failed to receive a response from the service. Try changing your request text or using different images."}
	dictionary[MsgText_BadRequest2] = MultiText{
		ru: "Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса.",
		en: "Failed to receive a response from the service. Try changing the request text."}
	dictionary[MsgText_BadRequest3] = MultiText{
		ru: "Не удалось получить ответ от сервиса. Попробуйте изменить текст вопроса или начать новый диалог.",
		en: "Failed to receive a response from the service. Try changing the question text or starting a new dialogue."}
	dictionary[MsgText_FailedGenerateImage1] = MultiText{
		ru: "Не удалось сгенерировать изображение, попробуйте позже.",
		en: "Failed to generate image, please try again later."}
	dictionary[MsgText_FailedGenerateImage2] = MultiText{
		ru: "Не удалось сгенерировать изображение. Попробуйте изменить текст описания картинки.",
		en: "Failed to generate image. Try changing the text of the picture description."}
	dictionary[MsgText_NotEnoughTokensWriteShorterTextLength] = MultiText{
		ru: "Недостаточно токенов, укажите текст меньшей длины.",
		en: "There are not enough tokens, please specify a shorter text length."}
	dictionary[MsgText_UnknownCommand] = MultiText{
		ru: "Неизвестная команда",
		en: "Unknown command"}
	dictionary[MsgText_LastOperationInProgress] = MultiText{
		ru: "Последняя операция ещё выполняется, дождитесь её завершения перед отправкой новых запросов.",
		en: "The last operation is still in progress, please wait until it completes before sending new requests."}
	dictionary[MsgText_DailyRequestLimitExceeded] = MultiText{
		ru: "Превышен дневной лимит запросов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.",
		en: "The daily request limit has been exceeded, wait until the limit is updated (%d hours %d min.) or use another neural network."}
	dictionary[MsgText_DailyTokenLimitExceeded] = MultiText{
		ru: "Превышен дневной лимит токенов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.",
		en: "The daily token limit has been exceeded, wait until the limit is updated (%d hours %d min.) or use another neural network."}
	dictionary[MsgText_ErrorSendingAudioFile] = MultiText{
		ru: "При отправке аудиофайла возникла ошибка, попробуйте ещё раз позже.",
		en: "There was an error sending the audio file, please try again later."}
	dictionary[MsgText_ErrorWhileSendingPicture] = MultiText{
		ru: "При отправке картинки возникла ошибка, попробуйте ещё раз позже.",
		en: "There was an error sending the picture, please try again later."}
	dictionary[MsgText_HelloCanIHelpYou] = MultiText{
		ru: "Привет! Чем могу помочь?",
		en: "Hello! How can I help?"}
	dictionary[MsgText_VoiceExamples] = MultiText{
		ru: "Примеры звучания голосов👇",
		en: "Voice examples👇"}
	dictionary[MsgText_UnexpectedError] = MultiText{
		ru: "Произошла непредвиденная ошибка, попробуйте позже.",
		en: "An unexpected error occurred, please try again later."}
	dictionary[MsgText_ResultImageGeneration] = MultiText{
		ru: `Результат генерации по запросу "%s", стиль: "%s"`,
		en: `Generation result for query "%s", style: "%s"`}
	dictionary[MsgText_ResultAudioGeneration] = MultiText{
		ru: `Результат генерации по тексту "%s", голос: "%s"`,
		en: `Generation result from text "%s", voice: "%s"`}
	dictionary[MsgText_GenerateAudioFromText] = MultiText{
		ru: "Сгенерировать аудио из текста",
		en: "Generate audio from text"}
	dictionary[MsgText_DescriptionTextNotExceed900Char] = MultiText{
		ru: "Текст описания картинки не должен превышать 900 символов.",
		en: "The description text of the picture should not exceed 900 characters."}
	dictionary[MsgText_AfterRecoveryProd] = MultiText{
		ru: "Функциональность бота восстановлена. Приносим извинения за неудобства.",
		en: "The bot's functionality has been restored. We apologize for the inconvenience."}
	dictionary[MsgText_AfterRecoveryDebug] = MultiText{
		ru: "Этот бот предназначен для тестирования и отладки, полностью рабочий и бесплатный находится здесь: @AI_free_chat_bot",
		en: "This bot is intended for testing and debugging, fully working and free, located here: @AI_free_chat_bot"}
	dictionary[MsgText_LanguageChanged] = MultiText{
		ru: "Язык успешено изменён!",
		en: "The language has been successfully changed!"}

	// buttons

	dictionary[BtnText_SendPictureWithText] = MultiText{ru: "🖼 AI Vision", en: "🖼 AI Vision"}
	dictionary[BtnText_ChangeQuestionText] = MultiText{ru: "Изменить вопрос", en: "Change question"}
	dictionary[BtnText_ChooseAnotherVoice] = MultiText{ru: "Изменить голос", en: "Change voice"}
	dictionary[BtnText_ChangeQuerryText] = MultiText{ru: "Изменить запрос", en: "Change request"}
	dictionary[BtnText_ChooseAnotherStyle] = MultiText{ru: "Изменить стиль", en: "Change style"}
	dictionary[BtnText_ChangeText] = MultiText{ru: "Изменить текст", en: "Change text"}
	dictionary[BtnText_EndDialog] = MultiText{ru: "Завершить диалог", en: "End dialog"}
	dictionary[BtnText_UploadNewImages] = MultiText{ru: "Загрузить новые фото", en: "Upload new images"}
	dictionary[BtnText_UploadNewImage] = MultiText{ru: "Загрузить новое фото", en: "Upload new image"}
	dictionary[BtnText_StartDialog] = MultiText{ru: "Начать диалог", en: "Start dialog"}
	dictionary[BtnText_GenerateAudioFromText] = MultiText{ru: "Озвучить текст", en: "Audio from text"}
	dictionary[BtnText_ClearContext] = MultiText{ru: "Очистить контекст", en: "Clear context"}
	dictionary[BtnText_Subscribe] = MultiText{ru: "✅Подписаться", en: "✅Subscribe"}

}

func GetText(key Text, lang string) string {

	element, exists := dictionary[key]
	if !exists {
		Logs <- NewLog(nil, "System", FatalError, "По ключу нет значения в словаре. Ключ:"+strconv.Itoa(int(key)))
		return "Not found"
	}

	if lang == "ru" || lang == "uk" {
		return element.ru
	} else {
		return element.en
	}

}

func textForStarting() MultiText {

	return MultiText{
		ru: `Привет, %s! 👋
		
Я бот для работы с нейросетями (v%s).
С моей помощью ты можешь использовать следующие модели:
	
<b>Gemini</b> - генерация текста и анализ изображений <i>(Google)</i>
<b>ChatGPT</b> - генерация текста, аудио и анализ изображений <i>(OpenAI)</i>
<b>Kandinsky</b> - создание изображений по текстовому описанию <i>(Sber AI)</i>
	
<u>Последние обновления:</u>
<i>09.01.24 - 🎧 добавлена генерация аудио из текста в ChatGPT.</i>
<i>14.01.24 - 🇺🇸 добавлена поддержка английского языка.</i>
<i>19.01.24 - 🏞 добавлена обработка картинок с вопросами в ChatGPT-4 Vision.</i>

Чтобы начать - просто выбери подходящую нейросеть и задай ей вопрос (или попроси сделать картинку), удачи 🔥`,

		en: `Hello, %s! 👋
		
I am a bot for working with neural networks (v%s).
With my help you can use the following models:
			
<b>Gemini</b> - text generation and image analysis <i>(Google)</i>
<b>ChatGPT</b> - text & audio generation and image analysis <i>(OpenAI)</i>
<b>Kandinsky</b> - creating images based on text description <i>(Sber AI)</i>
			
<u>Latest updates:</u>
<i>09.01.24 - 🎧 added generation of audio from text in ChatGPT.</i>
<i>14.01.24 - 🇺🇸 added English language support.</i>
<i>19.01.24 - 🏞 added processing of pictures with questions in ChatGPT-4 Vision.</i>

To get started, just choose a suitable neural network and ask it a question (or ask it to take a picture), good luck 🔥`,
	}

}
