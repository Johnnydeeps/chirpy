package main

import (
	"encoding/json"
	"net/http"
)

func handlerChirpsValidate(response http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValid struct {
		CleanedBody string `json:"cleaned_body"`
	}

	params := parameters{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondwithError(response, 500, "Server Error Decoding JSON", err)
		return
	}

	if len(params.Body) > 140 {
		respondwithError(response, 400, "Chirp is too long", nil)
		return
	}
	badwords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	respondWithJSON(response, 200, returnValid{
		CleanedBody: getCleanedBody(params.Body, badwords),
	})
}
