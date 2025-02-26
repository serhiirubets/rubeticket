package accounts

import (
	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/users"
	"github.com/serhiirubets/rubeticket/pkg/log"
	"github.com/serhiirubets/rubeticket/pkg/middleware"
	"github.com/serhiirubets/rubeticket/pkg/req"
	"github.com/serhiirubets/rubeticket/pkg/res"
	"net/http"
)

type AccountHandlerDeps struct {
	UserRepository *users.UserRepository
	Logger         log.ILogger
	Config         *config.Config
}

type AccountHandler struct {
	UserRepository *users.UserRepository
	Logger         log.ILogger
	Config         *config.Config
}

func NewAccountHandler(router *http.ServeMux, deps *AccountHandlerDeps) {
	handler := &AccountHandler{
		UserRepository: deps.UserRepository,
		Logger:         deps.Logger,
		Config:         deps.Config,
	}

	router.Handle("GET /account", middleware.Auth(handler.GetAccount(), deps.Config))
	router.Handle("PATCH /account", middleware.Auth(handler.UpdateAccountPatch(), deps.Config))
	router.Handle("PUT /account", middleware.Auth(handler.UpdateAccountPut(), deps.Config))
}

// GetAccount godoc
// @Summary Get account info
// @Description Return information about current user
// @Tags Account
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} GetAccountResponse "Success"
// @Failure 401 {object} string "Not authorized"
// @Router /account [get]
func (handler *AccountHandler) GetAccount() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, _ := r.Context().Value(middleware.ContextEmailKey).(string)

		user, userErr := handler.UserRepository.GetByEmail(email)

		if userErr != nil {
			handler.Logger.Error("Error getting user by email", userErr.Error())
			res.Json(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		body := GetAccountResponse{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Gender:    user.Gender,
			Birthday:  user.Birthday,
		}

		res.Json(w, body, http.StatusOK)
	})
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
// @Router /account [patch]
func (handler *AccountHandler) UpdateAccountPatch() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем email из контекста (добавлен middleware.Auth)
		email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			handler.Logger.Error("Email not found in context")
			res.Json(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		body, err := req.HandleBody[UpdateAccountRequestPatch](&w, r)

		if err != nil {
			return
		}

		user, err := handler.UserRepository.GetByEmail(email)
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

		if err := handler.UserRepository.DB.Model(&user).Updates(updates).Error; err != nil {
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
	})
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
// @Router /account [put]
func (handler *AccountHandler) UpdateAccountPut() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			handler.Logger.Error("Email not found in context")
			res.Json(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		body, errBody := req.HandleBody[UpdateAccountRequestPut](&w, r)

		if errBody != nil {
			return
		}

		user, errUser := handler.UserRepository.GetByEmail(email)
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
	})
}
