package filestorage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BaseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{BaseDir: baseDir}
}

func (s *LocalStorage) SaveFile(file multipart.File, header *multipart.FileHeader, uuid string) (string, error) {
	ext := filepath.Ext(header.Filename)
	fileName := uuid + ext
	filePath := filepath.Join(s.BaseDir, fileName)

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return "", err
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *LocalStorage) GetFile(filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}
