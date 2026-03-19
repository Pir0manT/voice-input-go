package fastflowlm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ModelInfo информация о модели FastFlowLM
type ModelInfo struct {
	ID         string  `json:"id"`      // ID для API (whisper-v3)
	Tag        string  `json:"tag"`     // Тег для pull/list (whisper-v3:turbo)
	Name       string  `json:"name"`
	Size       float64 `json:"size"`    // Размер в ГБ
	Downloaded bool    `json:"downloaded"`
}

// PullProgress прогресс скачивания модели
type PullProgress struct {
	Status    string `json:"status"`
	Total     int64  `json:"total"`
	Completed int64  `json:"completed"`
	Percent   int    `json:"-"`
}

// knownWhisperModels известные Whisper-модели FastFlowLM
var knownWhisperModels = []ModelInfo{
	{ID: "whisper-v3", Tag: "whisper-v3:turbo", Name: "Whisper V3 Turbo", Size: 0.6},
}

// GetTagByID возвращает pull-тег модели по её API ID
func GetTagByID(modelID string) string {
	for _, m := range knownWhisperModels {
		if m.ID == modelID {
			return m.Tag
		}
	}
	return modelID // fallback
}

// KnownWhisperModels возвращает копию списка известных Whisper-моделей (для fallback)
func KnownWhisperModels() []ModelInfo {
	result := make([]ModelInfo, len(knownWhisperModels))
	copy(result, knownWhisperModels)
	return result
}

// ollamaTagsResponse ответ GET /api/tags
type ollamaTagsResponse struct {
	Models []ollamaModel `json:"models"`
}

type ollamaModel struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// GetModels возвращает список Whisper-моделей (объединяет установленные с известными)
func GetModels(serverURL string) ([]ModelInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(serverURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tagsResp ollamaTagsResponse
	if err := json.Unmarshal(body, &tagsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Собираем множество установленных моделей
	installed := make(map[string]bool)
	for _, m := range tagsResp.Models {
		installed[m.Name] = true
	}

	// Объединяем с известным списком Whisper-моделей
	result := make([]ModelInfo, len(knownWhisperModels))
	copy(result, knownWhisperModels)
	for i := range result {
		if installed[result[i].Tag] {
			result[i].Downloaded = true
		}
	}

	return result, nil
}

// PullModel скачивает модель с прогрессом через NDJSON стриминг (Ollama pull API)
func PullModel(serverURL, modelTag string, onProgress func(PullProgress)) error {
	payload := map[string]interface{}{
		"name":   modelTag,
		"stream": true,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Minute}
	resp, err := client.Post(serverURL+"/api/pull", "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(respBody))
	}

	// Читаем NDJSON поток
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var progress PullProgress
		if err := json.Unmarshal([]byte(line), &progress); err != nil {
			continue
		}

		// Вычисляем процент
		if progress.Total > 0 {
			progress.Percent = int(progress.Completed * 100 / progress.Total)
		}

		if onProgress != nil {
			onProgress(progress)
		}

		if progress.Status == "error" {
			return fmt.Errorf("pull failed: %s", line)
		}
		if progress.Status == "success" {
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream read error: %w", err)
	}

	return nil
}
