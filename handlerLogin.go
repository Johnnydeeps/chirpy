package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Johnnydeeps/chirpy/internal/auth"
	"github.com/Johnnydeeps/chirpy/internal/database"
)

type parametersLoginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	signedToken, err := auth.MakeJWT(user.ID, configPtr.secretKey, maxExpiration)
	if err != nil {
		respondwithError(response, 500, "Token Initialization Failure", err)
		return
	}
	// create resfresh token expiration time of 60 days as time.Duration
	const refreshTokenExpiration = time.Hour * 24 * 60

	refreshToken := auth.MakeRefreshToken()

	storedRefreshToken, err := configPtr.databasePtr.CreateRefresgToken(request.Context(), database.CreateRefresgTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(refreshTokenExpiration),
		RevokedAt: sql.NullTime{},
	})

	respondWithJSON(response, 200, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        signedToken,
		RefreshToken: storedRefreshToken.Token,
	})

}
