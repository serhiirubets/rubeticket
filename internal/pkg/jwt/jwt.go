package jwt

import "github.com/golang-jwt/jwt/v5"

type Payload struct {
	Email string
	Id    uint
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{Secret: secret}
}

func (j *JWT) Create(data *Payload) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"email": data.Email, "id": data.Id})
	s, err := t.SignedString([]byte(j.Secret))

	if err != nil {
		return "", err
	}

	return s, nil
}

func (j *JWT) Parse(token string) (bool, *Payload) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		println(err.Error())
		return false, nil
	}

	email := t.Claims.(jwt.MapClaims)["email"].(string)
	id := t.Claims.(jwt.MapClaims)["id"].(float64)
	return t.Valid, &Payload{Email: email, Id: uint(id)}
}
