package accounts

import (
	"github.com/serhiirubets/rubeticket/internal/users"
	"time"
)

type GetAccountResponse struct {
	Email     string       `json:"email"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	Gender    users.Gender `json:"gender"`
	Birthday  time.Time    `json:"birthday"`
	PhotoUrl  string       `json:"photoUrl"`
	Address   string       `json:"address"`
}

type UpdateAccountRequestPatch struct {
	FirstName *string       `json:"firstName,omitempty"`
	LastName  *string       `json:"lastName,omitempty"`
	Gender    *users.Gender `json:"gender,omitempty"`
	Birthday  *time.Time    `json:"birthday,omitempty"`
	PhotoUrl  *string       `json:"photoUrl,omitempty"`
	Address   *string       `json:"address,omitempty"`
}

type UpdateAccountRequestPut struct {
	FirstName string       `json:"firstName" validate:"required,alphaunicode"`
	LastName  string       `json:"lastName" validate:"required,alphaunicode"`
	Gender    users.Gender `json:"gender" validate:"required,oneof=male female"`
	Birthday  time.Time    `json:"birthday" validate:"required"`
	PhotoUrl  string       `json:"photoUrl"`
	Address   string       `json:"address"`
}

type UpdateAccountResponse struct {
	Email     string       `json:"email"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	Gender    users.Gender `json:"gender"`
	PhotoUrl  string       `json:"photoUrl"`
	Address   string       `json:"address"`
	Birthday  time.Time    `json:"birthday"`
}
