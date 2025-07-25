package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(userID, secretKey string, exp time.Duration) (string, error) {
	if exp < 1 {
		exp = time.Hour * 24
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(exp).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}
