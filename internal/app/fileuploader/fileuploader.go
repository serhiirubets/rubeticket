package fileuploader

import (
	"errors"
	"github.com/google/uuid"
	"github.com/serhiirubets/rubeticket/internal/app/file"
	"github.com/serhiirubets/rubeticket/internal/pkg/db"
	"github.com/serhiirubets/rubeticket/internal/pkg/filestorage"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"mime/multipart"
	"os"
	"strings"
)

type Deps struct {
	Logger         log.ILogger
	DB             *db.Db
	Storage        filestorage.Storage
	AllowedTypes   []string
	MaxSizeMB      int64
	FileRepository *file.Repository
}

type FileUploader struct {
	Logger         log.ILogger
	DB             *db.Db
	Storage        filestorage.Storage
	AllowedTypes   []string
	MaxSizeMB      int64
	FileRepository *file.Repository
}

func NewFileUploader(deps *Deps) *FileUploader {
	return &FileUploader{
		Logger:         deps.Logger,
		DB:             deps.DB,
		Storage:        deps.Storage,
		AllowedTypes:   deps.AllowedTypes,
		MaxSizeMB:      deps.MaxSizeMB,
		FileRepository: deps.FileRepository,
	}
}

func (f *FileUploader) UploadFile(uploadFile multipart.File, header *multipart.FileHeader, userID uint, purpose string) (*file.File, error) {
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
		return nil, errors.New("invalid file type")
	}

	fileUUID := uuid.New().String()
	filePath, err := f.Storage.SaveFile(uploadFile, header, fileUUID)
	if err != nil {
		f.Logger.Error("Failed to save file to storage", "error", err.Error())
		return nil, err
	}

	fileModel := &file.File{
		UUID:     fileUUID,
		UserID:   userID,
		FilePath: filePath,
		Purpose:  purpose,
	}

	// TODO: Move to the file repository
	if err := f.DB.Create(fileModel).Error; err != nil {
		f.Logger.Error("Failed to save file metadata to DB", "error", err.Error())
		os.Remove(filePath)
		return nil, err
	}

	return fileModel, nil
}
