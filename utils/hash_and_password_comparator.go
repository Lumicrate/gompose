package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("invalid username or password")
	}
	return nil
}
