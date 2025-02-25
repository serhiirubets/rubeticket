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
}
