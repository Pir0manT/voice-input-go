package main

import (
	"fmt"
	"os"

	"github.com/go-audio/wav"
	"github.com/Pir0manT/voice-input-go/internal/autostart"
	"github.com/Pir0manT/voice-input-go/internal/clipboard"
	"github.com/Pir0manT/voice-input-go/internal/config"
	"github.com/Pir0manT/voice-input-go/internal/console"
	"github.com/Pir0manT/voice-input-go/internal/editor"
	"github.com/Pir0manT/voice-input-go/internal/hotkeys"
	"github.com/Pir0manT/voice-input-go/internal/i18n"
	"github.com/Pir0manT/voice-input-go/internal/lemonade"
	"github.com/Pir0manT/voice-input-go/internal/logger"
	"github.com/Pir0manT/voice-input-go/internal/logview"
	"github.com/Pir0manT/voice-input-go/internal/notify"
	"github.com/Pir0manT/voice-input-go/internal/paste"
	"github.com/Pir0manT/voice-input-go/internal/recorder"
	"github.com/Pir0manT/voice-input-go/internal/settings"
	"github.com/Pir0manT/voice-input-go/internal/singleton"
	"github.com/Pir0manT/voice-input-go/internal/tray"
)

var (
	rec  *recorder.Recorder
	lmn  *lemonade.Client
	cfg  *config.Config
	lang string = "ru"
	msg  *i18n.Messages
	ed   *editor.Editor
)

// getWavDuration возвращает длительность WAV файла в секундах
func getWavDuration(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Используем готовую библиотеку для декодирования WAV
	dec := wav.NewDecoder(file)

	// Получаем длительность в секундах
	duration, err := dec.Duration()
	if err != nil {
		return 0, err
	}
	return duration.Seconds(), nil
}

