package lemonade

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Client клиент Lemonade API
type Client struct {
	url        string
	httpClient *http.Client
}

// TranscribeRequest запрос на транскрибацию
type TranscribeRequest struct {
	Model    string `json:"model"`
	Language string `json:"language,omitempty"`
}

// TranscribeResponse ответ сервера
type TranscribeResponse struct {
	Text string `json:"text"`
}

// TranscribeResult результат транскрибации со статистикой
type TranscribeResult struct {
	Text        string
	Duration    float64 // Длительность аудио в секундах
	ProcessTime float64 // Время обработки в секундах
	Speed       float64 // Скорость (x реального времени)
	Backend     string  // Бэкенд (cpu/dxml)
}

// ModelInfo информация о модели из Lemonade API
type ModelInfo struct {
	ID         string   `json:"id"`
	Recipe     string   `json:"recipe"`
	Size       float64  `json:"size"`
	Labels     []string `json:"labels"`
	Downloaded bool     `json:"downloaded"`
}

// NewClient создаёт новый клиент
func NewClient(url string) *Client {
	return &Client{
		url: url,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // Долгий таймаут для транскрибации
		},
	}
}

// Transcribe отправляет аудио на транскрибацию
func (c *Client) Transcribe(audioPath string, model, language, prompt string, temperature float64) (string, error) {
	result, err := c.TranscribeWithStats(audioPath, model, language, prompt, temperature)
	if err != nil {
		return "", err
	}
	return result.Text, nil
}

// TranscribeWithStats отправляет аудио на транскрибацию и возвращает статистику
func (c *Client) TranscribeWithStats(audioPath string, model, language, prompt string, temperature float64) (*TranscribeResult, error) {
	startTime := time.Now()
	
	// Проверяем существование файла
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("audio file not found: %s", audioPath)
	}
	
	// Получаем длительность аудио из WAV файла
	duration, err := getWavDuration(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio duration: %w", err)
	}

	// Создаём multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем аудио файл
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Добавляем модель
	if err := writer.WriteField("model", model); err != nil {
		return nil, fmt.Errorf("failed to write model field: %w", err)
	}

	// Добавляем язык (опционально)
	if language != "" {
		if err := writer.WriteField("language", language); err != nil {
			return nil, fmt.Errorf("failed to write language field: %w", err)
		}
	}

	// Добавляем промпт (опционально)
	if prompt != "" {
		if err := writer.WriteField("prompt", prompt); err != nil {
			return nil, fmt.Errorf("failed to write prompt field: %w", err)
		}
	}

	// Добавляем температуру (если не дефолтная)
	if temperature > 0 {
		if err := writer.WriteField("temperature", fmt.Sprintf("%.2f", temperature)); err != nil {
			return nil, fmt.Errorf("failed to write temperature field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Создаём HTTP запрос
	url := c.url + "/api/v1/audio/transcriptions"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error (%d): %s", resp.StatusCode, string(respBody))
	}

	// Парсим JSON ответ
	var result TranscribeResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	processTime := time.Since(startTime).Seconds()
	speed := duration / processTime

	return &TranscribeResult{
		Text:        result.Text,
		Duration:    duration,
		ProcessTime: processTime,
		Speed:       speed,
		Backend:     "cpu", // По умолчанию
	}, nil
}

// CheckHealth проверяет доступность сервера
func (c *Client) CheckHealth() (bool, error) {
	url := c.url + "/health"
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to get health: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return true, nil
}

// GetWhisperModels возвращает список Whisper моделей (включая не скачанные)
func (c *Client) GetWhisperModels() ([]ModelInfo, error) {
	url := c.url + "/api/v1/models?show_all=true"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
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

	var result struct {
		Data []ModelInfo `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Фильтруем только Whisper модели (recipe = "whispercpp" или id содержит "whisper"/"Whisper")
	var whisperModels []ModelInfo
	for _, m := range result.Data {
		if m.Recipe == "whispercpp" || strings.Contains(strings.ToLower(m.ID), "whisper") {
			whisperModels = append(whisperModels, m)
		}
	}

	return whisperModels, nil
}

// PullProgress прогресс скачивания модели
type PullProgress struct {
	Status  string `json:"status"`  // "progress", "complete", "error"
	File    string `json:"file"`
	Percent int    `json:"percent"`
	Error   string `json:"error"`
}

// PullModel скачивает модель с прогрессом через SSE
// onProgress вызывается для каждого события прогресса
func (c *Client) PullModel(modelName string, onProgress func(PullProgress)) error {
	url := c.url + "/api/v1/pull"

	payload := map[string]interface{}{
		"model_name": modelName,
		"stream":     true,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Minute}
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(respBody))
	}

	// Читаем SSE поток
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		jsonData := strings.TrimPrefix(line, "data: ")

		var progress PullProgress
		if err := json.Unmarshal([]byte(jsonData), &progress); err != nil {
			continue
		}

		if onProgress != nil {
			onProgress(progress)
		}

		if progress.Status == "error" {
			return fmt.Errorf("pull failed: %s", progress.Error)
		}
		if progress.Status == "complete" {
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream read error: %w", err)
	}

	return nil
}

// LoadModel загружает модель в память (скачивает если нужно)
func (c *Client) LoadModel(modelName string) error {
	url := c.url + "/api/v1/load"

	payload := map[string]string{"model_name": modelName}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Загрузка модели может занять много времени (скачивание + загрузка в память)
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to load model: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// getWavDuration возвращает длительность WAV файла в секундах
func getWavDuration(filePath string) (float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Читаем WAV заголовок (44 байта)
	header := make([]byte, 44)
	if _, err := io.ReadFull(file, header); err != nil {
		return 0, err
	}

	// Читаем параметры из заголовка
	// Sample rate: байты 24-27 (little endian)
	sampleRate := uint32(header[24]) | uint32(header[25])<<8 | uint32(header[26])<<16 | uint32(header[27])<<24
	
	// Data size: байты 40-43 (little endian)
	dataSize := uint32(header[40]) | uint32(header[41])<<8 | uint32(header[42])<<16 | uint32(header[43])<<24
	
	// Bits per sample: байты 34-35
	bitsPerSample := uint32(header[34]) | uint32(header[35])<<8
	
	// Num channels: байты 22-23
	numChannels := uint32(header[22]) | uint32(header[23])<<8
	
	// Вычисляем длительность
	bytesPerSample := bitsPerSample / 8
	numSamples := dataSize / bytesPerSample
	duration := float64(numSamples) / float64(sampleRate*numChannels)
	
	return duration, nil
}
