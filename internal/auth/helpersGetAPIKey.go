package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("no API Key provided")
	}
	trimmed := strings.TrimPrefix(apiKey, "ApiKey")
	return strings.TrimSpace(trimmed), nil
}
