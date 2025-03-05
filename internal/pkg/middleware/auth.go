package middleware

import (
	"context"
	"fmt"
	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/pkg/jwt"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/res"
	"net/http"
)

type key string

type AuthContextData struct {
	Email  string
	UserID uint
}

const AuthKey = "authData"

func writeUnathed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuth(next http.Handler, conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			writeUnathed(w)
			return
		}

		token := cookie.Value
		if token == "" {
			writeUnathed(w)
			return
		}
		data, parseErr := jwt.NewJWT(conf.Auth.Secret).Parse(token)

		if parseErr != nil {
			// TODO: add logger
			writeUnathed(w)
			return
		}

		authData := AuthContextData{
			Email:  data.Email,
			UserID: data.Id,
		}

		ctx := context.WithValue(r.Context(), AuthKey, authData)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func GetAuthData(r *http.Request) (AuthContextData, error) {
	authData, ok := r.Context().Value(AuthKey).(AuthContextData)
	if !ok {
		return AuthContextData{}, fmt.Errorf("not authorized")
	}
	return authData, nil
}

func Auth(next http.HandlerFunc, config *config.Config, logger log.ILogger) http.Handler {
	return IsAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authData, err := GetAuthData(r)
		if err != nil {
			logger.Error("Auth failed", "error", err.Error())
			res.Json(w, err.Error(), http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), AuthKey, authData))
		next.ServeHTTP(w, r)
	}), config)
}
