package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Johnnydeeps/chirpy/internal/auth"
	"github.com/Johnnydeeps/chirpy/internal/database"
)

func (configPtr *apiConfig) handlerRevokeRefreshToken(response http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondwithError(response, 401, "Unauthorized", err)
		return
	}

	err = configPtr.databasePtr.RevokeRefreshToken(request.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		Token:     token,
	})
	if err != nil {
		respondwithError(response, 500, "DB error", err)
		return
	}

	// operation was sucessful but no content message http response: code 204.
	response.WriteHeader(204)
}