// transcribeAndCopy отправляет аудио на транскрибацию и копирует текст в буфер
func transcribeAndCopy(audioPath string) {
	msg := i18n.Get(lang)

	defer func() {
		if r := recover(); r != nil {
			logger.Error(msg.PanicTranscription, r)
			fmt.Printf(msg.PanicPrefix, r)
			tray.SetStatus(tray.StatusIdle)
		}
	}()

	// Удаляем временный файл при любом исходе
	defer func() {
		if err := os.Remove(audioPath); err != nil {
			logger.Error(msg.ErrorDeleteFile, err)
		} else {
			logger.Debug(msg.InfoFileDeleted, audioPath)
		}
	}()

	logger.Info(msg.StartTranscription)
	fmt.Println(msg.TranscribingAudio)
	logger.Debug(msg.AudioFile, audioPath)

	// Получаем и логируем длительность (только для лога, без эмодзи)
	var audioDuration float64
	if duration, err := getWavDuration(audioPath); err == nil {
		logger.Debug(msg.AudioDuration, duration)
		audioDuration = duration
	} else {
		logger.Debug(msg.ErrorAudioDuration, err)
	}

	// Проверяем что клиент и конфиг инициализированы
	if lmn == nil {
		logger.Error(msg.LemonadeNotInit)
		fmt.Println(msg.LemonadeNotInit)
		tray.SetStatus(tray.StatusIdle)
		return
	}

	if cfg == nil {
		logger.Error(msg.ConfigNotInit)
		fmt.Println(msg.ConfigNotInit)
		tray.SetStatus(tray.StatusIdle)
		return
	}

	logger.Debug(msg.LemonadeURL, cfg.Lemonade.URL)
	logger.Info(msg.ModelInfo, cfg.Lemonade.Model, cfg.Lemonade.Language)
	fmt.Printf(msg.ModelInfo+"\n", cfg.Lemonade.Model, cfg.Lemonade.Language)

	// Проверяем существование файла
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		logger.Error(msg.AudioFileNotFound, audioPath)
		fmt.Printf(msg.ErrorPrefix, msg.AudioFileNotFound)
		tray.SetStatus(tray.StatusIdle)
		return
	}

	result, err := lmn.TranscribeWithStats(audioPath, cfg.Lemonade.Model, cfg.Lemonade.Language, cfg.Lemonade.Prompt, cfg.Lemonade.Temperature)
	if err != nil {
		logger.Error(msg.TranscriptionError, err)
		fmt.Printf(msg.ErrorPrefix, msg.TranscriptionError)
		tray.SetStatus(tray.StatusIdle)
		return
	}

	// Выводим красивую статистику
	fmt.Println("============================================================")
	fmt.Println(msg.TranscriptionComplete)
	if audioDuration > 0 {
		fmt.Println(fmt.Sprintf(msg.AudioDurationLabel, audioDuration))
	}
	fmt.Println(fmt.Sprintf(msg.ProcessTime, result.ProcessTime))
	fmt.Println(fmt.Sprintf(msg.SpeedInfo, result.Speed))
	fmt.Printf(msg.Characters, len(result.Text))
	fmt.Printf(msg.Backend, result.Backend)
	fmt.Print(msg.TextCopied)
	fmt.Println("============================================================")

	// В лог пишем только кратко
	logger.Info(msg.TranscriptionStats, result.ProcessTime, result.Speed, len(result.Text))

	// Сохраняем текст в редакторе
	ed.SetText(result.Text)
	ed.AddToHistory(result.Text)

	// Копируем в буфер обмена
	if err := clipboard.Copy(result.Text); err != nil {
		logger.Error(msg.CopyError, err)
		fmt.Printf(msg.CopyFailed, err)
	} else if cfg.AutoPaste {
		paste.SimulateCtrlV()
	}

	// Воспроизводим звук уведомления
	if cfg.Notifications.Sound {
		if err := notify.PlaySound(); err != nil {
			logger.Error(msg.ErrorPlaySound, err)
		}
	}

	// Показываем toast-уведомление с превью текста
	if cfg.Notifications.Toast {
		preview := result.Text
		previewRunes := []rune(preview)
		if len(previewRunes) > 100 {
			preview = string(previewRunes[:100]) + "..."
		}
		if err := notify.ShowToast(msg.ToastTitle, preview); err != nil {
			logger.Error(msg.ErrorShowToast, err)
			fmt.Printf(msg.ErrorShowToast+"\n", err)
		}
	}

	// Возвращаем статус
	tray.SetStatus(tray.StatusIdle)
}

