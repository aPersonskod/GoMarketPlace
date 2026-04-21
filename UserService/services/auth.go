package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("my_important_secret")

type Claims struct {
	Id   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(id, role string) (string, error) {
	claims := Claims{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString(JwtSecret)
	return stringToken, err
}
