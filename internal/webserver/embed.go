package webserver

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
)

//go:embed ui/*
var uiFiles embed.FS

// GetFileSystem возвращает файловую систему для UI файлов
func GetFileSystem() http.FileSystem {
	// Получаем подкаталог ui
	sub, err := fs.Sub(uiFiles, "ui")
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}

// GetFile возвращает содержимое файла
func GetFile(path string) ([]byte, error) {
	return uiFiles.ReadFile("ui/" + path)
}

// FileExists проверяет существование файла
func FileExists(path string) bool {
	_, err := uiFiles.ReadFile("ui/" + path)
	return err == nil
}

// SaveFile сохраняет файл (для внешних файлов, не embed)
func SaveFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
