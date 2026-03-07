package editor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/Pir0manT/voice-input-go/internal/i18n"
)

// EditorInput данные для subprocess
type EditorInput struct {
	Text         string   `json:"text"`
	History      []string `json:"history"`
	HistoryIndex int      `json:"historyIndex"`
	Lang         string   `json:"lang"`
}

// EditorOutput результат от subprocess
type EditorOutput struct {
	Text           string `json:"text"`
	HistoryIndex   int    `json:"historyIndex"`
	DeletedIndices []int  `json:"deletedIndices,omitempty"`
	Cancelled      bool   `json:"cancelled,omitempty"`
}

// RunEditorGUI запускает Fyne GUI редактора (вызывается из subprocess с --editor)
func RunEditorGUI() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read stdin: %v\n", err)
		os.Exit(1)
	}

	var input EditorInput
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse input JSON: %v\n", err)
		os.Exit(1)
	}

	msg := i18n.Get(input.Lang)
	historyIndex := input.HistoryIndex
	history := append([]string{}, input.History...) // рабочая копия
	var deletedIndices []int                         // оригинальные индексы удалённых

	// Маппинг: текущий индекс в history → оригинальный индекс в input.History
	origIndices := make([]int, len(history))
	for i := range origIndices {
		origIndices[i] = i
	}

	a := app.New()
	w := a.NewWindow(msg.EditorWindowTitle)
	w.Resize(fyne.NewSize(600, 400))

	// Текстовое поле
	entry := widget.NewMultiLineEntry()
	entry.SetText(input.Text)
	entry.Wrapping = fyne.TextWrapWord

	// Счётчик символов
	charCount := widget.NewLabel(fmt.Sprintf(msg.EditorCharCount, len(input.Text)))
	entry.OnChanged = func(s string) {
		charCount.SetText(fmt.Sprintf(msg.EditorCharCount, len(s)))
	}

	// Навигация по истории
	historyLabel := widget.NewLabel("")
	prevBtn := widget.NewButton(msg.EditorPrev, nil)
	nextBtn := widget.NewButton(msg.EditorNext, nil)
	deleteBtn := widget.NewButton(msg.EditorDelete, nil)
	deleteBtn.Importance = widget.DangerImportance

	updateHistoryUI := func() {
		if len(history) == 0 {
			historyLabel.SetText(msg.EditorNoText)
			prevBtn.Disable()
			nextBtn.Disable()
			deleteBtn.Disable()
			return
		}
		historyLabel.SetText(fmt.Sprintf(msg.EditorHistoryPos, historyIndex+1, len(history)))
		deleteBtn.Enable()
		prevBtn.Enable()
		nextBtn.Enable()
		if historyIndex <= 0 {
			prevBtn.Disable()
		}
		if historyIndex >= len(history)-1 {
			nextBtn.Disable()
		}
	}

	prevBtn.OnTapped = func() {
		if historyIndex > 0 {
			historyIndex--
			entry.SetText(history[historyIndex])
			updateHistoryUI()
		}
	}

	nextBtn.OnTapped = func() {
		if historyIndex < len(history)-1 {
			historyIndex++
			entry.SetText(history[historyIndex])
			updateHistoryUI()
		}
	}

	deleteBtn.OnTapped = func() {
		if len(history) == 0 {
			return
		}

		// Запоминаем оригинальный индекс удалённой записи
		deletedIndices = append(deletedIndices, origIndices[historyIndex])

		// Удаляем из рабочих массивов
		history = append(history[:historyIndex], history[historyIndex+1:]...)
		origIndices = append(origIndices[:historyIndex], origIndices[historyIndex+1:]...)

		if len(history) == 0 {
			historyIndex = -1
			entry.SetText("")
		} else {
			if historyIndex >= len(history) {
				historyIndex = len(history) - 1
			}
			entry.SetText(history[historyIndex])
		}
		updateHistoryUI()
	}

	updateHistoryUI()

	// Кнопки действий
	cancelled := true

	copyBtn := widget.NewButton(msg.EditorCopyAndClose, func() {
		output := EditorOutput{
			Text:           entry.Text,
			HistoryIndex:   historyIndex,
			DeletedIndices: deletedIndices,
		}
		data, _ := json.Marshal(output)
		fmt.Fprint(os.Stdout, string(data))
		cancelled = false
		a.Quit()
	})
	copyBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton(msg.Cancel, func() {
		cancelled = true
		a.Quit()
	})

	// Панель навигации по истории
	historyBar := container.NewHBox(
		prevBtn,
		deleteBtn,
		layout.NewSpacer(),
		historyLabel,
		layout.NewSpacer(),
		nextBtn,
	)

	// Нижняя панель
	bottomBar := container.NewHBox(
		charCount,
		layout.NewSpacer(),
		cancelBtn,
		copyBtn,
	)

	// Компоновка
	content := container.NewBorder(
		historyBar, // top
		bottomBar,  // bottom
		nil,        // left
		nil,        // right
		entry,      // center
	)

	w.SetContent(content)

	w.SetOnClosed(func() {
		if cancelled {
			// Если были удаления — отправляем их даже при отмене
			if len(deletedIndices) > 0 {
				output := EditorOutput{
					Text:           entry.Text,
					HistoryIndex:   historyIndex,
					DeletedIndices: deletedIndices,
					Cancelled:      true,
				}
				data, _ := json.Marshal(output)
				fmt.Fprint(os.Stdout, string(data))
				os.Exit(0)
			}
			os.Exit(1)
		}
	})

	w.ShowAndRun()
}

// launchEditorWindow запускает GUI процесс редактора (сам себя с --editor)
func launchEditorWindow(input EditorInput) (*EditorOutput, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	inputData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	cmd := exec.Command(exePath, "--editor")
	cmd.Stdin = bytes.NewReader(inputData)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start editor: %w", err)
	}

	err = cmd.Wait()

	// Exit code 1 = отмена
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return nil, nil // отмена, не ошибка
			}
		}
		return nil, fmt.Errorf("editor process error: %w", err)
	}

	var output EditorOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, fmt.Errorf("failed to parse editor output: %w", err)
	}

	return &output, nil
}
