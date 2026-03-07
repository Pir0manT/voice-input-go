package recorder

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gordonklaus/portaudio"
)

// Recorder аудио рекордер
type Recorder struct {
	mu          sync.Mutex
	samplesMu   sync.Mutex // отдельный mutex для samples — используется в audioCallback
	recording   atomic.Bool // атомарный флаг для безопасного чтения из callback
	stream      *portaudio.Stream
	samples     []float32
	sampleRate  int
	channels    int
	startTime   time.Time
	language    string // Язык для сообщений
}

// Config конфигурация рекордера
type Config struct {
	SampleRate int
	Channels   int
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		SampleRate: 16000, // 16kHz - стандарт для Whisper
		Channels:   1,     // Mono
	}
}

// New создаёт новый рекордер
func New() *Recorder {
	return &Recorder{
		sampleRate: 16000,
		channels:   1,
		samples:    make([]float32, 0),
		language:   "ru", // По умолчанию русский
	}
}

// NewWithConfig создаёт рекордер с конфигурацией
func NewWithConfig(cfg Config, lang string) *Recorder {
	return &Recorder{
		sampleRate: cfg.SampleRate,
		channels:   cfg.Channels,
		samples:    make([]float32, 0),
		language:   lang,
	}
}

// SetLanguage устанавливает язык для сообщений
func (r *Recorder) SetLanguage(lang string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.language = lang
}

// Start начинает запись
func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.recording.Load() {
		return fmt.Errorf("already recording")
	}

	// Инициализируем portaudio
	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize portaudio: %w", err)
	}

	// Очищаем предыдущие сэмплы
	r.samplesMu.Lock()
	r.samples = make([]float32, 0)
	r.samplesMu.Unlock()

	// Создаём поток записи
	var err error
	r.stream, err = portaudio.OpenDefaultStream(
		r.channels,    // input channels
		0,             // output channels
		float64(r.sampleRate),
		0,             // frames per buffer (0 = use default)
		r.audioCallback,
	)
	if err != nil {
		portaudio.Terminate()
		return fmt.Errorf("failed to open stream: %w", err)
	}

	// Запускаем поток
	if err := r.stream.Start(); err != nil {
		r.stream.Close()
		portaudio.Terminate()
		return fmt.Errorf("failed to start stream: %w", err)
	}

	r.recording.Store(true)
	r.startTime = time.Now()

	return nil
}

// audioCallback вызывается PortAudio из аудиопотока ОС.
// Не должен захватывать r.mu — иначе deadlock с Stop().
// Используем atomic для флага и отдельный samplesMu для буфера.
func (r *Recorder) audioCallback(in []float32, out []float32) {
	if !r.recording.Load() {
		return
	}

	r.samplesMu.Lock()
	r.samples = append(r.samples, in...)
	r.samplesMu.Unlock()
}

// StopResult результат остановки записи
type StopResult struct {
	FilePath       string  // Путь к WAV файлу
	TrimmedSeconds float64 // Сколько секунд тишины обрезано (0 = ничего не обрезано)
}

// ErrSilentRecording запись содержит только тишину
var ErrSilentRecording = fmt.Errorf("silent recording")

// Stop останавливает запись и возвращает результат
func (r *Recorder) Stop() (*StopResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.recording.Load() {
		return nil, fmt.Errorf("not recording")
	}

	// Сначала сбрасываем флаг — callback перестанет писать в samples
	r.recording.Store(false)

	// Останавливаем поток (теперь callback не заблокируется на samplesMu)
	if r.stream != nil {
		if err := r.stream.Stop(); err != nil {
			// Пытаемся освободить ресурсы даже при ошибке
			r.stream.Close()
			portaudio.Terminate()
			return nil, fmt.Errorf("failed to stop stream: %w", err)
		}
		if err := r.stream.Close(); err != nil {
			portaudio.Terminate()
			return nil, fmt.Errorf("failed to close stream: %w", err)
		}
		portaudio.Terminate()
	}

	// Обрезаем тишину в начале и конце записи
	originalLen := len(r.samples)
	trimmed := trimSilence(r.samples, r.sampleRate)
	if len(trimmed) == 0 {
		return nil, ErrSilentRecording
	}
	r.samples = trimmed

	var trimmedSec float64
	if trimmedCount := originalLen - len(trimmed); trimmedCount > 0 {
		trimmedSec = float64(trimmedCount) / float64(r.sampleRate)
	}

	// Сохраняем в WAV файл (callback уже не пишет — samplesMu не нужен)
	filename, err := r.saveToWAV()
	if err != nil {
		return nil, fmt.Errorf("failed to save WAV: %w", err)
	}

	return &StopResult{
		FilePath:       filename,
		TrimmedSeconds: trimmedSec,
	}, nil
}

