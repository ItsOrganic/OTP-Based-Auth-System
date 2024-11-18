package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(phoneNumber string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["phone_number"] = phoneNumber
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	tokenString, err := token.SignedString([]byte("your-secret-key")) // Replace with secure key from env
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
