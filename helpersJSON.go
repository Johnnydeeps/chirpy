package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(response http.ResponseWriter, code int, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		response.WriteHeader(500)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(code)
	response.Write(body)
}

func respondwithError(response http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(response, code, errorResponse{
		Error: msg,
	})
}