func main() {
	// Если запущен с --editor, работаем как GUI subprocess
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--editor":
			editor.RunEditorGUI()
			return
		case "--settings":
			settings.RunSettingsGUI()
			return
		case "--logs":
			logview.RunLogsGUI()
			return
		}
	}

	// Сначала загружаем конфиг чтобы определить язык
	var err error
	cfg, err = config.Load()
	if err != nil {
		cfg = config.Default()
	}

	// Определяем язык из конфига
	lang = cfg.AppLanguage
	if lang == "" {
		lang = "ru"
	}

	// Получаем сообщения для текущего языка
	msg = i18n.Get(lang)

	// Скрываем/показываем консоль в зависимости от настройки
	console.SetVisible(cfg.ShowConsole)

	// Проверяем, запущено ли уже приложение
	appName := singleton.GetMutexName()
	appInstance, err := singleton.New(appName)
	if err != nil {
		fmt.Printf(msg.ErrorPrefix, msg.ErrorAlreadyRunning)
		fmt.Println(msg.ErrorAlreadyRunningWait)
		os.Exit(1)
	}
	defer appInstance.Release()

	// Инициализация логгера: если логирование включено — пишем всё, иначе только ошибки
	logPath := "voice-input-go.log"
	logLevel := logger.ERROR
	if cfg.Logging.Enabled {
		logLevel = logger.DEBUG
	}
	if err := logger.Init(logPath, logLevel, lang); err != nil {
		fmt.Printf(msg.ErrorInitLogger, err)
		os.Exit(1)
	}
	defer logger.Close()

	// Получаем полный путь к лог-файлу для окна логов
	logFilePath, _ := config.GetLogFilePath()

	logger.Info(msg.AppStarting)

	// Логируем конфиг (уже загружен выше)
	if err != nil {
		logger.Error(msg.ErrorConfigLoad, err)
		logger.Info(msg.InfoUsingDefaultConfig)
	}

	// Выводим конфиг в читаемом формате
	logger.Debug(msg.ConfigHeader)
	logger.Debug(msg.ConfigHotkeys,
		cfg.Hotkeys.Start, cfg.Hotkeys.Stop, cfg.Hotkeys.Editor)
	logger.Debug(msg.ConfigLemonade,
		cfg.Lemonade.URL, cfg.Lemonade.Model, cfg.Lemonade.Language)
	logger.Debug(msg.ConfigNotifications,
		cfg.Notifications.Sound, cfg.Notifications.Toast)
	logger.Debug(msg.ConfigAutostart, cfg.Autostart)
	logger.Debug(msg.ConfigLogging,
		cfg.Logging.Enabled, cfg.Logging.Level)
	logger.Debug(msg.ConfigLanguage, cfg.AppLanguage)

	// Создаем рекордер
	rec = recorder.New()
	rec.SetLanguage(lang)

	// Создаем редактор с персистентной историей
	ed = editor.New(lang, cfg.HistorySize)
	ed.LoadHistory()

	// Создаем Lemonade клиент
	lmn = lemonade.NewClient(cfg.Lemonade.URL)

	// Создаем менеджер хоткеев
	hkManager := hotkeys.New()

	// Создаем менеджер настроек и логов
	settingsMgr := settings.New(cfg, lang, logFilePath)

	// Callbacks для меню
	callbacks := map[string]func(){
		"start": func() {
			logger.Debug(msg.RecordingFromMenu)
			if err := rec.Start(); err != nil {
				logger.Error(msg.ErrorStartingRec, err)
				fmt.Printf(msg.ErrorPrefix, err)
			} else {
				tray.SetStatus(tray.StatusRecording)
				fmt.Println(msg.RecordingStarted)
			}
		},
		"stop": func() {
			logger.Debug(msg.StopRecordingFromMenu)
			result, err := rec.Stop()
			if err != nil {
				logger.Error(msg.ErrorStopRecording, err)
				fmt.Printf(msg.ErrorPrefix, err)
			} else {
				tray.SetStatus(tray.StatusProcessing)
				fmt.Printf(msg.RecordingSaved, result.FilePath)
				if result.TrimmedSeconds > 0 {
					m := i18n.Get(lang)
					fmt.Printf(m.SilenceTrimmed, result.TrimmedSeconds)
					logger.Debug(m.SilenceTrimmed, result.TrimmedSeconds)
				}

				// Запускаем транскрибацию в горутине
				go transcribeAndCopy(result.FilePath)
			}
		},
		"editor": func() {
			logger.Debug(msg.EditorFromMenu)
			if err := ed.Show(); err != nil {
				logger.Error(msg.EditorProcessError, err)
			}
		},
		"settings": func() {
			logger.Debug(msg.SettingsFromMenu)
			settingsMgr.Show()
		},
		"logs": func() {
			logger.Debug(msg.LogsFromMenu)
			settingsMgr.ShowLogs()
		},
	}

	// При открытии настроек — снимаем глобальные хоткеи, чтобы не мешали захвату
	settingsMgr.SetOnOpen(func() {
		hkManager.Unregister()
	})

	// При закрытии настроек — регистрируем хоткеи обратно (из актуального конфига)
	settingsMgr.SetOnClose(func() {
		if err := hkManager.Register(
			cfg.Hotkeys.Start, cfg.Hotkeys.Stop, cfg.Hotkeys.Editor,
			callbacks["start"], callbacks["stop"], callbacks["editor"], lang,
		); err != nil {
			logger.Error(msg.ErrorHotkeyRegister, err)
		}
		if err := hkManager.Start(lang); err != nil {
			logger.Error(msg.ErrorHotkeyListener, err)
		}
	})

	// Callback при сохранении настроек
	settingsMgr.SetOnConfigChange(func(newCfg *config.Config) {
		oldLang := cfg.AppLanguage
		oldURL := cfg.Lemonade.URL
		oldModel := cfg.Lemonade.Model

		// Обновляем глобальный конфиг
		cfg = newCfg

		// Если URL изменился — пересоздаём клиент
		if newCfg.Lemonade.URL != oldURL {
			lmn = lemonade.NewClient(newCfg.Lemonade.URL)
		}

		// Если модель изменилась — загружаем новую
		if newCfg.Lemonade.Model != oldModel {
			go func() {
				m := i18n.Get(lang)
				logger.Info(m.ModelActivating, newCfg.Lemonade.Model)
				fmt.Println(fmt.Sprintf(m.ModelActivating, newCfg.Lemonade.Model))
				if err := lmn.LoadModel(newCfg.Lemonade.Model); err != nil {
					logger.Error(m.ModelLoadError, err)
					fmt.Println(fmt.Sprintf(m.ModelLoadError, err))
				} else {
					logger.Info(m.ModelLoadSuccess, newCfg.Lemonade.Model)
					fmt.Println(fmt.Sprintf(m.ModelLoadSuccess, newCfg.Lemonade.Model))
				}
			}()
		}

		// Обновляем автозапуск
		if newCfg.Autostart {
			if err := autostart.Enable(); err != nil {
				logger.Error(msg.ErrorAutostartEnable, err)
			}
		} else {
			if err := autostart.Disable(); err != nil {
				logger.Error(msg.ErrorAutostartDisable, err)
			}
		}

		// Обновляем видимость консоли
		console.SetVisible(newCfg.ShowConsole)

		// Обновляем уровень логирования
		if newCfg.Logging.Enabled {
			logger.SetLevel(logger.DEBUG)
		} else {
			logger.SetLevel(logger.ERROR)
		}

		// Обновляем размер истории
		ed.SetHistorySize(newCfg.HistorySize)

		// Если язык изменился
		if newCfg.AppLanguage != oldLang {
			lang = newCfg.AppLanguage
			msg = i18n.Get(lang)
			rec.SetLanguage(lang)
			ed.SetLanguage(lang)
			hotkeys.SetLanguage(lang)
			settingsMgr.SetLanguage(lang)

			logger.Info(msg.LanguageSwitched, lang)
			tray.Restart(cfg, callbacks, lang)
		}

		// Хоткеи не перерегистрируем здесь — это сделает onClose с актуальным конфигом

		// Обновляем конфиг в менеджере настроек
		settingsMgr.SetConfig(cfg)
	})

	// Регистрируем хоткеи из конфига
	if err := hkManager.Register(cfg.Hotkeys.Start, cfg.Hotkeys.Stop, cfg.Hotkeys.Editor,
		callbacks["start"], callbacks["stop"], callbacks["editor"], lang); err != nil {
		logger.Error(msg.ErrorHotkeyRegister, err)
		fmt.Printf(msg.WarningHotkeys, err)
		fmt.Println(msg.WarningHotkeysDetail)
	}

	// Запускаем прослушивание хоткеев
	if err := hkManager.Start(lang); err != nil {
		logger.Error(msg.ErrorHotkeyListener, err)
	}

	// Запуск системного трея (блокирующий вызов)
	logger.Debug(msg.StartingTray)
	fmt.Println(msg.AppStarted)
	fmt.Printf(msg.LogFile, logPath)
	fmt.Println(msg.TrayMenu)
	fmt.Println(msg.PressCtrlC)
	fmt.Println(msg.TrayStarted)

	// Запускаем трей (это блокирующий вызов)
	tray.Start(cfg, callbacks, lang)
}
