package auth

import "github.com/serhiirubets/rubeticket/internal/app/users"

type LoginResponseDto struct {
	Id    uint
	Email string
	Role  users.Role
}
