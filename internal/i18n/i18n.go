package i18n

// Messages содержит все переводы приложения
type Messages struct {
	// Общие
	AppName        string
	AppNameRU      string
	SettingsTitle  string
	Save           string
	Cancel         string
	Confirm        string
	Delete         string
	Close          string
	
	// Вкладки
	TabGeneral     string
	TabNotifications string
	TabLogs        string
	
	// Горячие клавиши
	SectionHotkeys     string
	LabelHotkeyStart   string
	LabelHotkeyStop    string
	LabelHotkeyEditor  string
	
	// Lemonade Server
	SectionLemonade    string
	LabelURL           string
	LabelModel         string
	LabelLanguage      string
	ModelWhisperTurbo  string
	LangRussian        string
	LangEnglish        string
	
	// Автозапуск
	SectionAutostart   string
	CheckboxAutostart  string
	HintAutostart      string
	
	// Язык интерфейса
	SectionAppLanguage string
	LabelAppLanguage   string
	
	// Уведомления
	SectionNotifications string
	CheckboxSound        string
	CheckboxSoundOnRecord string
	CheckboxToast        string
	
	// Логирование
	SectionLogs          string
	CheckboxLogging      string
	LabelLoggingLevel    string
	LabelLastLogs        string
	BtnViewLogs          string
	BtnClearLogs         string
	BtnSaveLogs          string
	ConfirmClearLogs     string
	
	// Сообщения
	MsgSettingsSaved     string
	MsgSettingsError     string
	MsgLogsEmpty         string
	MsgLogsSaved         string
	MsgLogsClearConfirm  string
	
	// Трей
	TrayStart            string
	TrayStop             string
	TrayEditor           string
	TraySettings         string
	TrayLogs             string
	TrayQuit             string
	
	// Веб-сервер
	WebServerStarted     string
	WebServerSettings    string
	WebServerFileNotFound string
	WebServerMethodNotAllowed string
	WebServerConfigSaved string
	WebServerErrorSaving string
	
	// Логгер
	LoggerInitialized    string
	
	// Хоткеи
	HotkeyParseError     string
	HotkeyRegisterError  string
	HotkeyErrorsHeader   string
	HotkeysRegistered    string
	HotkeyStartInfo      string
	HotkeyStopInfo       string
	HotkeyEditorInfo     string
	ListeningHotkeys     string
	HotkeyStartPressed   string
	HotkeyStopPressed    string
	HotkeyEditorPressed  string
	
	// Main.go
	AppStarting          string
	AppStarted           string
	LogFile              string
	TrayMenu             string
	HotkeysInfo          string
	PressCtrlC           string
	TrayStarted          string
	StartingTray         string
	RecordingFromMenu    string
	StopRecordingFromMenu string
	EditorFromMenu       string
	SettingsFromMenu     string
	LogsFromMenu         string
	ErrorStartingRec     string
	RecordingStarted     string
	RecordingSaved       string
	OpenSettingsBrowser  string
	ViewLogsPlaceholder  string
	OpenEditorPlaceholder string
	
	// Ошибки
	ErrorPrefix          string
	ErrorStopRecording   string
	ErrorOpenSettings    string
	ErrorWebServer       string
	ErrorHotkeyRegister  string
	WarningHotkeys       string
	WarningHotkeysDetail string
	ErrorAlreadyRunning  string
	ErrorAlreadyRunningWait string
	ErrorInitLogger      string
	ErrorConfigLoad      string
	InfoUsingDefaultConfig string
	ConfigHeader         string
	LemonadeURL          string
	ErrorHotkeyListener  string
	ConfigHotkeys        string
	ConfigLemonade       string
	ConfigNotifications  string
	ConfigAutostart      string
	ConfigLogging        string
	ConfigLanguage       string
	RestartRequired      string
	RestartRequiredText  string
	LanguageChanged      string
	RestartingTray       string
	LanguageSwitched     string
	ErrorPlaySound       string
	ErrorShowToast       string
	ToastTitle           string
	ErrorDeleteFile      string
	InfoFileDeleted      string
	
	// Транскрибация
	PanicTranscription    string
	PanicPrefix           string
	StartTranscription    string
	TranscribingAudio     string
	AudioFile             string
	LemonadeNotInit       string
	ConfigNotInit         string
	ModelInfo             string
	AudioFileNotFound     string
	TranscriptionError    string
	TranscriptionComplete string
	ProcessTime           string
	SpeedInfo             string
	Characters            string
	Backend               string
	TextCopied            string
	TranscriptionStats    string
	CopyError             string
	CopyFailed            string
	AudioDuration         string
	AudioDurationLabel    string
	
	// Tray tooltip
	TrayTooltipRU         string
	TrayTooltipEN         string
	
	// Tray status
	StatusIdle            string
	StatusRecording       string
	StatusProcessing      string

	// Автозапуск (логирование)
	AutostartEnabled      string
	AutostartDisabled     string

	// Редактор (логирование)
	EditorTextSet         string
	EditorAddedToHistory  string
	EditorCopiedToClip    string
	EditorShowWindow      string
	EditorCleared         string

	// Редактор (GUI окно)
	EditorWindowTitle     string
	EditorCharCount       string
	EditorCopyAndClose    string
	EditorPrev            string
	EditorNext            string
	EditorHistoryPos      string
	EditorNoText          string
	EditorAlreadyOpen     string
	EditorProcessError    string
	EditorDelete          string

	// Хоткеи
	HotkeysUnregistered   string

	// Логгер
	LogsCleared           string

	// Окно настроек (Fyne)
	SettingsAlreadyOpen   string
	SettingsProcessError  string
	SettingsOpening       string

	// Окно логов (Fyne)
	LogsWindowTitle       string
	LogsAlreadyOpen       string
	LogsProcessError      string
	LogsOpening           string
	LogsFilterAll         string
	LogsRefresh           string
	LogsClear             string
	LogsLineCount         string
	LogsClearConfirm      string
	LogsEmpty             string
	LogsFileError         string

	// История
	SectionHistory        string
	LabelHistorySize      string

	// Консоль
	CheckboxShowConsole   string
	HintShowConsole       string

	// Модели
	ModelLoading          string
	ModelActivating       string
	ModelLoadError        string
	ModelLoadSuccess      string
	ModelFetchError       string
	ModelFetchEmpty       string
	ModelFetchOK          string
	ModelRefresh          string
	ModelDownloaded       string
	ModelNotDownloaded    string
	ModelSizeGB           string
	ModelInstall          string
	ModelInstalling       string
	ModelInstallProgress  string
	ModelInstallDone      string
	ModelInstallError     string
	ModelNotInstalled     string
	ModelReady            string

	// Параметры транскрибации
	LabelPrompt           string
	LabelTemperature      string
	HintPrompt            string
	HintTemperature       string

	// Автовставка
	CheckboxAutoPaste     string
	HintAutoPaste         string
	SectionBehavior       string

	// Обрезка тишины и диагностика аудио
	SilenceTrimmed        string
	RecordingSilent       string
	AudioStats            string
	AudioSilentWarning    string

	// Бэкенд транскрибации
	SectionBackend        string
	LabelBackend          string
	BackendLemonade       string
	BackendWhisperAPI     string
	HintBackendLemonade   string
	HintBackendWhisperAPI string

	// Whisper API
	TabWhisperAPI         string
	SectionWhisperAPI     string
	WhisperAPIHintURL     string
	WhisperAPIStatus      string
	WhisperAPIStatusOK    string
	WhisperAPIStatusError string
	WhisperAPICheckBtn    string
	WhisperAPINotInit     string
	BackendInfo           string
	ConfigWhisperAPI      string
	ConfigBackend         string

	// FastFlowLM
	BackendFastFlowLM       string
	TabFastFlowLM           string
	SectionFastFlowLM       string
	HintBackendFastFlowLM   string
	FastFlowLMHintURL       string
	FastFlowLMStatus        string
	FastFlowLMStatusOK      string
	FastFlowLMStatusError   string
	FastFlowLMCheckBtn      string
	FastFlowLMNotInit       string
	FastFlowLMLLMModel      string
	FastFlowLMHintLLMModel  string
	FastFlowLMNotInstalled  string
	FastFlowLMStarting      string
	FastFlowLMStarted       string
	FastFlowLMStartError    string
	ConfigFastFlowLM        string

	// Ошибки (внутренние, для логов)
	ErrorAudioDuration    string
	ErrorAutostartEnable  string
	ErrorAutostartDisable string
	ErrorReadStdin        string
	ErrorParseJSON        string
}

