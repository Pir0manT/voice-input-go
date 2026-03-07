package logview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/Pir0manT/voice-input-go/internal/i18n"
)

// LogView менеджер окна логов
type LogView struct {
	mu      sync.Mutex
	running bool
	lang    string
	logPath string
}

// New создаёт менеджер логов
func New(logPath, lang string) *LogView {
	return &LogView{
		logPath: logPath,
		lang:    lang,
	}
}

// SetLanguage обновляет язык
func (lv *LogView) SetLanguage(lang string) {
	lv.mu.Lock()
	defer lv.mu.Unlock()
	lv.lang = lang
}

// IsRunning проверяет открыто ли окно
func (lv *LogView) IsRunning() bool {
	lv.mu.Lock()
	defer lv.mu.Unlock()
	return lv.running
}

// Show открывает окно логов
func (lv *LogView) Show() {
	lv.mu.Lock()
	msg := i18n.Get(lv.lang)

	if lv.running {
		lv.mu.Unlock()
		fmt.Println(msg.LogsAlreadyOpen)
		return
	}

	lv.running = true
	input := LogsInput{
		LogPath: lv.logPath,
		Lang:    lv.lang,
	}
	lv.mu.Unlock()

	fmt.Println(msg.LogsOpening)

	go func() {
		defer func() {
			lv.mu.Lock()
			lv.running = false
			lv.mu.Unlock()
		}()

		if err := launchLogsWindow(input); err != nil {
			fmt.Printf(msg.LogsProcessError+"\n", err)
		}
	}()
}

// launchLogsWindow запускает GUI процесс логов (сам себя с --logs)
func launchLogsWindow(input LogsInput) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	inputData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	cmd := exec.Command(exePath, "--logs")
	cmd.Stdin = bytes.NewReader(inputData)
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start log viewer: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return nil
			}
		}
		return fmt.Errorf("log viewer process error: %w", err)
	}

	return nil
}
