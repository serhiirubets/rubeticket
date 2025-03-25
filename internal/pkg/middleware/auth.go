package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/app/users"
	"github.com/serhiirubets/rubeticket/internal/pkg/jwt"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/res"
)

type contextKey string

type AuthContextData struct {
	Email  string
	UserID uint
	Role   users.Role
}

const AuthKey contextKey = "authData"

type AuthMiddleware struct {
	conf       *config.Config
	logger     log.ILogger
	openRoutes map[string]struct{}
	apiPrefix  string // Example: "/api/v1"
}

func NewAuthMiddleware(conf *config.Config, logger log.ILogger, openRoutes []string, apiPrefix string) *AuthMiddleware {
	openRoutesMap := make(map[string]struct{})
	for _, route := range openRoutes {
		normalizedRoute := "/" + strings.Trim(route, "/")
		openRoutesMap[normalizedRoute] = struct{}{}
	}
	return &AuthMiddleware{
		conf:       conf,
		logger:     logger,
		openRoutes: openRoutesMap,
		apiPrefix:  apiPrefix,
	}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, является ли маршрут открытым
		path := r.URL.Path
		normalizedPath := "/" + strings.Trim(path, "/")

		isOpen := false
		for route := range m.openRoutes {
			if strings.Contains(route, "{") && strings.Contains(route, "}") {
				// Check for public routes with params
				parts := strings.Split(route, "{")
				prefix := parts[0]
				if strings.HasPrefix(normalizedPath, prefix) {
					isOpen = true
					break
				}
			} else if normalizedPath == route {
				// Exact match
				isOpen = true
				break
			}
		}

		if isOpen {
			// Skip token check for open routes
			next.ServeHTTP(w, r)
			return
		}

		// Check token for closed routes
		cookie, err := r.Cookie("token")
		if err != nil {
			m.logger.Debug("No token cookie found", "error", err.Error())
			writeUnathed(w)
			return
		}

		token := cookie.Value
		if token == "" {
			m.logger.Debug("Token is empty")
			writeUnathed(w)
			return
		}

		data, parseErr := jwt.NewJWT(m.conf.Auth.Secret).Parse(token)
		if parseErr != nil {
			m.logger.Error("Token parse failed", "error", parseErr.Error())
			writeUnathed(w)
			return
		}

		authData := AuthContextData{
			Email:  data.Email,
			UserID: data.Id,
			Role:   users.UserRole,
		}

		ctx := context.WithValue(r.Context(), AuthKey, authData)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func (m *AuthMiddleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authData, err := GetAuthData(r)
		if err != nil {
			m.logger.Error("Admin auth failed", "error", err.Error())
			res.Json(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if authData.Role != users.AdminRole {
			m.logger.Error("Admin access required", "role", authData.Role)
			res.Json(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeUnathed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func GetAuthData(r *http.Request) (AuthContextData, error) {
	authData, ok := r.Context().Value(AuthKey).(AuthContextData)
	if !ok {
		return AuthContextData{}, fmt.Errorf("not authorized")
	}
	return authData, nil
}
