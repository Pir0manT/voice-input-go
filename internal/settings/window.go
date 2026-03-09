package settings

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/Pir0manT/voice-input-go/internal/config"
	"github.com/Pir0manT/voice-input-go/internal/i18n"
	"github.com/Pir0manT/voice-input-go/internal/lemonade"
)

// SettingsInput данные для subprocess
type SettingsInput struct {
	Config     *config.Config `json:"config"`
	Lang       string         `json:"lang"`
	LogPath    string         `json:"logPath"`
	InitialTab int            `json:"initialTab"`
}

// SettingsOutput результат от subprocess
type SettingsOutput struct {
	Config *config.Config `json:"config"`
}

// RunSettingsGUI запускает Fyne GUI настроек (вызывается из subprocess с --settings)
func RunSettingsGUI() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, i18n.Get("en").ErrorReadStdin, err)
		os.Exit(1)
	}

	var input SettingsInput
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, i18n.Get("en").ErrorParseJSON, err)
		os.Exit(1)
	}

	cfg := input.Config
	msg := i18n.Get(input.Lang)

	a := app.New()
	w := a.NewWindow(msg.SettingsTitle)
	w.Resize(fyne.NewSize(600, 580))

	// === Выбор бэкенда ===
	backendOptions := []string{msg.BackendLemonade, msg.BackendWhisperAPI}
	backendValues := []string{"lemonade", "whisper-api"}
	currentBackendDisplay := msg.BackendLemonade
	for i, v := range backendValues {
		if v == cfg.Backend {
			currentBackendDisplay = backendOptions[i]
			break
		}
	}
	backendSelect := widget.NewSelect(backendOptions, nil)
	backendSelect.SetSelected(currentBackendDisplay)

	// === Вкладка "Основные" ===

	// Горячие клавиши
	hotkeyStartVal := ConfigToDisplay(cfg.Hotkeys.Start)
	hotkeyStopVal := ConfigToDisplay(cfg.Hotkeys.Stop)
	hotkeyEditorVal := ConfigToDisplay(cfg.Hotkeys.Editor)

	hotkeyStart := NewHotkeyEntry(hotkeyStartVal, func(v string) { hotkeyStartVal = v })
	hotkeyStop := NewHotkeyEntry(hotkeyStopVal, func(v string) { hotkeyStopVal = v })
	hotkeyEditor := NewHotkeyEntry(hotkeyEditorVal, func(v string) { hotkeyEditorVal = v })

	hotkeyCol1 := container.NewVBox(widget.NewLabel(msg.LabelHotkeyStart), hotkeyStart)
	hotkeyCol2 := container.NewVBox(widget.NewLabel(msg.LabelHotkeyStop), hotkeyStop)
	hotkeyCol3 := container.NewVBox(widget.NewLabel(msg.LabelHotkeyEditor), hotkeyEditor)

	hotkeysSection := container.NewVBox(
		widget.NewLabelWithStyle(msg.SectionHotkeys, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(3, hotkeyCol1, hotkeyCol2, hotkeyCol3),
		widget.NewSeparator(),
	)

	// Lemonade Server
	urlEntry := widget.NewEntry()
	urlEntry.SetText(cfg.Lemonade.URL)

	// --- Модель: динамический список с сервера ---
	var modelList []lemonade.ModelInfo // полный список Whisper моделей

	modelOptions := []string{cfg.Lemonade.Model}
	modelSelect := widget.NewSelect(modelOptions, nil)
	modelSelect.SetSelected(cfg.Lemonade.Model)

	modelStatus := widget.NewLabel(msg.ModelLoading)
	installBtn := widget.NewButton(msg.ModelInstall, nil)
	installBtn.Hide()

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// Вспомогательная функция: извлечь чистый ID из строки "ModelID (статус)"
	extractModelID := func(display string) string {
		if idx := strings.Index(display, " ("); idx > 0 {
			return display[:idx]
		}
		return display
	}

	// Проверяет, установлена ли модель по ID
	isModelDownloaded := func(modelID string) bool {
		for _, m := range modelList {
			if m.ID == modelID {
				return m.Downloaded
			}
		}
		// Если модели нет в списке (сервер недоступен) — считаем установленной
		// чтобы не блокировать сохранение с текущей моделью
		return true
	}

	// Обновляет UI в зависимости от выбранной модели
	var saveBtn *widget.Button // forward declaration, будет присвоен ниже

	updateModelUI := func(modelID string) {
		downloaded := isModelDownloaded(modelID)
		if downloaded {
			modelStatus.SetText(msg.ModelReady)
			installBtn.Hide()
			progressBar.Hide()
			if saveBtn != nil {
				saveBtn.Enable()
			}
		} else {
			modelStatus.SetText(msg.ModelNotInstalled)
			installBtn.Show()
			progressBar.Hide()
			if saveBtn != nil {
				saveBtn.Disable()
			}
		}
	}

	// При выборе модели в Select
	modelSelect.OnChanged = func(selected string) {
		modelID := extractModelID(selected)
		updateModelUI(modelID)
	}

	// Функция загрузки списка моделей с сервера
	fetchModels := func(serverURL string) {
		modelStatus.SetText(msg.ModelLoading)
		installBtn.Hide()
		progressBar.Hide()
		go func() {
			client := lemonade.NewClient(serverURL)
			models, err := client.GetWhisperModels()
			if err != nil {
				errText := fmt.Sprintf(msg.ModelFetchError, err)
				fyne.Do(func() { modelStatus.SetText(errText) })
				return
			}
			if len(models) == 0 {
				fyne.Do(func() { modelStatus.SetText(msg.ModelFetchEmpty) })
				return
			}

			modelList = models

			// Формируем список: "ModelID (статус)"
			options := make([]string, 0, len(models))
			for _, m := range models {
				label := m.ID
				if m.Downloaded {
					label += " " + msg.ModelDownloaded
				} else if m.Size > 0 {
					label += " " + fmt.Sprintf(msg.ModelNotDownloaded, m.Size)
				}
				options = append(options, label)
			}

			// Находим текущую модель в списке
			currentSelected := ""
			for _, opt := range options {
				if strings.HasPrefix(opt, cfg.Lemonade.Model) {
					currentSelected = opt
					break
				}
			}

			fyne.Do(func() {
				modelSelect.Options = options
				modelSelect.Refresh()
				if currentSelected != "" {
					modelSelect.SetSelected(currentSelected)
				}
				selectedID := extractModelID(modelSelect.Selected)
				updateModelUI(selectedID)
			})
		}()
	}

	// Кнопка установки модели
	installBtn.OnTapped = func() {
		modelID := extractModelID(modelSelect.Selected)
		installBtn.Disable()
		progressBar.Hide()
		modelStatus.SetText(fmt.Sprintf(msg.ModelInstalling, modelID))

		go func() {
			progressShown := false
			client := lemonade.NewClient(urlEntry.Text)
			err := client.PullModel(modelID, func(p lemonade.PullProgress) {
				val := float64(p.Percent) / 100.0
				text := fmt.Sprintf(msg.ModelInstallProgress, p.File, p.Percent)
				show := !progressShown
				progressShown = true
				fyne.Do(func() {
					if show {
						progressBar.SetValue(0)
						progressBar.Show()
					}
					progressBar.SetValue(val)
					modelStatus.SetText(text)
				})
			})

			if err != nil {
				errText := fmt.Sprintf(msg.ModelInstallError, err)
				fyne.Do(func() {
					modelStatus.SetText(errText)
					installBtn.Enable()
					progressBar.Hide()
				})
				return
			}

			// Обновляем статус в локальном списке
			for i := range modelList {
				if modelList[i].ID == modelID {
					modelList[i].Downloaded = true
					break
				}
			}

			// Перестраиваем опции Select
			options := make([]string, 0, len(modelList))
			for _, m := range modelList {
				label := m.ID
				if m.Downloaded {
					label += " " + msg.ModelDownloaded
				} else if m.Size > 0 {
					label += " " + fmt.Sprintf(msg.ModelNotDownloaded, m.Size)
				}
				options = append(options, label)
			}

			doneText := fmt.Sprintf(msg.ModelInstallDone, modelID)
			fyne.Do(func() {
				progressBar.SetValue(1.0)
				modelStatus.SetText(doneText)
				progressBar.Hide()
				installBtn.Hide()
				if saveBtn != nil {
					saveBtn.Enable()
				}

				modelSelect.Options = options
				modelSelect.Refresh()
				for _, opt := range options {
					if strings.HasPrefix(opt, modelID) {
						modelSelect.SetSelected(opt)
						break
					}
				}
			})
		}()
	}

	refreshModelsBtn := widget.NewButton(msg.ModelRefresh, func() {
		fetchModels(urlEntry.Text)
	})

	// Запускаем первую загрузку
	fetchModels(cfg.Lemonade.URL)

	langDisplayNames := []string{msg.LangRussian, msg.LangEnglish}
	langValues := []string{"ru", "en"}
	currentLangDisplay := msg.LangRussian
	for i, v := range langValues {
		if v == cfg.Lemonade.Language {
			currentLangDisplay = langDisplayNames[i]
			break
		}
	}
	lemonLangSelect := widget.NewSelect(langDisplayNames, nil)
	lemonLangSelect.SetSelected(currentLangDisplay)

	modelRow := container.NewBorder(nil, nil, nil, refreshModelsBtn, modelSelect)
	modelActionRow := container.NewBorder(nil, nil, installBtn, nil, progressBar)

	// Температура
	tempEntry := widget.NewEntry()
	tempEntry.SetText(fmt.Sprintf("%.1f", cfg.Lemonade.Temperature))

	// Промпт (подсказка для модели)
	promptEntry := widget.NewMultiLineEntry()
	promptEntry.SetText(cfg.Lemonade.Prompt)
	promptEntry.SetPlaceHolder(msg.HintPrompt)
	promptEntry.SetMinRowsVisible(2)

	lemonadeSection := container.NewVBox(
		widget.NewLabelWithStyle(msg.SectionLemonade, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.New(layout.NewFormLayout(),
			widget.NewLabel(msg.LabelURL), urlEntry,
			widget.NewLabel(msg.LabelModel), modelRow,
		),
		modelStatus,
		modelActionRow,
		container.New(layout.NewFormLayout(),
			widget.NewLabel(msg.LabelLanguage), lemonLangSelect,
			widget.NewLabel(msg.LabelTemperature), tempEntry,
		),
		widget.NewLabel(msg.LabelPrompt),
		promptEntry,
		widget.NewSeparator(),
	)

	// Автозапуск
	autostartCheck := widget.NewCheck(msg.CheckboxAutostart, nil)
	autostartCheck.SetChecked(cfg.Autostart)

	// Автовставка
	autoPasteCheck := widget.NewCheck(msg.CheckboxAutoPaste, nil)
	autoPasteCheck.SetChecked(cfg.AutoPaste)

	// Консоль
	showConsoleCheck := widget.NewCheck(msg.CheckboxShowConsole, nil)
	showConsoleCheck.SetChecked(cfg.ShowConsole)

	autostartSection := container.NewVBox(
		widget.NewLabelWithStyle(msg.SectionBehavior, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(3, autostartCheck, autoPasteCheck, showConsoleCheck),
		widget.NewSeparator(),
	)

	// Язык интерфейса + История — в одну строку
	appLangDisplayNames := []string{msg.LangRussian, msg.LangEnglish}
	appLangValues := []string{"ru", "en"}
	currentAppLangDisplay := msg.LangRussian
	for i, v := range appLangValues {
		if v == cfg.AppLanguage {
			currentAppLangDisplay = appLangDisplayNames[i]
			break
		}
	}
	appLangSelect := widget.NewSelect(appLangDisplayNames, nil)
	appLangSelect.SetSelected(currentAppLangDisplay)

	historySizeOptions := []string{"5", "10", "20", "50", "100"}
	historySizeSelect := widget.NewSelect(historySizeOptions, nil)
	historySizeSelect.SetSelected(fmt.Sprintf("%d", cfg.HistorySize))

	miscSection := container.NewVBox(
		widget.NewLabelWithStyle(msg.SectionAppLanguage, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(2,
			container.New(layout.NewFormLayout(), widget.NewLabel(msg.LabelAppLanguage), appLangSelect),
			container.New(layout.NewFormLayout(), widget.NewLabel(msg.LabelHistorySize), historySizeSelect),
		),
	)

	backendSection := container.NewVBox(
		widget.NewLabelWithStyle(msg.SectionBackend, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.New(layout.NewFormLayout(),
			widget.NewLabel(msg.LabelBackend), backendSelect,
		),
		widget.NewSeparator(),
	)

	generalTab := container.NewVBox(
		hotkeysSection,
		backendSection,
		autostartSection,
		miscSection,
	)

	// === Вкладка "Уведомления" ===

	soundCheck := widget.NewCheck(msg.CheckboxSound, nil)
	soundCheck.SetChecked(cfg.Notifications.Sound)

	toastCheck := widget.NewCheck(msg.CheckboxToast, nil)
	toastCheck.SetChecked(cfg.Notifications.Toast)

	notifSection := container.NewVBox(
		widget.NewLabelWithStyle(msg.SectionNotifications, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		soundCheck,
		toastCheck,
		widget.NewSeparator(),
	)

	// Логирование
	loggingCheck := widget.NewCheck(msg.CheckboxLogging, nil)
	loggingCheck.SetChecked(cfg.Logging.Enabled)

	loggingSection := container.NewVBox(
		widget.NewLabelWithStyle(msg.SectionLogs, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		loggingCheck,
	)

	notifTab := container.NewVBox(
		notifSection,
		loggingSection,
	)

	// === Вкладка "Логи" ===

	const maxLogLines = 500

	logEntry := widget.NewMultiLineEntry()
	logEntry.Wrapping = fyne.TextWrapOff

	lineCount := widget.NewLabel("")
	logsLoaded := false

	// Функция загрузки логов (показывает последние maxLogLines строк)
	loadLogs := func(filterLevel string) {
		content, err := os.ReadFile(input.LogPath)
		if err != nil {
			logEntry.SetText(fmt.Sprintf(msg.LogsFileError, err))
			lineCount.SetText(fmt.Sprintf(msg.LogsLineCount, 0))
			return
		}

		text := string(content)
		if text == "" {
			logEntry.SetText(msg.LogsEmpty)
			lineCount.SetText(fmt.Sprintf(msg.LogsLineCount, 0))
			return
		}

		// Разбиваем на строки
		var allLines []string
		scanner := bufio.NewScanner(strings.NewReader(text))
		for scanner.Scan() {
			allLines = append(allLines, scanner.Text())
		}

		// Фильтрация по уровню
		var lines []string
		if filterLevel != "" && filterLevel != msg.LogsFilterAll {
			for _, line := range allLines {
				if strings.Contains(line, "["+filterLevel+"]") {
					lines = append(lines, line)
				}
			}
		} else {
			lines = allLines
		}

		totalLines := len(lines)
		if totalLines == 0 {
			logEntry.SetText(msg.LogsEmpty)
			lineCount.SetText(fmt.Sprintf(msg.LogsLineCount, 0))
			return
		}

		// Берём только последние maxLogLines строк
		if len(lines) > maxLogLines {
			lines = lines[len(lines)-maxLogLines:]
		}

		result := strings.Join(lines, "\n")
		logEntry.SetText(result)
		logEntry.CursorRow = len(lines) - 1
		lineCount.SetText(fmt.Sprintf(msg.LogsLineCount, totalLines))
	}

	currentFilter := msg.LogsFilterAll
	filterOptions := []string{msg.LogsFilterAll, "DEBUG", "INFO", "WARN", "ERROR"}
	filterSelect := widget.NewSelect(filterOptions, func(selected string) {
		currentFilter = selected
		loadLogs(selected)
	})
	filterSelect.SetSelected(msg.LogsFilterAll)

	refreshBtn := widget.NewButton(msg.LogsRefresh, func() {
		loadLogs(currentFilter)
	})

	clearBtn := widget.NewButton(msg.LogsClear, func() {
		dialog.ShowConfirm(msg.LogsClearConfirm, msg.LogsClearConfirm, func(ok bool) {
			if ok {
				os.WriteFile(input.LogPath, []byte{}, 0644)
				loadLogs(currentFilter)
			}
		}, w)
	})

	logTopBar := container.NewHBox(
		widget.NewLabel(msg.LabelLoggingLevel),
		filterSelect,
		layout.NewSpacer(),
		refreshBtn,
		clearBtn,
	)

	logsTab := container.NewBorder(
		logTopBar,  // top
		lineCount,  // bottom
		nil,        // left
		nil,        // right
		logEntry,   // center
	)

	// === Вкладка "Whisper API" ===
	whisperURLEntry := widget.NewEntry()
	whisperURLEntry.SetText(cfg.WhisperAPI.URL)
	whisperURLEntry.SetPlaceHolder(msg.WhisperAPIHintURL)

	whisperLangDisplayNames := []string{msg.LangRussian, msg.LangEnglish}
	whisperLangValues := []string{"ru", "en"}
	currentWhisperLangDisplay := msg.LangRussian
	for i, v := range whisperLangValues {
		if v == cfg.WhisperAPI.Language {
			currentWhisperLangDisplay = whisperLangDisplayNames[i]
			break
		}
	}
	whisperLangSelect := widget.NewSelect(whisperLangDisplayNames, nil)
	whisperLangSelect.SetSelected(currentWhisperLangDisplay)

	whisperPromptEntry := widget.NewMultiLineEntry()
	whisperPromptEntry.SetText(cfg.WhisperAPI.Prompt)
	whisperPromptEntry.SetPlaceHolder(msg.HintPrompt)
	whisperPromptEntry.SetMinRowsVisible(2)

	whisperStatus := widget.NewLabel("")
	whisperCheckBtn := widget.NewButton(msg.WhisperAPICheckBtn, func() {
		whisperStatus.SetText(msg.WhisperAPIStatus)
		go func() {
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Get(whisperURLEntry.Text + "/")
			if err != nil {
				errText := fmt.Sprintf(msg.WhisperAPIStatusError, err)
				fyne.Do(func() { whisperStatus.SetText(errText) })
				return
			}
			resp.Body.Close()
			fyne.Do(func() { whisperStatus.SetText(msg.WhisperAPIStatusOK) })
		}()
	})

	whisperAPITab := container.NewVBox(
		widget.NewLabelWithStyle(msg.SectionWhisperAPI, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.New(layout.NewFormLayout(),
			widget.NewLabel(msg.LabelURL), whisperURLEntry,
			widget.NewLabel(msg.LabelLanguage), whisperLangSelect,
		),
		widget.NewLabel(msg.LabelPrompt),
		whisperPromptEntry,
		widget.NewSeparator(),
		container.NewHBox(whisperCheckBtn, whisperStatus),
	)

	// === Табы ===
	tabs := container.NewAppTabs(
		container.NewTabItem(msg.TabGeneral, generalTab),
		container.NewTabItem(msg.SectionLemonade, container.NewVScroll(lemonadeSection)),
		container.NewTabItem(msg.TabWhisperAPI, whisperAPITab),
		container.NewTabItem(msg.TabNotifications, notifTab),
		container.NewTabItem(msg.TabLogs, logsTab),
	)

	// Ленивая загрузка логов — только при переходе на вкладку
	tabs.OnSelected = func(tab *container.TabItem) {
		if tab.Text == msg.TabLogs && !logsLoaded {
			logsLoaded = true
			loadLogs(currentFilter)
		}
	}

	// Если запрошена конкретная вкладка (например, логи из трея)
	if input.InitialTab > 0 && input.InitialTab < len(tabs.Items) {
		tabs.SelectIndex(input.InitialTab)
	}

	// === Кнопки действий ===
	cancelled := true

	saveBtn = widget.NewButton(msg.Save, func() {
		// Извлекаем чистый ID модели (до первого пробела со скобкой)
		selectedModel := modelSelect.Selected
		if idx := strings.Index(selectedModel, " ("); idx > 0 {
			selectedModel = selectedModel[:idx]
		}

		newCfg := &config.Config{
			Hotkeys: config.HotkeysConfig{
				Start:  ComboToConfig(hotkeyStart.Text),
				Stop:   ComboToConfig(hotkeyStop.Text),
				Editor: ComboToConfig(hotkeyEditor.Text),
			},
			Backend: langDisplayToValue(backendSelect.Selected, backendOptions, backendValues),
			Lemonade: config.LemonadeConfig{
				URL:         urlEntry.Text,
				Model:       selectedModel,
				Language:    langDisplayToValue(lemonLangSelect.Selected, langDisplayNames, langValues),
				Prompt:      promptEntry.Text,
				Temperature: parseTemperature(tempEntry.Text),
			},
			WhisperAPI: config.WhisperAPIConfig{
				URL:      whisperURLEntry.Text,
				Language: langDisplayToValue(whisperLangSelect.Selected, whisperLangDisplayNames, whisperLangValues),
				Prompt:   whisperPromptEntry.Text,
			},
			Notifications: config.NotificationsConfig{
				Sound: soundCheck.Checked,
				Toast: toastCheck.Checked,
			},
			Autostart: autostartCheck.Checked,
			AutoPaste: autoPasteCheck.Checked,
			Logging: config.LoggingConfig{
				Enabled: loggingCheck.Checked,
				Level:   "info",
			},
			AppLanguage: langDisplayToValue(appLangSelect.Selected, appLangDisplayNames, appLangValues),
			HistorySize: historySizeToInt(historySizeSelect.Selected),
			ShowConsole: showConsoleCheck.Checked,
		}

		output := SettingsOutput{Config: newCfg}
		data, _ := json.Marshal(output)
		fmt.Fprint(os.Stdout, string(data))
		cancelled = false
		a.Quit()
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton(msg.Cancel, func() {
		a.Quit()
	})

	bottomBar := container.NewHBox(
		layout.NewSpacer(),
		cancelBtn,
		saveBtn,
	)

	content := container.NewBorder(
		nil,       // top
		bottomBar, // bottom
		nil,       // left
		nil,       // right
		tabs,      // center
	)

	w.SetContent(content)

	w.SetOnClosed(func() {
		if cancelled {
			os.Exit(1)
		}
	})

	w.ShowAndRun()
}

// langDisplayToValue преобразует отображаемое имя языка в значение
func langDisplayToValue(display string, displayNames, values []string) string {
	for i, name := range displayNames {
		if name == display {
			return values[i]
		}
	}
	return values[0]
}

// parseTemperature парсит строку температуры в float64 (0.0 - 1.0)
func parseTemperature(s string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil || v < 0 {
		return 0.0
	}
	if v > 1.0 {
		return 1.0
	}
	return v
}

// historySizeToInt преобразует строку размера истории в int
func historySizeToInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	if n <= 0 {
		return 20
	}
	return n
}

// launchSettingsWindow запускает GUI процесс настроек (сам себя с --settings)
func launchSettingsWindow(input SettingsInput) (*SettingsOutput, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	inputData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	cmd := exec.Command(exePath, "--settings")
	cmd.Stdin = bytes.NewReader(inputData)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start settings: %w", err)
	}

	err = cmd.Wait()

	// Exit code 1 = отмена
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("settings process error: %w", err)
	}

	var output SettingsOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, fmt.Errorf("failed to parse settings output: %w", err)
	}

	return &output, nil
}
