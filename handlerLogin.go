package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Johnnydeeps/chirpy/internal/auth"
)

type parametersLoginJSON struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (configPtr *apiConfig) handlerLogin(response http.ResponseWriter, request *http.Request) {
	params := parametersLoginJSON{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondwithError(response, 400, "Bad Resquest: Error decoding JSON", err)
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

	const maxExpiration = time.Hour
	expiration := maxExpiration
	if params.ExpiresInSeconds > 0 && time.Duration(params.ExpiresInSeconds)*time.Second < maxExpiration {
		expiration = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	signedToken, err := auth.MakeJWT(user.ID, configPtr.secretKey, expiration)
	if err != nil {
		respondwithError(response, 500, "Token Initialization Failure", err)
		return
	}

	respondWithJSON(response, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     signedToken,
	})

}
