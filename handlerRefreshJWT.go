package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Johnnydeeps/chirpy/internal/auth"
)

type responseParametersRefreshJSON struct {
	Token string `json:"token"`
}

// The short-lived JWT access token expired after 60 minuets, but the long-lived refresh token
// exists for 60 days — Based on a confirmed and existing long-lived refresh token, refresh or
// grant a new shortlived access token.
func (configPtr *apiConfig) handlerRefreshJWT(response http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondwithError(response, 401, "Unauthorized", err)
		return
	}
	user, err := configPtr.databasePtr.GetUserFromRefreshToken(request.Context(), token)
	if errors.Is(err, sql.ErrNoRows) {
		respondwithError(response, 401, "Failed to retrieve user", err)
		return
	} else if err != nil {
		respondwithError(response, 401, "No matching user token", err)
		return
	}
	if user.ExpiresAt.Before(time.Now()) || user.RevokedAt.Valid == true {
		respondwithError(response, 401, "token expired", nil)
		return
	}

	const maxExpiration = time.Hour
	signedToken, err := auth.MakeJWT(user.ID, configPtr.secretKey, maxExpiration)
	if err != nil {
		respondwithError(response, 500, "Token Initialization Failure", err)
		return
	}

	respondWithJSON(response, 200, responseParametersRefreshJSON{
		Token: signedToken,
	})
}