// Get возвращает переводы для указанного языка
func Get(lang string) *Messages {
	switch lang {
	case "ru":
		return &Messages{
			// Общие
			AppName:        "Voice Input Go",
			AppNameRU:      "Голосовой ввод",
			SettingsTitle:  "⚙️ Настройки — Голосовой ввод",
			Save:           "Сохранить",
			Cancel:         "Отмена",
			Confirm:        "Подтвердить",
			Delete:         "Удалить",
			Close:          "Закрыть",
			
			// Вкладки
			TabGeneral:     "Основные",
			TabNotifications: "Уведомления",
			TabLogs:        "Логи",
			
			// Горячие клавиши
			SectionHotkeys:     "Горячие клавиши",
			LabelHotkeyStart:   "Начать запись:",
			LabelHotkeyStop:    "Остановить:",
			LabelHotkeyEditor:  "Редактор:",
			
			// Lemonade Server
			SectionLemonade:    "Lemonade Server",
			LabelURL:           "URL:",
			LabelModel:         "Модель:",
			LabelLanguage:      "Язык:",
			ModelWhisperTurbo:  "Whisper-Large-v3-Turbo",
			LangRussian:        "Русский",
			LangEnglish:        "English",
			
			// Автозапуск
			SectionAutostart:   "Автозапуск",
			CheckboxAutostart:  "Запускать при старте системы",
			HintAutostart:      "Приложение будет запускаться автоматически при входе в систему",
			
			// Язык интерфейса
			SectionAppLanguage: "Язык интерфейса",
			LabelAppLanguage:   "Язык приложения:",
			
			// Уведомления
			SectionNotifications: "Уведомления",
			CheckboxSound:        "Звук после распознавания",
			CheckboxSoundOnRecord: "Звук при начале записи",
			CheckboxToast:        "Всплывающее уведомление",
			
			// Логирование
			SectionLogs:          "Логирование",
			CheckboxLogging:      "Включить логирование",
			LabelLoggingLevel:    "Уровень:",
			LabelLastLogs:        "Последние логи:",
			BtnViewLogs:          "🔄 Обновить логи",
			BtnClearLogs:         "🗑 Очистить",
			BtnSaveLogs:          "💾 Сохранить в файл",
			ConfirmClearLogs:     "Вы уверены, что хотите очистить логи?",
			
			// Сообщения
			MsgSettingsSaved:     "Настройки сохранены!",
			MsgSettingsError:     "Ошибка при сохранении настроек: ",
			MsgLogsEmpty:         "Логи пусты",
			MsgLogsSaved:         "Логи сохранены!",
			MsgLogsClearConfirm:  "Вы уверены, что хотите очистить логи?",
			
			// Трей
			TrayStart:        "Начать запись",
			TrayStop:         "Остановить запись",
			TrayEditor:       "Редактор",
			TraySettings:     "Настройки",
			TrayLogs:         "Логи",
			TrayQuit:         "Выход",
			
			// Веб-сервер
			WebServerStarted:     "🌐 Веб-сервер запущен на http://localhost:%d",
			WebServerSettings:    "   Настройки: http://localhost:%d/settings",
			WebServerFileNotFound: "Файл не найден",
			WebServerMethodNotAllowed: "Метод не разрешён",
			WebServerConfigSaved: "Конфигурация сохранена",
			WebServerErrorSaving: "Ошибка сохранения конфигурации",
			
			// Логгер
			LoggerInitialized:    "Логгер инициализирован (уровень: %s)",
			
			// Хоткеи
			HotkeyParseError:     "ошибка парсинга комбинации '%s': %v",
			HotkeyRegisterError:  "ошибка регистрации хоткея '%s': %v (возможно занят)",
			HotkeyErrorsHeader:   "ошибки регистрации хоткеев",
			HotkeysRegistered:    "✅ Глобальные хоткеи зарегистрированы:",
			HotkeyStartInfo:      "   %s - Начать запись\n",
			HotkeyStopInfo:       "   %s - Остановить запись\n",
			HotkeyEditorInfo:     "   %s - Открыть редактор\n",
			ListeningHotkeys:     "🎧 Прослушивание глобальных хоткеев...",
			HotkeyStartPressed:   "🔥 Хоткей нажат: старт\n",
			HotkeyStopPressed:    "🔥 Хоткей нажат: стоп\n",
			HotkeyEditorPressed:  "🔥 Хоткей нажат: редактор\n",
			
			// Main.go
			AppStarting:          "Voice Input Go запускается...",
			AppStarted:           "✅ Voice Input Go запущен!",
			LogFile:              "📝 Лог файл: %s\n",
			TrayMenu:             "🖱️  Меню в системном трее\n",
			HotkeysInfo:          "⌨️  Горячие клавиши:\n   %s - Начать запись\n   %s - Остановить запись\n   %s - Открыть редактор\n",
			PressCtrlC:           "Нажми Ctrl+C для выхода...",
			TrayStarted:          "✅ Системный трей запущен с меню\n",
			StartingTray:         "Запуск системного трея...",
			RecordingFromMenu:    "🎤 Запись из меню",
			StopRecordingFromMenu: "⏹ Остановка записи из меню",
			EditorFromMenu:       "📝 Редактор из меню",
			SettingsFromMenu:     "⚙️ Настройки из меню",
			LogsFromMenu:         "📊 Логи из меню",
			ErrorStartingRec:     "Ошибка запуска записи: %v",
			RecordingStarted:     "✅ Запись начата",
			RecordingSaved:       "✅ Запись сохранена: %s\n",
			OpenSettingsBrowser:  "⚙️ Открыть настройки в браузере",
			ViewLogsPlaceholder:  "📊 Просмотр логов (заглушка)",
			OpenEditorPlaceholder: "📝 Открыть редактор (заглушка)",
			
			// Ошибки
			ErrorPrefix:          "❌ Ошибка: %v\n",
			ErrorStopRecording:   "Ошибка остановки записи: %v",
			ErrorOpenSettings:    "Не удалось открыть настройки: %v",
			ErrorWebServer:       "Ошибка веб-сервера: %v",
			ErrorHotkeyRegister:  "Ошибка регистрации хоткеев: %v",
			WarningHotkeys:       "⚠️  Предупреждение: %v\n",
			WarningHotkeysDetail: "Некоторые хоткеи могут быть заняты другим приложением.\n   Закройте конфликтующее приложение или измените хоткеи в конфиге.",
			ErrorAlreadyRunning:  "Приложение уже запущено",
			ErrorAlreadyRunningWait: "Если приложение закрылось некорректно, подождите несколько секунд.\n   If the application closed incorrectly, wait a few seconds.",
			ErrorInitLogger:      "Не удалось инициализировать логгер: %v",
			ErrorConfigLoad:      "Ошибка загрузки конфига: %v",
			InfoUsingDefaultConfig: "Используется конфиг по умолчанию",
			ConfigHeader:         "📋 Конфигурация:",
			LemonadeURL:          "Lemonade URL: %s",
			ErrorHotkeyListener:  "Ошибка запуска прослушивания хоткеев: %v",
			ConfigHotkeys:        "   🔑 Хоткеи: Start=%s, Stop=%s, Editor=%s",
			ConfigLemonade:       "   🤖 Lemonade: URL=%s, Model=%s, Language=%s",
			ConfigNotifications:  "   🔔 Уведомления: Звук=%v, Всплывашки=%v",
			ConfigAutostart:      "   🚀 Автозапуск: %v",
			ConfigLogging:        "   📝 Логирование: Включено=%v, Уровень=%s",
			ConfigLanguage:       "   🌐 Язык интерфейса: %s",
			RestartRequired:      "⚠️ Требуется перезапуск",
			RestartRequiredText:  "Для применения языка необходимо перезапустить приложение.\n\nПерезапустить сейчас?",
			LanguageChanged:      "Язык изменён: %s → %s",
			RestartingTray:       "Перезапуск трея с новым языком...",
			LanguageSwitched:     "✅ Язык переключён на: %s",
			ErrorPlaySound:       "Ошибка воспроизведения звука: %v",
			ErrorShowToast:       "Ошибка показа уведомления: %v",
			ToastTitle:           "Голосовой ввод",
			ErrorDeleteFile:      "Ошибка удаления временного файла: %v",
			InfoFileDeleted:      "Временный файл удалён: %s",
			
			// Транскрибация
			PanicTranscription:    "ПАНИКА при транскрибации: %v",
			PanicPrefix:           "❌ ПАНИКА: %v\n",
			StartTranscription:    "🎤 Начало транскрибации...",
			TranscribingAudio:     "🔄 Транскрибация аудио...",
			AudioFile:             "Аудио файл: %s",
			LemonadeNotInit:       "Lemonade клиент не инициализирован!",
			ConfigNotInit:         "Конфигурация не инициализирована!",
			ModelInfo:             "Модель: %s, Язык: %s",
			AudioFileNotFound:     "Аудио файл не найден: %s",
			TranscriptionError:    "Ошибка транскрибации: %v",
			TranscriptionComplete: "✅ Транскрибация завершена",
			ProcessTime:           "🕐 Время обработки: %.2f сек",
			SpeedInfo:             "🚀 Скорость: %.2fx реального времени",
			Characters:            "📝 Символов: %d\n",
			Backend:               "💻 Бэкенд: %s\n",
			TextCopied:            "📋 Текст скопирован в буфер обмена!\n",
			TranscriptionStats:    "📊 Транскрибация | Обработка: %.2f сек | Скорость: %.2fx | Символов: %d",
			CopyError:             "Ошибка копирования в буфер: %v",
			CopyFailed:            "❌ Ошибка копирования: %v\n",
			AudioDuration:         "Длительность: %.2f сек",
			AudioDurationLabel:    "🎵 Длительность: %.2f сек",
			
			// Tray tooltip
			TrayTooltipRU:         "Voice Input Go - голосовая транскрибация",
			TrayTooltipEN:         "Voice Input Go - voice transcription",
			
			// Tray status
			StatusIdle:       "🟢 Ожидание",
			StatusRecording:  "🔴 Запись",
			StatusProcessing: "🟡 Обработка",

			// Автозапуск
			AutostartEnabled:  "Автозапуск включён (%s)",
			AutostartDisabled: "Автозапуск отключён (%s)",

			// Редактор
			EditorTextSet:        "Текст установлен (%d символов)",
			EditorAddedToHistory: "Добавлено в историю (всего: %d)",
			EditorCopiedToClip:   "Скопировано в буфер (%d символов)",
			EditorShowWindow:     "Открытие окна редактора",
			EditorCleared:        "Редактор очищен",

			// Редактор (GUI)
			EditorWindowTitle:    "Голосовой ввод - Редактор",
			EditorCharCount:      "Символов: %d",
			EditorCopyAndClose:   "Копировать и закрыть",
			EditorPrev:           "< Назад",
			EditorNext:           "Далее >",
			EditorHistoryPos:     "%d из %d",
			EditorNoText:         "Нет текста",
			EditorAlreadyOpen:    "Редактор уже открыт",
			EditorProcessError:   "Ошибка открытия редактора: %v",
			EditorDelete:         "Удалить",

			// Хоткеи
			HotkeysUnregistered: "👋 Хоткеи отменены",

			// Логгер
			LogsCleared: "Логи очищены",

			// Окно настроек (Fyne)
			SettingsAlreadyOpen:  "Настройки уже открыты",
			SettingsProcessError: "Ошибка открытия настроек: %v",
			SettingsOpening:      "Открытие окна настроек",

			// Окно логов (Fyne)
			LogsWindowTitle:  "Логи — Голосовой ввод",
			LogsAlreadyOpen:  "Окно логов уже открыто",
			LogsProcessError: "Ошибка открытия логов: %v",
			LogsOpening:      "Открытие окна логов",
			LogsFilterAll:    "Все",
			LogsRefresh:      "Обновить",
			LogsClear:        "Очистить",
			LogsLineCount:    "Строк: %d",
			LogsClearConfirm: "Очистить файл логов?",
			LogsEmpty:        "Логи пусты",
			LogsFileError:    "Ошибка чтения логов: %v",

			// История
			SectionHistory:   "История",
			LabelHistorySize: "Количество записей:",

			// Консоль
			CheckboxShowConsole: "Показывать консольное окно",
			HintShowConsole:     "Для отладки и просмотра работы приложения",

			// Модели
			ModelLoading:        "Загрузка списка моделей...",
			ModelActivating:     "Загрузка модели %s...",
			ModelLoadError:      "Ошибка загрузки модели: %v",
			ModelLoadSuccess:    "Модель %s загружена",
			ModelFetchError:     "Сервер недоступен: %v",
			ModelFetchEmpty:     "Whisper модели не найдены",
			ModelFetchOK:        "Найдено моделей: %d",
			ModelRefresh:        "Обновить",
			ModelDownloaded:     "(установлена)",
			ModelNotDownloaded:  "(%.1f ГБ)",
			ModelSizeGB:         "(%.1f ГБ)",
			ModelInstall:        "Установить",
			ModelInstalling:     "Установка %s...",
			ModelInstallProgress: "Скачивание: %s — %d%%",
			ModelInstallDone:    "Модель %s установлена!",
			ModelInstallError:   "Ошибка установки: %v",
			ModelNotInstalled:   "Модель не установлена на сервере. Нажмите \"Установить\" для скачивания.",
			ModelReady:          "Модель установлена и готова к работе",

			// Параметры транскрибации
			LabelPrompt:      "Подсказка:",
			LabelTemperature: "Температура:",
			HintPrompt:       "Подсказка для модели: стиль пунктуации, имена, термины",
			HintTemperature:  "0.0 — точнее, 0.2-0.4 — разнообразнее (пунктуация, смена языка)",

			// Автовставка
			CheckboxAutoPaste: "Автовставка текста (Ctrl+V)",
			HintAutoPaste:     "После распознавания текст автоматически вставляется в позицию курсора",
			SectionBehavior:   "Поведение",

			// Обрезка тишины и диагностика аудио
			SilenceTrimmed:     "✂️ Обрезано тишины: %.1f сек\n",
			RecordingSilent:    "Запись содержит только тишину",
			AudioStats:         "🎙️ Аудио: %.1f сек, %d сэмплов, пик: %.4f, RMS: %.4f\n",
			AudioSilentWarning: "⚠️ Микрофон не захватывает звук (пик=0). Проверьте разрешение микрофона в настройках ОС.",

			// Бэкенд транскрибации
			SectionBackend:        "Бэкенд транскрибации",
			LabelBackend:          "Бэкенд:",
			BackendLemonade:       "Lemonade Server (локальный)",
			BackendWhisperAPI:     "Whisper API (внешний сервер)",
			HintBackendLemonade:   "Локальный AI-сервер с управлением моделями",
			HintBackendWhisperAPI: "Внешний Whisper сервер (Docker, сеть)",

			// Whisper API
			TabWhisperAPI:         "Whisper API",
			SectionWhisperAPI:     "Whisper API",
			WhisperAPIHintURL:     "Адрес сервера, например http://192.168.1.50:9000",
			WhisperAPIStatus:      "Статус: проверка...",
			WhisperAPIStatusOK:    "Сервер доступен",
			WhisperAPIStatusError: "Сервер недоступен: %v",
			WhisperAPICheckBtn:    "Проверить",
			WhisperAPINotInit:     "Whisper API клиент не инициализирован!",
			BackendInfo:           "Бэкенд: %s",
			ConfigWhisperAPI:      "   🌐 Whisper API: URL=%s, Language=%s",
			ConfigBackend:         "   🔀 Бэкенд: %s",

			// FastFlowLM
			BackendFastFlowLM:       "FastFlowLM (NPU, локальный)",
			TabFastFlowLM:           "FastFlowLM",
			SectionFastFlowLM:       "FastFlowLM",
			HintBackendFastFlowLM:   "AMD Ryzen AI NPU — локальная транскрибация",
			FastFlowLMHintURL:       "Адрес сервера, по умолчанию http://localhost:52625",
			FastFlowLMStatus:        "Статус: проверка...",
			FastFlowLMStatusOK:      "Сервер доступен",
			FastFlowLMStatusError:   "Сервер недоступен: %v",
			FastFlowLMCheckBtn:      "Проверить",
			FastFlowLMNotInit:       "FastFlowLM клиент не инициализирован!",
			FastFlowLMLLMModel:      "LLM модель:",
			FastFlowLMHintLLMModel:  "LLM модель для запуска сервера (напр. llama3.2:1b)",
			FastFlowLMNotInstalled:  "FLM не найден в PATH. Установите FastFlowLM: https://github.com/amd/FastFlowLM",
			FastFlowLMStarting:      "Запуск FLM-сервера...",
			FastFlowLMStarted:       "FLM-сервер запущен",
			FastFlowLMStartError:    "Ошибка запуска FLM-сервера: %v",
			ConfigFastFlowLM:        "   ⚡ FastFlowLM: URL=%s, Model=%s, Language=%s",

			// Ошибки (внутренние, для логов)
			ErrorAudioDuration:    "Не удалось получить длительность аудио: %v",
			ErrorAutostartEnable:  "Ошибка включения автозапуска: %v",
			ErrorAutostartDisable: "Ошибка отключения автозапуска: %v",
			ErrorReadStdin:        "Ошибка чтения stdin: %v\n",
			ErrorParseJSON:        "Ошибка разбора JSON: %v\n",
		}

	case "en":
		return &Messages{
			// Общие
			AppName:        "Voice Input Go",
			AppNameRU:      "Voice Input",
			SettingsTitle:  "⚙️ Settings — Voice Input",
			Save:           "Save",
			Cancel:         "Cancel",
			Confirm:        "Confirm",
			Delete:         "Delete",
			Close:          "Close",
			
			// Вкладки
			TabGeneral:     "General",
			TabNotifications: "Notifications",
			TabLogs:        "Logs",
			
			// Горячие клавиши
			SectionHotkeys:     "Hotkeys",
			LabelHotkeyStart:   "Start recording:",
			LabelHotkeyStop:    "Stop:",
			LabelHotkeyEditor:  "Editor:",
			
			// Lemonade Server
			SectionLemonade:    "Lemonade Server",
			LabelURL:           "URL:",
			LabelModel:         "Model:",
			LabelLanguage:      "Language:",
			ModelWhisperTurbo:  "Whisper-Large-v3-Turbo",
			LangRussian:        "Русский",
			LangEnglish:        "English",
			
			// Автозапуск
			SectionAutostart:   "Autostart",
			CheckboxAutostart:  "Launch at system startup",
			HintAutostart:      "The application will start automatically when you log in",
			
			// Язык интерфейса
			SectionAppLanguage: "Interface Language",
			LabelAppLanguage:   "Application language:",
			
			// Уведомления
			SectionNotifications: "Notifications",
			CheckboxSound:        "Sound after transcription",
			CheckboxSoundOnRecord: "Sound on recording start",
			CheckboxToast:        "Popup notification",
			
			// Логирование
			SectionLogs:          "Logging",
			CheckboxLogging:      "Enable logging",
			LabelLoggingLevel:    "Level:",
			LabelLastLogs:        "Recent logs:",
			BtnViewLogs:          "🔄 Refresh logs",
			BtnClearLogs:         "🗑 Clear",
			BtnSaveLogs:          "💾 Save to file",
			ConfirmClearLogs:     "Are you sure you want to clear logs?",
			
			// Сообщения
			MsgSettingsSaved:     "Settings saved!",
			MsgSettingsError:     "Error saving settings: ",
			MsgLogsEmpty:         "Logs are empty",
			MsgLogsSaved:         "Logs saved!",
			MsgLogsClearConfirm:  "Are you sure you want to clear logs?",
			
			// Трей
			TrayStart:        "Start recording",
			TrayStop:         "Stop recording",
			TrayEditor:       "Editor",
			TraySettings:     "Settings",
			TrayLogs:         "Logs",
			TrayQuit:         "Quit",
			
			// Веб-сервер
			WebServerStarted:     "🌐 Web server started on http://localhost:%d",
			WebServerSettings:    "   Settings: http://localhost:%d/settings",
			WebServerFileNotFound: "File not found",
			WebServerMethodNotAllowed: "Method not allowed",
			WebServerConfigSaved: "Configuration saved",
			WebServerErrorSaving: "Error saving configuration",
			
			// Логгер
			LoggerInitialized:    "Logger initialized (level: %s)",
			
			// Хоткеи
			HotkeyParseError:     "failed to parse combo '%s': %v",
			HotkeyRegisterError:  "failed to register hotkey '%s': %v (may be in use)",
			HotkeyErrorsHeader:   "hotkey registration errors",
			HotkeysRegistered:    "✅ Global hotkeys registered:",
			HotkeyStartInfo:      "   %s - Start recording\n",
			HotkeyStopInfo:       "   %s - Stop recording\n",
			HotkeyEditorInfo:     "   %s - Open editor\n",
			ListeningHotkeys:     "🎧 Listening for global hotkeys...",
			HotkeyStartPressed:   "🔥 Hotkey pressed: start\n",
			HotkeyStopPressed:    "🔥 Hotkey pressed: stop\n",
			HotkeyEditorPressed:  "🔥 Hotkey pressed: editor\n",
			
			// Main.go
			AppStarting:          "Voice Input Go starting...",
			AppStarted:           "✅ Voice Input Go started!",
			LogFile:              "📝 Log file: %s\n",
			TrayMenu:             "🖱️  Check system tray for menu\n",
			HotkeysInfo:          "⌨️  Hotkeys:\n   %s - Start recording\n   %s - Stop recording\n   %s - Open editor\n",
			PressCtrlC:           "Press Ctrl+C to exit...",
			TrayStarted:          "✅ System tray started with menu\n",
			StartingTray:         "Starting system tray...",
			RecordingFromMenu:    "🎤 Recording from menu",
			StopRecordingFromMenu: "⏹ Stop recording from menu",
			EditorFromMenu:       "📝 Editor from menu",
			SettingsFromMenu:     "⚙️ Settings from menu",
			LogsFromMenu:         "📊 Logs from menu",
			ErrorStartingRec:     "Failed to start recording: %v",
			RecordingStarted:     "✅ Recording started",
			RecordingSaved:       "✅ Recording saved: %s\n",
			OpenSettingsBrowser:  "⚙️ Open settings in browser",
			ViewLogsPlaceholder:  "📊 View logs (placeholder)",
			OpenEditorPlaceholder: "📝 Open editor (placeholder)",
			
			// Ошибки
			ErrorPrefix:          "❌ Error: %v\n",
			ErrorStopRecording:   "Failed to stop recording: %v",
			ErrorOpenSettings:    "Failed to open settings: %v",
			ErrorWebServer:       "Web server error: %v",
			ErrorHotkeyRegister:  "Hotkey registration error: %v",
			WarningHotkeys:       "⚠️  WARNING: %v\n",
			WarningHotkeysDetail: "Some hotkeys may be already in use by another application.\n   Close the conflicting application or change hotkeys in config.",
			ErrorAlreadyRunning:  "Application is already running",
			ErrorAlreadyRunningWait: "If the application closed incorrectly, wait a few seconds.",
			ErrorInitLogger:      "Failed to initialize logger: %v",
			ErrorConfigLoad:      "Config load error: %v",
			InfoUsingDefaultConfig: "Using default config",
			ConfigHeader:         "📋 Configuration:",
			LemonadeURL:          "Lemonade URL: %s",
			ErrorHotkeyListener:  "Failed to start hotkey listener: %v",
			ConfigHotkeys:        "   🔑 Hotkeys: Start=%s, Stop=%s, Editor=%s",
			ConfigLemonade:       "   🤖 Lemonade: URL=%s, Model=%s, Language=%s",
			ConfigNotifications:  "   🔔 Notifications: Sound=%v, Toast=%v",
			ConfigAutostart:      "   🚀 Autostart: %v",
			ConfigLogging:        "   📝 Logging: Enabled=%v, Level=%s",
			ConfigLanguage:       "   🌐 Interface language: %s",
			RestartRequired:      "⚠️ Restart required",
			RestartRequiredText:  "Application restart is required to apply the language.\n\nRestart now?",
			LanguageChanged:      "Language changed: %s → %s",
			RestartingTray:       "Restarting tray with new language...",
			LanguageSwitched:     "✅ Language switched to: %s",
			ErrorPlaySound:       "Failed to play notification sound: %v",
			ErrorShowToast:       "Failed to show toast notification: %v",
			ToastTitle:           "Voice Input",
			ErrorDeleteFile:      "Failed to delete temporary file: %v",
			InfoFileDeleted:      "Temporary file deleted: %s",
			
			// Транскрибация
			PanicTranscription:    "PANIC during transcription: %v",
			PanicPrefix:           "❌ PANIC: %v\n",
			StartTranscription:    "🎤 Starting transcription...",
			TranscribingAudio:     "🔄 Transcribing audio...",
			AudioFile:             "Audio file: %s",
			LemonadeNotInit:       "Lemonade client not initialized!",
			ConfigNotInit:         "Config not initialized!",
			ModelInfo:             "Model: %s, Language: %s",
			AudioFileNotFound:     "Audio file not found: %s",
			TranscriptionError:    "Transcription error: %v",
			TranscriptionComplete: "✅ Transcription complete",
			ProcessTime:           "🕐 Process time: %.2f sec",
			SpeedInfo:             "🚀 Speed: %.2fx real-time",
			Characters:            "📝 Characters: %d\n",
			Backend:               "💻 Backend: %s\n",
			TextCopied:            "📋 Text copied to clipboard!\n",
			TranscriptionStats:    "📊 Transcription | Process: %.2f sec | Speed: %.2fx | Characters: %d",
			CopyError:             "Failed to copy to clipboard: %v",
			CopyFailed:            "❌ Copy failed: %v\n",
			AudioDuration:         "Duration: %.2f sec",
			AudioDurationLabel:    "🎵 Duration: %.2f sec",
			
			// Tray tooltip
			TrayTooltipRU:         "Voice Input Go - голосовая транскрибация",
			TrayTooltipEN:         "Voice Input Go - voice transcription",
			
			// Tray status
			StatusIdle:       "🟢 Idle",
			StatusRecording:  "🔴 Recording",
			StatusProcessing: "🟡 Processing",

			// Автозапуск
			AutostartEnabled:  "Autostart enabled (%s)",
			AutostartDisabled: "Autostart disabled (%s)",

			// Редактор
			EditorTextSet:        "Text set (%d chars)",
			EditorAddedToHistory: "Added to history (total: %d)",
			EditorCopiedToClip:   "Copied to clipboard (%d chars)",
			EditorShowWindow:     "Show editor window",
			EditorCleared:        "Editor cleared",

			// Редактор (GUI)
			EditorWindowTitle:    "Voice Input - Editor",
			EditorCharCount:      "Characters: %d",
			EditorCopyAndClose:   "Copy and Close",
			EditorPrev:           "< Prev",
			EditorNext:           "Next >",
			EditorHistoryPos:     "%d of %d",
			EditorNoText:         "No text",
			EditorAlreadyOpen:    "Editor already open",
			EditorProcessError:   "Editor error: %v",
			EditorDelete:         "Delete",

			// Хоткеи
			HotkeysUnregistered: "👋 Hotkeys unregistered",

			// Логгер
			LogsCleared: "Logs cleared",

			// Окно настроек (Fyne)
			SettingsAlreadyOpen:  "Settings already open",
			SettingsProcessError: "Settings error: %v",
			SettingsOpening:      "Opening settings window",

			// Окно логов (Fyne)
			LogsWindowTitle:  "Logs — Voice Input",
			LogsAlreadyOpen:  "Log viewer already open",
			LogsProcessError: "Log viewer error: %v",
			LogsOpening:      "Opening log viewer",
			LogsFilterAll:    "All",
			LogsRefresh:      "Refresh",
			LogsClear:        "Clear",
			LogsLineCount:    "Lines: %d",
			LogsClearConfirm: "Clear log file?",
			LogsEmpty:        "Logs are empty",
			LogsFileError:    "Error reading logs: %v",

			// История
			SectionHistory:   "History",
			LabelHistorySize: "Number of entries:",

			// Консоль
			CheckboxShowConsole: "Show console window",
			HintShowConsole:     "For debugging and monitoring app activity",

			// Модели
			ModelLoading:        "Loading models list...",
			ModelActivating:     "Loading model %s...",
			ModelLoadError:      "Failed to load model: %v",
			ModelLoadSuccess:    "Model %s loaded",
			ModelFetchError:     "Server unavailable: %v",
			ModelFetchEmpty:     "No Whisper models found",
			ModelFetchOK:        "Models found: %d",
			ModelRefresh:        "Refresh",
			ModelDownloaded:     "(installed)",
			ModelNotDownloaded:  "(%.1f GB)",
			ModelSizeGB:         "(%.1f GB)",
			ModelInstall:        "Install",
			ModelInstalling:     "Installing %s...",
			ModelInstallProgress: "Downloading: %s — %d%%",
			ModelInstallDone:    "Model %s installed!",
			ModelInstallError:   "Install error: %v",
			ModelNotInstalled:   "Model not installed on server. Click \"Install\" to download.",
			ModelReady:          "Model installed and ready",

			// Параметры транскрибации
			LabelPrompt:      "Prompt:",
			LabelTemperature: "Temperature:",
			HintPrompt:       "Hint for model: punctuation style, names, terms",
			HintTemperature:  "0.0 — precise, 0.2-0.4 — diverse (punctuation, language switch)",

			// Автовставка
			CheckboxAutoPaste: "Auto-paste text (Ctrl+V)",
			HintAutoPaste:     "After transcription, text is automatically pasted at cursor position",
			SectionBehavior:   "Behavior",

			// Обрезка тишины и диагностика аудио
			SilenceTrimmed:     "✂️ Silence trimmed: %.1f sec\n",
			RecordingSilent:    "Recording contains only silence",
			AudioStats:         "🎙️ Audio: %.1f sec, %d samples, peak: %.4f, RMS: %.4f\n",
			AudioSilentWarning: "⚠️ Microphone not capturing audio (peak=0). Check microphone permissions in OS settings.",

			// Бэкенд транскрибации
			SectionBackend:        "Transcription Backend",
			LabelBackend:          "Backend:",
			BackendLemonade:       "Lemonade Server (local)",
			BackendWhisperAPI:     "Whisper API (external server)",
			HintBackendLemonade:   "Local AI server with model management",
			HintBackendWhisperAPI: "External Whisper server (Docker, network)",

			// Whisper API
			TabWhisperAPI:         "Whisper API",
			SectionWhisperAPI:     "Whisper API",
			WhisperAPIHintURL:     "Server address, e.g. http://192.168.1.50:9000",
			WhisperAPIStatus:      "Status: checking...",
			WhisperAPIStatusOK:    "Server available",
			WhisperAPIStatusError: "Server unavailable: %v",
			WhisperAPICheckBtn:    "Check",
			WhisperAPINotInit:     "Whisper API client not initialized!",
			BackendInfo:           "Backend: %s",
			ConfigWhisperAPI:      "   🌐 Whisper API: URL=%s, Language=%s",
			ConfigBackend:         "   🔀 Backend: %s",

			// FastFlowLM
			BackendFastFlowLM:       "FastFlowLM (NPU, local)",
			TabFastFlowLM:           "FastFlowLM",
			SectionFastFlowLM:       "FastFlowLM",
			HintBackendFastFlowLM:   "AMD Ryzen AI NPU — local transcription",
			FastFlowLMHintURL:       "Server address, default http://localhost:52625",
			FastFlowLMStatus:        "Status: checking...",
			FastFlowLMStatusOK:      "Server available",
			FastFlowLMStatusError:   "Server unavailable: %v",
			FastFlowLMCheckBtn:      "Check",
			FastFlowLMNotInit:       "FastFlowLM client not initialized!",
			FastFlowLMLLMModel:      "LLM model:",
			FastFlowLMHintLLMModel:  "LLM model for server launch (e.g. llama3.2:1b)",
			FastFlowLMNotInstalled:  "FLM not found in PATH. Install FastFlowLM: https://github.com/amd/FastFlowLM",
			FastFlowLMStarting:      "Starting FLM server...",
			FastFlowLMStarted:       "FLM server started",
			FastFlowLMStartError:    "FLM server start error: %v",
			ConfigFastFlowLM:        "   ⚡ FastFlowLM: URL=%s, Model=%s, Language=%s",

			// Ошибки (внутренние, для логов)
			ErrorAudioDuration:    "Failed to get audio duration: %v",
			ErrorAutostartEnable:  "Autostart enable error: %v",
			ErrorAutostartDisable: "Autostart disable error: %v",
			ErrorReadStdin:        "Failed to read stdin: %v\n",
			ErrorParseJSON:        "Failed to parse JSON: %v\n",
		}

	default:
		// По умолчанию английский
		return Get("en")
	}
}
