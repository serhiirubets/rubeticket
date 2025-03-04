package middleware

import (
	"context"
	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/pkg/jwt"
	"net/http"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
	ContextIdKey    key = "ContextIdKey"
)

type AuthContextData struct {
	Email  string
	UserID uint
}

const AuthKey = "authData"

func writeUnathed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func Auth(next http.Handler, conf *config.Config) http.Handler {
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
		isValid, data := jwt.NewJWT(conf.Auth.Secret).Parse(token)

		if !isValid {
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
