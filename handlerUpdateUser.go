package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Johnnydeeps/chirpy/internal/auth"
	"github.com/Johnnydeeps/chirpy/internal/database"
	"github.com/google/uuid"
)

type parametersLUpdateUserJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type responseParametersUpdateUserJSON struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}

func (configPtr *apiConfig) handlerUpdateUser(response http.ResponseWriter, request *http.Request) {
	params := parametersLUpdateUserJSON{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondwithError(response, 400, "Bad Resquest: Error decoding JSON", err)
		return
	}

	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondwithError(response, 401, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, configPtr.secretKey)
	if err != nil {
		respondwithError(response, 401, "Token validation error", err)
		return
	}

	newHashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondwithError(response, 500, "DB error hashing new password", err)
		return
	}

	updatedUser, err := configPtr.databasePtr.UpdateUserHashedPasswordOrEmail(request.Context(),
		database.UpdateUserHashedPasswordOrEmailParams{
			HashedPassword: newHashedPassword,
			Email:          params.Email,
			UpdatedAt:      time.Now(),
			ID:             userID,
		})
	if errors.Is(err, sql.ErrNoRows) {
		respondwithError(response, 401, "Failed to update user", err)
		return
	} else if err != nil {
		respondwithError(response, 500, "No matching user token", err)
		return
	}

	respondWithJSON(response, 200, responseParametersUpdateUserJSON{
		ID:             userID,
		CreatedAt:      updatedUser.CreatedAt,
		UpdatedAt:      updatedUser.UpdatedAt,
		Email:          updatedUser.Email,
		HashedPassword: newHashedPassword,
		IsChirpyRed:    updatedUser.IsChirpyRed.Bool,
	})

}
