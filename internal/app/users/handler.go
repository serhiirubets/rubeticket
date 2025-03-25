package users

import (
	"net/http"

	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/res"
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

	router.HandleFunc("GET /users", handler.Find())
	router.HandleFunc("GET /users/{id}", handler.GetById())
}

func (handler *UserHandler) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement find users
	}
}

func (handler *UserHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		user, userErr := handler.UserRepository.GetById(id)

		if userErr != nil {
			res.Json(w, nil, http.StatusNotFound)
			return
		}

		res.Json(w, user, http.StatusOK)
	}
}
