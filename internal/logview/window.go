package logview

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/Pir0manT/voice-input-go/internal/i18n"
)

// LogsInput данные для subprocess
type LogsInput struct {
	LogPath string `json:"logPath"`
	Lang    string `json:"lang"`
}

// RunLogsGUI запускает Fyne GUI просмотра логов (вызывается из subprocess с --logs)
func RunLogsGUI() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, i18n.Get("en").ErrorReadStdin, err)
		os.Exit(1)
	}

	var input LogsInput
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, i18n.Get("en").ErrorParseJSON, err)
		os.Exit(1)
	}

	msg := i18n.Get(input.Lang)

	a := app.New()
	w := a.NewWindow(msg.LogsWindowTitle)
	w.Resize(fyne.NewSize(700, 500))

	// Текстовое поле для логов
	logEntry := widget.NewMultiLineEntry()
	logEntry.Wrapping = fyne.TextWrapWord

	// Счётчик строк
	lineCount := widget.NewLabel("")

	// Функция загрузки логов
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

		// Фильтрация по уровню
		if filterLevel != "" && filterLevel != msg.LogsFilterAll {
			var filtered []string
			scanner := bufio.NewScanner(strings.NewReader(text))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "["+filterLevel+"]") {
					filtered = append(filtered, line)
				}
			}
			if len(filtered) == 0 {
				logEntry.SetText(msg.LogsEmpty)
				lineCount.SetText(fmt.Sprintf(msg.LogsLineCount, 0))
				return
			}
			text = strings.Join(filtered, "\n")
			lineCount.SetText(fmt.Sprintf(msg.LogsLineCount, len(filtered)))
		} else {
			lines := strings.Count(text, "\n")
			if len(text) > 0 && !strings.HasSuffix(text, "\n") {
				lines++
			}
			lineCount.SetText(fmt.Sprintf(msg.LogsLineCount, lines))
		}

		logEntry.SetText(text)
		// Скролл вниз
		logEntry.CursorRow = strings.Count(text, "\n")
	}

	// Фильтр по уровню
	currentFilter := msg.LogsFilterAll
	filterOptions := []string{msg.LogsFilterAll, "DEBUG", "INFO", "WARN", "ERROR"}
	filterSelect := widget.NewSelect(filterOptions, func(selected string) {
		currentFilter = selected
		loadLogs(selected)
	})
	filterSelect.SetSelected(msg.LogsFilterAll)

	// Кнопка обновить
	refreshBtn := widget.NewButton(msg.LogsRefresh, func() {
		loadLogs(currentFilter)
	})

	// Кнопка очистить
	clearBtn := widget.NewButton(msg.LogsClear, func() {
		dialog.ShowConfirm(msg.LogsClearConfirm, msg.LogsClearConfirm, func(ok bool) {
			if ok {
				os.WriteFile(input.LogPath, []byte{}, 0644)
				loadLogs(currentFilter)
			}
		}, w)
	})

	// Верхняя панель
	topBar := container.NewHBox(
		widget.NewLabel(msg.LabelLoggingLevel),
		filterSelect,
		layout.NewSpacer(),
		refreshBtn,
		clearBtn,
	)

	// Нижняя панель
	closeBtn := widget.NewButton(msg.Close, func() {
		a.Quit()
	})

	bottomBar := container.NewHBox(
		lineCount,
		layout.NewSpacer(),
		closeBtn,
	)

	content := container.NewBorder(
		topBar,    // top
		bottomBar, // bottom
		nil,       // left
		nil,       // right
		logEntry,  // center
	)

	w.SetContent(content)

	// Загружаем логи
	loadLogs(msg.LogsFilterAll)

	w.ShowAndRun()
}
