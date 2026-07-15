package main

import (
	"strings"
)

func getCleanedBody(body string, badwords map[string]struct{}) string {
	splitBody := strings.Split(body, " ")
	for i, word := range splitBody {
		_, exists := badwords[strings.ToLower(word)]
		if exists {
			splitBody[i] = "****"
		}
	}
	return strings.Join(splitBody, " ")
}
