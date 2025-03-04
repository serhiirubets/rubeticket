package auth

import (
	"github.com/serhiirubets/rubeticket/internal/app/users"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Success bool `json:"success"`
}

type RegisterRequest struct {
	Email     string       `json:"email" validate:"required,email"`
	FirstName string       `json:"firstName" validate:"required,alphaunicode"`
	LastName  string       `json:"lastName" validate:"required,alphaunicode"`
	Password  string       `json:"password" validate:"required,min=6"`
	Gender    users.Gender `json:"gender" validate:"required,oneof=male female"`
	Birthday  string       `json:"birthday" validate:"required,date_iso8601"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}