// saveToWAV сохраняет сэмплы в WAV файл
func (r *Recorder) saveToWAV() (string, error) {
	// Используем AppData для временных файлов
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		// Фоллбэк на системную temp директорию
		appDataDir = os.TempDir()
	}
	
	tempDir := filepath.Join(appDataDir, "voice-input-go", "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Генерируем имя файла с timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(tempDir, fmt.Sprintf("recording_%s.wav", timestamp))

	// Открываем файл
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// WAV заголовок
	// RIFF header
	file.Write([]byte("RIFF"))
	file.Write(uint32ToBytes(36 + uint32(len(r.samples)*2))) // File size - 8
	file.Write([]byte("WAVE"))

	// fmt subchunk
	file.Write([]byte("fmt "))
	file.Write(uint32ToBytes(16))            // Subchunk1Size (16 for PCM)
	file.Write(uint16ToBytes(1))             // AudioFormat (1 = PCM)
	file.Write(uint16ToBytes(uint16(r.channels))) // NumChannels
	file.Write(uint32ToBytes(uint32(r.sampleRate))) // SampleRate
	file.Write(uint32ToBytes(uint32(r.sampleRate * r.channels * 2))) // ByteRate
	file.Write(uint16ToBytes(uint16(r.channels * 2))) // BlockAlign
	file.Write(uint16ToBytes(16))            // BitsPerSample (16 = 2 bytes)

	// data subchunk
	file.Write([]byte("data"))
	file.Write(uint32ToBytes(uint32(len(r.samples) * 2))) // Subchunk2Size

	// Записываем сэмплы (конвертируем float32 → int16)
	for _, sample := range r.samples {
		// Конвертируем float32 [-1.0, 1.0] → int16 [-32768, 32767]
		s := int16(sample * 32767.0)
		file.Write(int16ToBytes(s))
	}

	return filename, nil
}

// IsRecording возвращает статус записи
func (r *Recorder) IsRecording() bool {
	return r.recording.Load()
}

// GetDuration возвращает длительность текущей записи
func (r *Recorder) GetDuration() time.Duration {
	if !r.recording.Load() {
		return 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return time.Since(r.startTime)
}

// GetSampleCount возвращает количество сэмплов
func (r *Recorder) GetSampleCount() int {
	r.samplesMu.Lock()
	defer r.samplesMu.Unlock()
	return len(r.samples)
}

// Clear очищает буфер сэмплов
func (r *Recorder) Clear() {
	r.samplesMu.Lock()
	defer r.samplesMu.Unlock()
	r.samples = make([]float32, 0)
}

// Вспомогательные функции для конвертации чисел в байты

func uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
	}
}

func uint16ToBytes(n uint16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}

func int16ToBytes(n int16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}

// trimSilence обрезает тишину в начале и конце записи.
// Использует RMS (среднеквадратичное) по окнам ~20ms.
// Оставляет небольшой запас (tailMs) чтобы не обрезать окончания слов.
func trimSilence(samples []float32, sampleRate int) []float32 {
	const (
		thresholdDB = -40.0 // Порог тишины в dB
		tailMs      = 200   // Запас после последнего звука (мс)
		windowMs    = 20    // Размер окна анализа (мс)
	)

	windowSize := sampleRate * windowMs / 1000 // ~320 сэмплов при 16kHz
	tailSize := sampleRate * tailMs / 1000     // ~3200 сэмплов при 16kHz

	if len(samples) < windowSize {
		return samples
	}

	// Порог RMS из dB: threshold = 10^(dB/20)
	threshold := float32(math.Pow(10, thresholdDB/20.0)) // ~0.01

	// Ищем первое окно со звуком (с начала)
	startIdx := 0
	for i := 0; i+windowSize <= len(samples); i += windowSize {
		if rmsWindow(samples[i:i+windowSize]) > threshold {
			// Отступаем назад на одно окно для запаса
			startIdx = i - windowSize
			if startIdx < 0 {
				startIdx = 0
			}
			break
		}
		// Если дошли до конца — вся запись тихая
		if i+windowSize >= len(samples) {
			return nil
		}
	}

	// Ищем последнее окно со звуком (с конца)
	endIdx := len(samples)
	for i := len(samples) - windowSize; i >= 0; i -= windowSize {
		if rmsWindow(samples[i:i+windowSize]) > threshold {
			// Добавляем хвост чтобы не обрезать окончание слова
			endIdx = i + windowSize + tailSize
			if endIdx > len(samples) {
				endIdx = len(samples)
			}
			break
		}
	}

	if startIdx >= endIdx {
		return nil
	}

	return samples[startIdx:endIdx]
}

// rmsWindow вычисляет RMS (среднеквадратичное) для окна сэмплов
func rmsWindow(window []float32) float32 {
	var sum float32
	for _, s := range window {
		sum += s * s
	}
	return float32(math.Sqrt(float64(sum / float32(len(window)))))
}
