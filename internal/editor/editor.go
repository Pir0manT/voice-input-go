package editor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Pir0manT/voice-input-go/internal/clipboard"
	"github.com/Pir0manT/voice-input-go/internal/i18n"
)

// historyData формат файла истории на диске
type historyData struct {
	Entries []string `json:"entries"`
	Index   int      `json:"index"`
}

// Editor редактор текста
type Editor struct {
	mu          sync.Mutex
	history     []string
	current     string
	index       int
	lang        string
	running     bool
	historySize int
}

// New создаёт новый редактор
func New(lang string, historySize int) *Editor {
	if historySize <= 0 {
		historySize = 20
	}
	return &Editor{
		history:     make([]string, 0, historySize),
		index:       -1,
		lang:        lang,
		historySize: historySize,
	}
}

// getHistoryPath возвращает путь к файлу истории
func getHistoryPath() (string, error) {
	appData, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appData, "voice-input-go", "history.json"), nil
}

// LoadHistory загружает историю из файла
func (e *Editor) LoadHistory() {
	e.mu.Lock()
	defer e.mu.Unlock()

	path, err := getHistoryPath()
	if err != nil {
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var h historyData
	if err := json.Unmarshal(data, &h); err != nil {
		return
	}

	e.history = h.Entries
	e.index = h.Index

	// Ограничиваем по текущему historySize
	if len(e.history) > e.historySize {
		e.history = e.history[len(e.history)-e.historySize:]
		e.index = len(e.history) - 1
	}

	if e.index >= len(e.history) {
		e.index = len(e.history) - 1
	}
	if e.index >= 0 {
		e.current = e.history[e.index]
	}
}

// saveHistory сохраняет историю в файл (вызывать под мьютексом)
func (e *Editor) saveHistory() {
	path, err := getHistoryPath()
	if err != nil {
		return
	}

	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	h := historyData{
		Entries: e.history,
		Index:   e.index,
	}

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(path, data, 0644)
}

// SetHistorySize обновляет максимальный размер истории
func (e *Editor) SetHistorySize(size int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if size <= 0 {
		size = 20
	}
	e.historySize = size

	// Обрезаем если нужно
	if len(e.history) > size {
		e.history = e.history[len(e.history)-size:]
		if e.index >= len(e.history) {
			e.index = len(e.history) - 1
		}
		e.saveHistory()
	}
}

// SetLanguage устанавливает язык для сообщений
func (e *Editor) SetLanguage(lang string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.lang = lang
}

// SetText устанавливает текст
func (e *Editor) SetText(text string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.current = text
	msg := i18n.Get(e.lang)
	fmt.Printf(msg.EditorTextSet+"\n", len(text))
}

// GetText возвращает текст
func (e *Editor) GetText() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.current
}

// AddToHistory добавляет текст в историю
func (e *Editor) AddToHistory(text string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Добавляем в историю
	e.history = append(e.history, text)
	e.index = len(e.history) - 1

	// Ограничиваем историю
	if len(e.history) > e.historySize {
		e.history = e.history[len(e.history)-e.historySize:]
		e.index = len(e.history) - 1
	}

	e.saveHistory()

	msg := i18n.Get(e.lang)
	fmt.Printf(msg.EditorAddedToHistory+"\n", len(e.history))
}

// GetHistory возвращает историю
func (e *Editor) GetHistory() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.history
}

// GetHistoryIndex возвращает текущий индекс в истории
func (e *Editor) GetHistoryIndex() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.index
}

// CopyToClipboard копирует текст в буфер
func (e *Editor) CopyToClipboard() error {
	e.mu.Lock()
	text := e.current
	lang := e.lang
	e.mu.Unlock()

	if err := clipboard.Copy(text); err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	msg := i18n.Get(lang)
	fmt.Printf(msg.EditorCopiedToClip+"\n", len(text))
	return nil
}

// IsRunning проверяет запущен ли редактор
func (e *Editor) IsRunning() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.running
}

// Show открывает окно редактора
func (e *Editor) Show() error {
	e.mu.Lock()
	msg := i18n.Get(e.lang)

	if e.running {
		e.mu.Unlock()
		fmt.Println(msg.EditorAlreadyOpen)
		return nil
	}

	e.running = true
	input := EditorInput{
		Text:         e.current,
		History:      append([]string{}, e.history...),
		HistoryIndex: e.index,
		Lang:         e.lang,
	}
	e.mu.Unlock()

	fmt.Println(msg.EditorShowWindow)

	go func() {
		defer func() {
			e.mu.Lock()
			e.running = false
			e.mu.Unlock()
		}()

		output, err := launchEditorWindow(input)
		if err != nil {
			fmt.Printf(msg.EditorProcessError+"\n", err)
			return
		}

		if output == nil {
			return
		}

		// Обновляем состояние из результата
		e.mu.Lock()
		// Удаляем записи из истории (с конца, чтобы индексы не сдвигались)
		if len(output.DeletedIndices) > 0 {
			// Сортируем по убыванию
			sorted := make([]int, len(output.DeletedIndices))
			copy(sorted, output.DeletedIndices)
			for i := 0; i < len(sorted); i++ {
				for j := i + 1; j < len(sorted); j++ {
					if sorted[j] > sorted[i] {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				}
			}
			for _, idx := range sorted {
				if idx >= 0 && idx < len(e.history) {
					e.history = append(e.history[:idx], e.history[idx+1:]...)
				}
			}
		}
		e.index = output.HistoryIndex
		if e.index >= len(e.history) {
			e.index = len(e.history) - 1
		}
		if !output.Cancelled {
			e.current = output.Text
		} else if e.index >= 0 && e.index < len(e.history) {
			e.current = e.history[e.index]
		} else {
			e.current = ""
		}
		e.saveHistory()
		e.mu.Unlock()

		// При отмене не копируем в буфер обмена
		if output.Cancelled {
			return
		}

		// Копируем в буфер обмена
		if err := clipboard.Copy(output.Text); err != nil {
			fmt.Printf(msg.EditorProcessError+"\n", err)
			return
		}
		fmt.Printf(msg.EditorCopiedToClip+"\n", len(output.Text))
	}()

	return nil
}

// NavigateHistory навигация по истории
func (e *Editor) NavigateHistory(direction int) string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.history) == 0 {
		return e.current
	}

	newIndex := e.index + direction
	if newIndex < 0 {
		newIndex = 0
	}
	if newIndex >= len(e.history) {
		newIndex = len(e.history) - 1
	}

	e.index = newIndex
	e.current = e.history[e.index]

	return e.current
}

// Clear очищает редактор
func (e *Editor) Clear() {
	e.mu.Lock()
	lang := e.lang
	e.current = ""
	e.mu.Unlock()

	msg := i18n.Get(lang)
	fmt.Println(msg.EditorCleared)
}
