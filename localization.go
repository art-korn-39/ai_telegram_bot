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

	MsgText_Start                     Text = iota
	MsgText_Account                        // Описание аккаунта
	MsgText_LastOperationInProgress        // Последняя операция ещё выполняется, дождитесь её завершения перед отправкой новых запросов.
	MsgText_SubscribeForUsing              // Для продолжения использования бота необходимо подписаться на канал👇
	MsgText_UnexpectedError                // Произошла непредвиденная ошибка, попробуйте позже.
	MsgText_AiNotSelected                  // Не выбрана нейросеть для обработки запроса.
	MsgText_AfterRecoveryProd              // Функциональность бота восстановлена. Приносим извинения за неудобства.
	MsgText_AfterRecoveryDebug             // Этот бот предназначен для тестирования и отладки, полностью рабочий и бесплатный находится здесь: @AI_free_chat_bot
	MsgText_HelloCanIHelpYou               // Привет! Чем могу помочь?
	MsgText_SelectOption                   // Выберите один из предложенных вариантов:
	MsgText_UnknownCommand                 // Неизвестная команда
	MsgText_EndDialog                      // Завершить диалог
	MsgText_LanguageChanged                // Язык успешено изменён!
	MsgText_DailyRequestLimitExceeded      // Достигнут дневной лимит запросов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.
	MsgText_APIdead                        // Сервис временно недоступен из-за технических неполадок :(\nПриносим изменения за неудобства.
	MsgText_AvailiableImageFormats         // Некорректный формат файла, поддерживаются изображения с расширениями: png и jpeg.
	MsgText_WrongDataType                  // Некорректный тип данных.
	MsgText_ProcessingRequest              // Обработка запроса...
	MsgText_nil

	// GEMINI

	MsgText_GeminiHello                 // Вас приветствует Gemini Pro от компании Google 🚀
	MsgText_WriteQuestionToImages       // Напишите свой вопрос к загруженным изображениям
	MsgText_UploadImages                // Загрузите одну или несколько картинок
	MsgText_PhotosUploadedWriteQuestion // Загружено фотографий: %d\nНапишите свой вопрос.\nНапример:\n\"Пришли текст с картинки\"\n\"Переведи на русский\"
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
	MsgText_DailyTokenLimitExceeded               // Достигнут дневной лимит токенов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.
	MsgText_ErrorWhileProcessingRequest           // Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже.
	MsgText_WriteQuestionToImage                  // Напишите свой вопрос к загруженному изображению.
	MsgText_UploadImage                           // Загрузите картинку
	MsgText_PhotoUploadedWriteQuestion            // Напишите свой запрос.\nНапример:\n\"Реши тест на картинке\"\n\"Как называется это блюдо?\"
	MsgText_LoadingImage                          // Выполняется загрузка изображения...

	// KANDINSKY

	MsgText_EnterDescriptionOfPicture       // Введите описание картинки:
	MsgText_DescriptionTextNotExceed900Char // Текст описания картинки не должен превышать 900 символов.
	MsgText_SelectStyleForImage             // Выберите стиль, в котором генерировать изображение.
	MsgText_SelectStyleFromOptions          // Выберите стиль из предложенных вариантов.
	MsgText_ImageGenerationStarted1         // Запущена генерация картинки, среднее время выполнения 30-40 секунд.
	MsgText_ResultImageGeneration           // Результат генерации по запросу "%s", стиль: "%s"
	MsgText_ErrorWhileSendingPicture        // При отправке картинки возникла ошибка, попробуйте ещё раз позже.
	MsgText_FailedGenerateImage1            // Не удалось сгенерировать изображение, попробуйте позже.
	MsgText_FailedGenerateImage2            // Не удалось сгенерировать изображение. Попробуйте изменить текст описания картинки.

	// SDXL

	MsgText_SDXLinfo                         // Осталось генераций и улучшений: <b>%d</b> <i>(обновится через: %d ч. %d мин.)</i>
	MsgText_DescriptionTextNotExceed2000Char // Текст описания картинки не должен превышать 2000 символов.
	MsgText_ErrorTranslatingIntoEnglish      // Возникла ошибка при переводе на английский язык, попробуйте изменить текст запроса.
	MsgText_ImageGenerationStarted2          // Запущена генерация картинки...
	MsgText_ImageProcessingStarted           // Запущена обработка картинки...
	MsgText_NoImageFoundToProcess            // Не найдена картинка для обработки.
	MsgText_FailedImageUpscale               // Не удалось повысить качество картинки, попробуйте другое изображение.
	MsgText_UploadImage2                     // Загрузите картинку (рекомендуется с разрешением не больше 1024х1024)

	// FACESWAP
	MsgText_FSinfo      // Осталось генераций: <b>%d</b> <i>(обновится через: %d ч. %d мин.)</i>
	MsgText_FSimage1    // Загрузите картинку из которой необходимо взять лицо.
	MsgText_FSimage2    // Загрузите картинку в которой нужно заменить лицо на отправленное ранее.
	MsgText_NoFaceFound // Не обнаружено лицо на фотографии

	// BAD REQUEST

	MsgText_BadRequest1 // Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса или использовать другие изображения.
	MsgText_BadRequest2 // Не удалось получить ответ от сервиса. Попробуйте изменить текст запроса.
	MsgText_BadRequest3 // Не удалось получить ответ от сервиса. Попробуйте изменить текст вопроса или начать новый диалог.
	MsgText_BadRequest4 // Запрос был заблокирован по соображениям безопасности. Попробуйте изменить текст запроса.

	// BUTTONS

	BtnText_Gemini    // 🚀 Gemini
	BtnText_ChatGPT   // 🤖 ChatGPT
	BtnText_Kandinsky // 🗿 Kandinsky
	BtnText_SDXL      // 🏔 SDXL 1.0
	BtnText_Faceswap  // 🎭 Face Swap

	BtnText_Subscribe             // ✅ Подписаться
	BtnText_SendPictureWithText   // 🖼 AI Vision
	BtnText_ChooseAnotherVoice    // Изменить голос
	BtnText_ChangeQuerryText      // 🎮 Изменить запрос
	BtnText_ChooseAnotherStyle    // 🎨 Изменить стиль
	BtnText_ChangeText            // 📝 Изменить текст
	BtnText_UploadNewImages       // Загрузить новые фото
	BtnText_UploadNewImage        // Загрузить новое фото
	BtnText_EndDialog             // 🏁 Завершить диалог
	BtnText_StartDialog           // 💭 Начать диалог
	BtnText_GenerateAudioFromText // 🗣 Озвучить текст
	BtnText_ClearContext          // 🧻 Очистить контекст
	BtnText_Upscale               // ⭐️ Улучшить (SDXL)
	BtnText_Upscale2              // ⭐️ Улучшить мою картинку
	BtnText_GenerateImage         // 🏞 Создать картинку

	//BtnText_ChangeQuestionText    // Изменить вопрос
)

