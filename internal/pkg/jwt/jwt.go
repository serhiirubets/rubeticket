package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/serhiirubets/rubeticket/internal/app/users"
)

type Payload struct {
	Email string
	Id    uint
	Role  users.Role
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{Secret: secret}
}

func (j *JWT) Create(data *Payload) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		"id":    data.Id,
		"role":  data.Role,
	})
	s, err := t.SignedString([]byte(j.Secret))

	if err != nil {
		return "", err
	}

	return s, nil
}

func (j *JWT) Parse(token string) (*Payload, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		println(err.Error())
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !t.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type, expected jwt.MapClaims")
	}

	email, ok := claims["email"]
	if !ok {
		return nil, fmt.Errorf("email field is missing")
	}

	emailStr, ok := email.(string)
	if !ok {
		return nil, fmt.Errorf("email field must be a string")
	}

	id, ok := claims["id"]
	if !ok {
		return nil, fmt.Errorf("id field is missing")
	}
	idFloat, ok := id.(float64)
	if !ok {
		return nil, fmt.Errorf("id field must be a number")
	}

	role, ok := claims["role"]
	if !ok {
		return nil, fmt.Errorf("role field is missing")
	}
	roleStr, ok := role.(string)
	if !ok {
		return nil, fmt.Errorf("role field must be a string")
	}

	return &Payload{
		Email: emailStr,
		Id:    uint(idFloat),
		Role:  users.Role(roleStr),
	}, nil
}
