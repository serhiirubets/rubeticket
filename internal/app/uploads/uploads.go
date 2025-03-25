package uploads

import (
	"net/http"

	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/app/file"
	"github.com/serhiirubets/rubeticket/internal/app/fileuploader"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/middleware"
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
	router.HandleFunc("GET /uploads/{fileName}", handler.GetPhoto())
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
	return func(w http.ResponseWriter, r *http.Request) {
		authData, err := middleware.GetAuthData(r)
		if err != nil {
			handler.Logger.Error("Error getting auth data", "error", err.Error())
			return
		}

		fileName := r.PathValue("fileName")
		if fileName == "" {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		// Check if file exists in database
		var fileModel file.File
		var userID uint
		if authData.UserID != 0 {
			userID = authData.UserID
		}
		query := handler.FileUploader.DB.Where("file_path = ?", fileName)
		if userID != 0 {
			query = query.Where("user_id = ?", userID)
		}
		if err := query.First(&fileModel).Error; err != nil {
			handler.Logger.Warn("File not found in database", "file_path", fileName, "user_id", userID, "error", err.Error())
			http.Error(w, "Forbidden or file not found", http.StatusForbidden)
			return
		}

		filePath := "uploads/" + fileModel.FilePath

		// Check if file exists on disk
		if _, err := http.Dir(".").Open(filePath); err != nil {
			handler.Logger.Warn("File not found on disk", "file_path", filePath, "error", err.Error())
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, filePath)
	}
}