package filestorage

import (
	"io"
	"mime/multipart"
)

type Storage interface {
	SaveFile(file multipart.File, header *multipart.FileHeader, uuid string) (string, error)
	GetFile(filePath string) (io.ReadCloser, error)
}
