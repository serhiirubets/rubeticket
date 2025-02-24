package users

import (
	"github.com/serhiirubets/rubeticket/pkg/log"
	"net/http"
)

type UserHandlerDeps struct {
	UserRepository *UserRepository
	Logger         log.ILogger
}

type UserHandler struct {
	UserRepository *UserRepository
	Logger         log.ILogger
}

func NewUserHandler(router *http.ServeMux, deps *UserHandlerDeps) {
	handler := &UserHandler{
		UserRepository: deps.UserRepository,
		Logger:         deps.Logger,
	}

	router.Handle("Get /users", handler.Find())
}

func (handler *UserHandler) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
