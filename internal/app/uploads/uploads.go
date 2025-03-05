package uploads

import (
	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/app/file"
	"github.com/serhiirubets/rubeticket/internal/app/fileuploader"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/middleware"
	"net/http"
)

type HandlerDeps struct {
	Logger       log.ILogger
	Config       *config.Config
	FileUploader *fileuploader.FileUploader
}

type Handler struct {
	Logger       log.ILogger
	Config       *config.Config
	FileUploader *fileuploader.FileUploader
}

func NewUploadsHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{
		Logger:       deps.Logger,
		Config:       deps.Config,
		FileUploader: deps.FileUploader,
	}

	router.Handle("GET /uploads/{fileName}", middleware.Auth(handler.GetPhoto(), deps.Config, deps.Logger))
}

// GetPhoto godoc
// @Summary Get a file by path
// @Description Retrieve a file by its path for the authenticated user
// @Tags Account
// @Security ApiKeyAuth
// @Produce application/octet-stream
// @Param fileName path string true "File path (e.g., b60b4dd7-6dda-49fc-830f-020fa5fe4817.png)"
// @Success 200 {file} file "File content"
// @Failure 400 {object} string "Invalid file path"
// @Failure 401 {object} string "Not authorized"
// @Failure 403 {object} string "Forbidden or file not found"
// @Failure 500 {object} string "Internal server error"
// @Router /uploads/{fileName} [get]
func (handler *Handler) GetPhoto() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authData, _ := middleware.GetAuthData(r)

		fileName := r.PathValue("fileName")
		if fileName == "" {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		// Check that user is a file owner
		var fileModel file.File
		if err := handler.FileUploader.DB.
			Where("file_path = ? AND user_id = ?", fileName, authData.UserID).
			First(&fileModel).Error; err != nil {
			handler.Logger.Warn("Unauthorized file access attempt", "file_path", fileName, "user_id", authData.UserID)
			http.Error(w, "Forbidden or file not found", http.StatusForbidden)
			return
		}

		filePath := "uploads/" + fileModel.FilePath

		http.ServeFile(w, r, filePath)
	})
}
