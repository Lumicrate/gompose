package utils

import (
	"errors"
	"strings"
)

func ExtractBearerToken(header string) (string, error) {
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return "", errors.New("missing or invalid Authorization header")
	}
	return strings.TrimPrefix(header, "Bearer "), nil
}
