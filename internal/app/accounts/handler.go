package accounts

import (
	"mime/multipart"
	"net/http"

	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/app/file"
	"github.com/serhiirubets/rubeticket/internal/app/fileuploader"
	"github.com/serhiirubets/rubeticket/internal/app/users"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/middleware"
	"github.com/serhiirubets/rubeticket/internal/pkg/req"
	"github.com/serhiirubets/rubeticket/internal/pkg/res"
)

type AccountHandlerDeps struct {
	UserRepository *users.UserRepository
	Logger         log.ILogger
	Config         *config.Config
	FileUploader   *fileuploader.FileUploader
}

type AccountHandler struct {
	UserRepository *users.UserRepository
	Logger         log.ILogger
	Config         *config.Config
	FileUploader   *fileuploader.FileUploader
}

func NewAccountHandler(router *http.ServeMux, deps *AccountHandlerDeps) {
	handler := &AccountHandler{
		UserRepository: deps.UserRepository,
		Logger:         deps.Logger,
		Config:         deps.Config,
		FileUploader:   deps.FileUploader,
	}

	router.HandleFunc("GET /account", handler.GetAccount())
	router.HandleFunc("PATCH /account", handler.UpdateAccountPatch())
	router.HandleFunc("PUT /account", handler.UpdateAccountPut())
	router.HandleFunc("POST /account/photo", handler.UploadPhoto())
}

// GetAccount godoc
// @Summary Get account info
// @Description Return information about current user
// @Tags Account
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} GetAccountResponse "Success"
// @Failure 401 {object} string "Not authorized"
// @Router /v1/account [get]
func (handler *AccountHandler) GetAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authData, _ := middleware.GetAuthData(r)

		user, userErr := handler.UserRepository.GetByEmail(authData.Email)

		if userErr != nil {
			handler.Logger.Error("Error getting user by email", userErr.Error())
			res.Json(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		var photo file.File
		photoErr := handler.FileUploader.DB.
			Where("user_id = ? AND purpose = ?", authData.UserID, "profile").
			Last(&photo).Error

		photoUrl := ""

		if photoErr == nil {
			photoUrl = photo.FilePath
		}

		body := GetAccountResponse{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Gender:    user.Gender,
			Birthday:  user.Birthday,
			PhotoUrl:  photoUrl,
		}

		res.Json(w, body, http.StatusOK)
	}
}

// UpdateAccountPatch godoc
// @Summary Update account info
// @Description Update specific fields of the current user's account
// @Tags Account
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body UpdateAccountRequestPatch true "Fields to update"
// @Success 200 {object} UpdateAccountResponse "Success"
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} string "Not authorized"
// @Failure 500 {object} string "Internal server error"
// @Router /v1/account [patch]
func (handler *AccountHandler) UpdateAccountPatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authData, _ := middleware.GetAuthData(r)

		body, err := req.HandleBody[UpdateAccountRequestPatch](&w, r)

		if err != nil {
			return
		}

		user, err := handler.UserRepository.GetByEmail(authData.Email)
		if err != nil {
			handler.Logger.Error("Error getting user by email", err.Error())
			res.Json(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		updates := make(map[string]interface{})
		if body.FirstName != nil {
			updates["first_name"] = *body.FirstName
		}
		if body.LastName != nil {
			updates["last_name"] = *body.LastName
		}
		if body.Gender != nil {
			updates["gender"] = *body.Gender
		}
		if body.Birthday != nil {
			updates["birthday"] = *body.Birthday
		}
		if body.Address != nil {
			updates["address"] = *body.Address
		}

		if len(updates) == 0 {
			res.Json(w, "No fields to update", http.StatusBadRequest)
			return
		}

		if err := handler.UserRepository.Update(user, updates); err != nil {
			handler.Logger.Error("Failed to update user", err.Error())
			res.Json(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := UpdateAccountResponse{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Gender:    user.Gender,
			Birthday:  user.Birthday,
		}

		res.Json(w, response, http.StatusOK)
	}
}

// UpdateAccountPut godoc
// @Summary Update full account info
// @Description Replace all fields of the current user's account with provided values
// @Tags Account
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body UpdateAccountRequestPut true "Full account details to update"
// @Success 200 {object} UpdateAccountResponse "Success"
// @Failure 400 {object} string "Invalid request body"
// @Failure 401 {object} string "Not authorized"
// @Failure 500 {object} string "Internal server error"
// @Router /v1/account [put]
func (handler *AccountHandler) UpdateAccountPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authData, _ := middleware.GetAuthData(r)

		body, errBody := req.HandleBody[UpdateAccountRequestPut](&w, r)

		if errBody != nil {
			return
		}

		user, errUser := handler.UserRepository.GetByEmail(authData.Email)
		if errUser != nil {
			handler.Logger.Error("Error getting user by email", errUser.Error())
			res.Json(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		user.FirstName = body.FirstName
		user.LastName = body.LastName
		user.Gender = body.Gender
		user.Birthday = body.Birthday

		if err := handler.UserRepository.DB.Save(&user).Error; err != nil {
			handler.Logger.Error("Failed to update user", err.Error())
			res.Json(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := UpdateAccountResponse{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Gender:    user.Gender,
			Birthday:  user.Birthday,
		}
		res.Json(w, response, http.StatusOK)
	}
}

// UploadPhoto godoc
// @Summary Upload a photo
// @Description Upload a photo file for the current user
// @Tags Account
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param photo formData file true "Photo file to upload"
// @Success 200 {object} map[string]string "Success"
// @Failure 400 {object} string "Bad request"
// @Failure 401 {object} string "Not authorized"
// @Failure 500 {object} string "Internal server error"
// @Router /v1/account/photo [post]
func (handler *AccountHandler) UploadPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authData, _ := middleware.GetAuthData(r)

		r.Body = http.MaxBytesReader(w, r.Body, handler.FileUploader.MaxSizeMB<<20)
		if err := r.ParseMultipartForm(handler.FileUploader.MaxSizeMB << 20); err != nil {
			handler.Logger.Error("Failed to parse multipart form", "error", err.Error())
			res.Json(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		photo, header, err := r.FormFile("photo")
		if err != nil {
			handler.Logger.Error("Failed to get file from form", "error", err.Error())
			res.Json(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer func(photo multipart.File) {
			err := photo.Close()
			if err != nil {
				handler.Logger.Error("Failed to close file", "error", err.Error())
			}
		}(photo)

		fileModel, err := handler.FileUploader.UploadFile(photo, header, authData.UserID, "profile")
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, map[string]string{"message": "Photo uploaded", "uuid": fileModel.UUID}, http.StatusOK)
	}
}
