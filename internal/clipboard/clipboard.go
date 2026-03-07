package clipboard

// Copy копирует текст в буфер обмена
func Copy(text string) error {
	return copyPlatform(text)
}

// Paste возвращает текст из буфера обмена
func Paste() (string, error) {
	return pastePlatform()
}
