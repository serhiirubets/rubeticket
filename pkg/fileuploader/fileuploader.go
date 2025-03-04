package fileuploader

import (
	"errors"
	"github.com/serhiirubets/rubeticket/pkg/log"
	"golang.org/x/mod/sumdb/storage"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type FileUploaderDeps struct {
	Logger       log.ILogger
	UploadDir    string
	MaxSizeMB    int64
	AllowedTypes []string
	Storage      storage.Storage
}

type FileUploader struct {
	Logger       log.ILogger
	UploadDir    string
	MaxSizeMB    int64
	AllowedTypes []string
}

func NewFileUploader(deps *FileUploaderDeps) *FileUploader {
	return &FileUploader{
		Logger:       deps.Logger,
		UploadDir:    deps.UploadDir,
		MaxSizeMB:    deps.MaxSizeMB,
		AllowedTypes: deps.AllowedTypes,
	}
}

// UploadFile — универсальная функция загрузки
// entityType: "user", "group", "concert"
// entityID: ID сущности (например, user.ID)
// purpose: "profile", "cover", "poster"
func (f *FileUploader) UploadFile(file multipart.File, header *multipart.FileHeader, entityType, entityID, purpose string) (string, error) {
	// Проверка MIME-типа (только изображения)
	contentType := header.Header.Get("Content-Type")
	allowed := false
	for _, allowedType := range f.AllowedTypes {
		if strings.HasPrefix(contentType, allowedType) {
			allowed = true
			break
		}
	}
	if !allowed {
		f.Logger.Warn("Invalid file type attempted", "type", contentType)
		return "", errors.New("invalid file type")
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(header.Filename)
	filename := entityType + "_" + entityID + "_" + purpose + ext
	fPath := filepath.Join(f.UploadDir, entityType, filename)

	// Создаём директорию
	if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
		f.Logger.Error("Failed to create directory", "error", err.Error())
		return "", err
	}

	// Создаём файл
	dst, err := os.Create(fPath)
	if err != nil {
		f.Logger.Error("Failed to create file", "error", err.Error())
		return "", err
	}
	defer dst.Close()

	// Копируем содержимое
	if _, err := io.Copy(dst, file); err != nil {
		f.Logger.Error("Failed to save file", "error", err.Error())
		return "", err
	}

	return fPath, nil
}