func init() {

	// common
	dictionary[MsgText_Start] = textForStarting()
	dictionary[MsgText_Account] = textForAccount()
	dictionary[MsgText_nil] = MultiText{ru: "", en: ""}

	dictionary[MsgText_ChatGPTHello] = MultiText{
		ru: "Вас приветствует ChatGPT 3.5 Turbo 🤖\n\nТекущий остаток токенов: <b>%d</b> <i>(обновится через: %d ч. %d мин.)</i>",
		en: "Welcome to ChatGPT 3.5 Turbo 🤖\n\nCurrent balance of tokens: <b>%d</b> <i>(updated in: %d hours %d min.)</i>"}
	dictionary[MsgText_GeminiHello] = MultiText{
		ru: "Вас приветствует Gemini Pro от компании Google 🚀",
		en: "Welcome to Gemini Pro from Google 🚀"}
	dictionary[MsgText_EnterDescriptionOfPicture] = MultiText{
		ru: "Введите описание картинки:",
		en: "Enter a description of the picture:"}
	dictionary[MsgText_EnterTextForAudio] = MultiText{
		ru: "Введите текст для аудио:",
		en: "Enter text for audio:"}
	dictionary[MsgText_ErrorWhileProcessingRequest] = MultiText{
		ru: "Во время обработки запроса произошла ошибка. Пожалуйста, попробуйте ещё раз позже.",
		en: "An error occurred while processing the request. Please try again later."}
	dictionary[MsgText_ErrorTranslatingIntoEnglish] = MultiText{
		ru: "Возникла ошибка при переводе на английский язык, попробуйте изменить текст запроса.",
		en: "There was an error translating into English, try changing the text of the request."}
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
	dictionary[MsgText_DailyRequestLimitExceeded] = MultiText{
		ru: "Достигнут дневной лимит запросов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.",
		en: "The daily request limit has been exceeded, wait until the limit is updated (%d hours %d min.) or use another neural network."}
	dictionary[MsgText_DailyTokenLimitExceeded] = MultiText{
		ru: "Достигнут дневной лимит токенов, дождитесь обновления лимита (%d ч. %d мин.) или воспользуйтесь другой нейросетью.",
		en: "The daily token limit has reached, wait until the limit is updated (%d hours %d min.) or use another neural network."}
	dictionary[MsgText_EndDialog] = MultiText{
		ru: "Завершить диалог",
		en: "End dialog"}
	dictionary[MsgText_PhotosUploadedWriteQuestion] = MultiText{
		ru: "Загружено фотографий: %d\nНапишите свой вопрос.\nНапример:\n\"Напиши текст из картинки\"\n\"Переведи на русский\"",
		en: "Photos uploaded: %d\nWrite your question.\nFor example:\n\"Send text from picture\"\n\"Translate to English\""}
	dictionary[MsgText_UploadImages] = MultiText{
		ru: "Загрузите одну или несколько картинок",
		en: "Upload one or more images"}
	dictionary[MsgText_UploadImage] = MultiText{
		ru: "Загрузите картинку",
		en: "Upload image"}
	dictionary[MsgText_UploadImage2] = MultiText{
		ru: "Загрузите картинку (рекомендуется с разрешением не больше 1024х1024)",
		en: "Upload image (recommended with a resolution of no more than 1024x1024)"}
	dictionary[MsgText_FSimage1] = MultiText{
		ru: "Загрузите картинку из которой необходимо взять лицо.",
		en: "Upload a picture from which you need to take a face."}
	dictionary[MsgText_FSimage2] = MultiText{
		ru: "Загрузите картинку в которой нужно заменить лицо на отправленное ранее.",
		en: "Upload a picture in which you need to replace the face with the one sent earlier."}
	dictionary[MsgText_BadRequest4] = MultiText{
		ru: "Запрос был заблокирован по соображениям безопасности. Попробуйте изменить текст запроса.",
		en: "The request was blocked for security reasons. Try changing the request text."}
	dictionary[MsgText_ChatGPTDialogStarted] = MultiText{
		ru: `Запущен диалог с СhatGPT, чтобы очистить контекст от предыдущих сообщений - нажмите кнопку "Очистить контекст". Это позволяет сократить расход токенов.`,
		en: `A dialog has started with ChatGPT, to clear the context from previous messages - click the "Clear context" button. This allows you to reduce the consumption of tokens.`}
	dictionary[MsgText_ImageGenerationStarted1] = MultiText{
		ru: "Запущена генерация картинки, среднее время выполнения 30-40 секунд.",
		en: "Generation of the image has started, the average execution time is 30-40 seconds."}
	dictionary[MsgText_ImageGenerationStarted2] = MultiText{
		ru: "Запущена генерация картинки...",
		en: "Generation of the image has started..."}
	dictionary[MsgText_ImageProcessingStarted] = MultiText{
		ru: "Запущена обработка картинки...",
		en: "Processing of the image has started..."}
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
	dictionary[MsgText_PhotoUploadedWriteQuestion] = MultiText{
		ru: "Напишите свой запрос.\nНапример:\n\"Реши тест на картинке\"\n\"Как называется это блюдо?\"",
		en: "Write your request.\nFor example:\n\"Solve the test in the picture\"\n\"What is the name of this dish?\""}
	dictionary[MsgText_WriteTextForVoicing] = MultiText{
		ru: "Напишите текст для озвучивания:",
		en: "Write the text for voicing:"}
	dictionary[MsgText_AiNotSelected] = MultiText{
		ru: "Не выбрана нейросеть для обработки запроса.",
		en: "The neural network for processing requests has not been selected."}
	dictionary[MsgText_NoImageFoundToProcess] = MultiText{
		ru: "Не найдена картинка для обработки.",
		en: "No image found to process."}
	dictionary[MsgText_NoFaceFound] = MultiText{
		ru: "Не обнаружено лицо на фотографии.",
		en: "No face found in photo."}
	dictionary[MsgText_FailedLoadImages] = MultiText{
		ru: "Не удалось загрузить изображение, попробуйте ещё раз.",
		en: "Failed to load image, try again."}
	dictionary[MsgText_FailedImageUpscale] = MultiText{
		ru: "Не удалось повысить качество картинки, попробуйте другое изображение.",
		en: "Could not improve picture quality, try another image."}
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
	dictionary[MsgText_WrongDataType] = MultiText{
		ru: "Некорректный тип данных",
		en: "Wrong data type"}
	dictionary[MsgText_AvailiableImageFormats] = MultiText{
		ru: "Некорректный формат файла, поддерживаются изображения с расширениями: png и jpeg.",
		en: "Incorrect file format, supported images with extensions: png and jpeg."}
	dictionary[MsgText_ProcessingRequest] = MultiText{
		ru: "Обработка запроса...",
		en: "Processing request..."}
	dictionary[MsgText_SDXLinfo] = MultiText{
		ru: "Осталось генераций и улучшений: <b>%d</b> <i>(обновится через: %d ч. %d мин.)</i>",
		en: "Generations and upscales left: <b>%d</b> <i>(updated in: %d hours %d min.)</i>"}
	dictionary[MsgText_FSinfo] = MultiText{
		ru: "Осталось генераций: <b>%d</b> <i>(обновится через: %d ч. %d мин.)</i>",
		en: "Generations left: <b>%d</b> <i>(updated in: %d hours %d min.)</i>"}
	dictionary[MsgText_LastOperationInProgress] = MultiText{
		ru: "Последняя операция ещё выполняется, дождитесь её завершения перед отправкой новых запросов.",
		en: "The last operation is still in progress, please wait until it completes before sending new requests."}
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
	dictionary[MsgText_APIdead] = MultiText{
		ru: "Сервис временно недоступен из-за технических неполадок :(\nПриносим изменения за неудобства.",
		en: "The service is temporarily unavailable due to technical problems :(\nWe apologize for the inconvenience."}
	dictionary[MsgText_DescriptionTextNotExceed900Char] = MultiText{
		ru: "Текст описания картинки не должен превышать 900 символов.",
		en: "The description text of the picture should not exceed 900 characters."}
	dictionary[MsgText_DescriptionTextNotExceed2000Char] = MultiText{
		ru: "Текст описания картинки не должен превышать 2000 символов.",
		en: "The description text of the picture should not exceed 2000 characters."}
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

	dictionary[BtnText_Gemini] = MultiText{ru: "🚀 Gemini", en: "🚀 Gemini"}
	dictionary[BtnText_ChatGPT] = MultiText{ru: "🤖 ChatGPT", en: "🤖 ChatGPT"}
	dictionary[BtnText_Kandinsky] = MultiText{ru: "🗿 Kandinsky", en: "🗿 Kandinsky"}
	dictionary[BtnText_SDXL] = MultiText{ru: "🏔 Stable Diffusion XL", en: "🏔 Stable Diffusion XL"}
	dictionary[BtnText_Faceswap] = MultiText{ru: "🎭 Face Swap", en: "🎭 Face Swap"}

	dictionary[BtnText_SendPictureWithText] = MultiText{ru: "🖼 AI Vision", en: "🖼 AI Vision"}
	dictionary[BtnText_ChooseAnotherVoice] = MultiText{ru: "Изменить голос", en: "Change voice"}
	dictionary[BtnText_ChangeQuerryText] = MultiText{ru: "🎮 Изменить запрос", en: "🎮 Change request"}
	dictionary[BtnText_ChooseAnotherStyle] = MultiText{ru: "🎨 Изменить стиль", en: "🎨 Change style"}
	dictionary[BtnText_ChangeText] = MultiText{ru: "📝 Изменить текст", en: "📝 Change text"}
	dictionary[BtnText_EndDialog] = MultiText{ru: "🏁 Завершить диалог", en: "🏁 End dialog"}
	dictionary[BtnText_UploadNewImages] = MultiText{ru: "Загрузить новые фото", en: "Upload new images"}
	dictionary[BtnText_UploadNewImage] = MultiText{ru: "Загрузить новое фото", en: "Upload new image"}
	dictionary[BtnText_StartDialog] = MultiText{ru: "💭 Начать диалог", en: "💭 Start dialog"}
	dictionary[BtnText_GenerateAudioFromText] = MultiText{ru: "🗣 Озвучить текст", en: "🗣 Audio from text"}
	dictionary[BtnText_ClearContext] = MultiText{ru: "🧻 Очистить контекст", en: "🧻 Clear context"}
	dictionary[BtnText_Subscribe] = MultiText{ru: "✅ Подписаться", en: "✅ Subscribe"}
	dictionary[BtnText_Upscale] = MultiText{ru: "⭐️ Улучшить (SDXL)", en: "⭐️ Upscale (SDXL)"}
	dictionary[BtnText_Upscale2] = MultiText{ru: "⭐️ Улучшить мою картинку", en: "⭐ Upscale my picture"}
	dictionary[BtnText_GenerateImage] = MultiText{ru: "🏞 Создать картинку", en: "🏞 Create a picture"}

	//dictionary[BtnText_ChangeQuestionText] = MultiText{ru: "Изменить вопрос", en: "Change question"}

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

func GetLevelName(level UserLevel, lang string) string {

	var result string
	if level == Basic {
		if lang == "ru" || lang == "uk" {
			result = "Базовый"
		} else {
			result = "Basic"
		}
	} else if level == Advanced {
		if lang == "ru" || lang == "uk" {
			result = "Продвинутый"
		} else {
			result = "Advanced"
		}
	}

	return result

}

func textForStarting() MultiText {

	return MultiText{
		ru: `Привет, %s! 👋
		
Я бот для работы с нейросетями.
С моей помощью ты можешь использовать следующие модели:
	
🚀 <b>Gemini</b> - генерация текста и анализ изображений <i>(Google)</i>
🤖 <b>ChatGPT</b> - генерация текста, аудио и анализ изображений <i>(OpenAI)</i>
🗿 <b>Kandinsky</b> - создание изображений по текстовому описанию <i>(Sber AI)</i>
🏔 <b>Stable Diffusion XL</b> - создание изображений по текстовому описанию
🎭 <b>Face Swap</b> - замена лица у фотографий
	
<u>Последние обновления:</u>
<i>21.02.24 - добавлена опция по улучшению собственных картинок при помощи Stable Diffusion.</i>
<i>10.03.24 - добавлена замена лица (Face Swap).</i>

Бот полностью бесплатный, удачных генераций 🔥`,

		en: `Hello, %s! 👋
		
I am a bot for working with neural networks.
With my help you can use the following models:
			
🚀 <b>Gemini</b> - text generation and image analysis <i>(Google)</i>
🤖 <b>ChatGPT</b> - text & audio generation and image analysis <i>(OpenAI)</i>
🗿 <b>Kandinsky</b> - creating images based on text description <i>(Sber AI)</i>
🏔 <b>Stable Diffusion XL</b> - creating images based on text description
🎭 <b>Face Swap</b> - face replacement for photos
			
<u>Latest updates:</u>
<i>21.02.24 - added an option to improve your own pictures using Stable Diffusion.</i>
<i>10.03.24 - added face swap.</i>

Bot is absolutely free, successful generations 🔥`,
	}

}

func textForAccount() MultiText {

	return MultiText{
		ru: `
👤 ID Пользователя: <b>%d</b>
⭐️ Уровень: <b>%s</b>
✌️ Посещений подряд (дней): <b>%d</b>
✅ Дата первого использования: <b>%s</b>
----------------------------------------------
Дневные лимиты:     
🚀 Gemini запросы: <b>%d</b> (осталось <b>%d</b>)
🤖 ChatGPT токены: <b>%d</b> (осталось <b>%d</b>)
🗿 Kandinsky: <b>без ограничений</b>
🏔 Stable Diffusion: <b>%d</b> (осталось <b>%d</b>)
🎭 Face Swap: <b>%d</b> (осталось <b>%d</b>)
----------------------------------------------                
		
<i>Лимиты обновятся через : %d ч. %d мин.</i>
			
Регулярные пользователи бота (%d дней подряд и более) получают <b>%s</b> уровень, на котором доступно по <b>%d</b> генераций в Stable Diffusion и Face Swap + <b>%d</b> токенов ChatGPT в сутки 🔥`,

		en: `
👤 User ID: <b>%d</b>
⭐️ Level: <b>%s</b>
✌️ Consecutive visits (days): <b>%d</b>
✅ Date of first use: <b>%s</b>
----------------------------------------------
Daily limits:
🚀 Gemini requests: <b>%d</b> (<b>%d</b> left)
🤖 ChatGPT tokens: <b>%d</b> (<b>%d</b> left)
🗿 Kandinsky: <b>no limits</b>
🏔 Stable Diffusion: <b>%d</b> (<b>%d</b> left)
🎭 Face Swap: <b>%d</b> (<b>%d</b> left)
----------------------------------------------
		
<i>Limits will be updated in: %d hours %d minutes</i>
		
Regular users of the bot (%d days in a row or more) receive the <b>%s</b> level at which <b>%d</b> generation is available in Stable Diffusion and Face Swap + <b>%d</b> ChatGPT tokens per day 🔥`,
	}

}

// не используется
func textForAccount_tmp() MultiText {

	return MultiText{
		ru: `
👤 ID Пользователя: <b>%d</b>
⭐️ Уровень: <b>%s</b>
✌️ Посещений подряд (дней): <b>%d</b>
✅ Дата первого использования: <b>%s</b>
----------------------------------------------
Дневные лимиты:     
🚀 Gemini запросы: <b>%d</b> (осталось <b>%d</b>)
🤖 ChatGPT токены: <b>%d</b> (осталось <b>%d</b>)
🗿 Kandinsky: <b>без ограничений</b>
🏔 Stable Diffusion: <b>%d</b> (осталось <b>%d</b>)
----------------------------------------------                
		
<i>Лимиты обновятся через : %d ч. %d мин.</i>`,

		en: `
👤 User ID: <b>%d</b>
⭐️ Level: <b>%s</b>
✌️ Consecutive visits (days): <b>%d</b>
✅ Date of first use: <b>%s</b>
----------------------------------------------
Daily limits:
🚀 Gemini requests: <b>%d</b> (<b>%d</b> left)
🤖 ChatGPT tokens: <b>%d</b> (<b>%d</b> left)
🗿 Kandinsky: <b>no limits</b>
🏔 Stable Diffusion: <b>%d</b> (<b>%d</b> left)
----------------------------------------------
		
<i>Limits will be updated in: %d hours %d minutes</i>`,
	}

}
