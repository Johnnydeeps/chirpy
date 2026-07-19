package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no value provided")
	}
	trimmed := strings.TrimPrefix(auth, "Bearer")
	return strings.TrimSpace(trimmed), nil
}
