package accounts

import (
	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/users"
	"github.com/serhiirubets/rubeticket/pkg/log"
	"github.com/serhiirubets/rubeticket/pkg/middleware"
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
}

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
