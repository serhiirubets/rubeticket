package auth

import (
	configs "github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/pkg/jwt"
	"github.com/serhiirubets/rubeticket/pkg/log"
	"github.com/serhiirubets/rubeticket/pkg/req"
	"github.com/serhiirubets/rubeticket/pkg/res"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
	Logger log.ILogger
}

type AuthHandler struct {
	*configs.Config
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

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)

		if err != nil {
			return
		}

		email, err := handler.AuthService.Login(body.Email, body.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(&jwt.JWTData{Email: email})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := LoginResponse{
			Token: token,
		}

		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)

		if err != nil {
			handler.Logger.Error("Registration: request body parsing failed")
			return
		}

		_, registerErr := handler.AuthService.Register(body)

		if registerErr != nil {
			handler.Logger.WithFields(log.WithFields{
				"user_email": body.Email,
			}).Error("Registration failed")

			http.Error(w, registerErr.Error(), http.StatusUnauthorized)
			return
		}

		token, jwtErr := jwt.NewJWT(handler.Config.Auth.Secret).Create(&jwt.JWTData{Email: body.Email})

		if jwtErr != nil {
			handler.Logger.WithFields(log.WithFields{
				"user_email": body.Email,
			}).Error("Creating jwt failed")

			http.Error(w, jwtErr.Error(), http.StatusInternalServerError)
			return
		}

		data := RegisterResponse{
			Token: token,
		}

		res.Json(w, data, http.StatusOK)
	}
}
