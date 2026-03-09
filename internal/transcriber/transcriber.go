package transcriber

// Result результат транскрибации со статистикой
type Result struct {
	Text        string
	Duration    float64 // Длительность аудио в секундах
	ProcessTime float64 // Время обработки в секундах
	Speed       float64 // Скорость (x реального времени)
	Backend     string  // Бэкенд (lemonade/whisper-api)
}

// Transcriber интерфейс для бэкендов транскрибации
type Transcriber interface {
	// TranscribeWithStats отправляет аудио на транскрибацию и возвращает результат со статистикой
	TranscribeWithStats(audioPath string) (*Result, error)

	// CheckHealth проверяет доступность сервера
	CheckHealth() (bool, error)
}
