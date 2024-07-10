package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAuthToken(email string,id int) (*string, error) {
	claims := jwt.MapClaims{
		"id": id,
		"email": email,
		"user":  true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("modi_sarkar"))

	if err != nil {
		return nil, err
	}

	return &t, nil
}
