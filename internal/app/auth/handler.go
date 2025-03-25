package auth

import (
	"net/http"

	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/app/users"
	"github.com/serhiirubets/rubeticket/internal/pkg/jwt"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/req"
	"github.com/serhiirubets/rubeticket/internal/pkg/res"
)

type AuthHandlerDeps struct {
	*config.Config
	*AuthService
	Logger log.ILogger
}

type AuthHandler struct {
	*config.Config
	*AuthService
	Logger log.ILogger
}

func NewAuthHandler(router *http.ServeMux, deps *AuthHandlerDeps) {
	handler := AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
		Logger:      deps.Logger,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

// Login godoc
// @Summary Login a user
// @Description Login user, set token in cookie and return success: true
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "LoginRequest credentials"
// @Success 200 {object} LoginResponse "Successfully logged in"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/auth/login [post]
func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)

		if err != nil {
			return
		}

		loginDto, err := handler.AuthService.Login(body.Email, body.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(&jwt.Payload{Email: loginDto.Email, Id: loginDto.Id, Role: loginDto.Role})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.SetToken(w, token)
		res.Json(w, &LoginResponse{Success: true}, http.StatusOK)
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user and return authentication token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration credentials"
// @Success 200 {object} RegisterResponse "Successfully registered"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/auth/register [post]
func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)

		if err != nil {
			handler.Logger.Error("Registration: request body parsing failed")
			return
		}

		id, registerErr := handler.AuthService.Register(body)

		if registerErr != nil {
			handler.Logger.WithFields(log.WithFields{
				"user_email": body.Email,
			}).Error("Registration failed")

			http.Error(w, registerErr.Error(), http.StatusUnauthorized)
			return
		}

		token, jwtErr := jwt.NewJWT(handler.Config.Auth.Secret).Create(&jwt.Payload{Email: body.Email, Id: id, Role: users.UserRole})

		if jwtErr != nil {
			handler.Logger.WithFields(log.WithFields{
				"user_email": body.Email,
			}).Error("Creating jwt failed")

			http.Error(w, jwtErr.Error(), http.StatusInternalServerError)
			return
		}

		res.SetToken(w, token)
		res.Json(w, &RegisterResponse{Success: true}, http.StatusOK)
	}
}
