package fastflowlm

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	mu         sync.Mutex
	serverProc *os.Process // процесс FLM-сервера, если мы его запускали
)

// IsInstalled проверяет наличие flm в PATH
func IsInstalled() bool {
	_, err := exec.LookPath("flm")
	return err == nil
}

// EnsureRunning проверяет health → если не отвечает → запускает flm serve <llmModel> --asr 1
func EnsureRunning(url, llmModel string) error {
	// Сначала проверяем, работает ли сервер
	if checkHealthQuick(url) {
		return nil
	}

	// Сервер не отвечает — пробуем запустить
	if !IsInstalled() {
		return fmt.Errorf("flm not found in PATH")
	}

	mu.Lock()
	defer mu.Unlock()

	// Повторная проверка под мьютексом
	if checkHealthQuick(url) {
		return nil
	}

	// Запускаем flm serve [llmModel] --asr 1
	var cmd *exec.Cmd
	if llmModel != "" {
		cmd = exec.Command("flm", "serve", llmModel, "--asr", "1")
	} else {
		cmd = exec.Command("flm", "serve", "--asr", "1")
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	setSysProcAttr(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start flm serve: %w", err)
	}

	serverProc = cmd.Process

	// Ждём готовности (poll health до 30 сек)
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(500 * time.Millisecond)
		if checkHealthQuick(url) {
			return nil
		}
	}

	return fmt.Errorf("flm server did not become ready within 30 seconds")
}

// StopServer убивает процесс FLM-сервера и всё дерево дочерних процессов
func StopServer() {
	mu.Lock()
	defer mu.Unlock()

	if serverProc == nil {
		return
	}

	// Убиваем всё дерево процессов (NPU runtime и т.д.)
	_ = killProcessTree(serverProc)

	// Ждём завершения, чтобы ресурсы (NPU, порты) точно освободились
	done := make(chan struct{})
	go func() {
		serverProc.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		// Если не завершился за 5 сек — принудительно kill
		_ = serverProc.Kill()
	}

	serverProc = nil
}

// checkHealthQuick быстрая проверка доступности сервера
func checkHealthQuick(url string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url + "/v1/models")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
