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

	router.HandleFunc("GET /users/{id}", handler.GetById())
}

// GetById godoc
// @Summary Get user by ID
// @Description Returns a user by their ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} GetUserResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/{id} [get]
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
