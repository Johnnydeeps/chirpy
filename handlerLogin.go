package main

import (
	"encoding/json"
	"net/http"

	"github.com/Johnnydeeps/chirpy/internal/auth"
)

func (configPtr *apiConfig) handlerLogin(response http.ResponseWriter, request *http.Request) {
	params := parametersUserJSON{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondwithError(response, 500, "Error decoding JSON", err)
		return
	}

	user, err := configPtr.databasePtr.GetUserByEmailLogin(request.Context(), params.Email)
	if err != nil {
		respondwithError(response, 401, "Unauthorized: Incorrect email or password", err)
		return
	}

	isValid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondwithError(response, 401, "Unauthorized: Incorrect email or password", err)
		return
	}
	if isValid != true {
		respondwithError(response, 401, "Unauthorized: Incorrect email or password", nil)
		return
	}

	respondWithJSON(response, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
