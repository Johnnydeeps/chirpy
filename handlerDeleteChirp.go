package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Johnnydeeps/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (configPtr *apiConfig) handlerDeleteChirp(response http.ResponseWriter, request *http.Request) {
	chirpIDString := request.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondwithError(response, 400, "Unable to parse ID or bad/invalid ID format", err)
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

	chirp, err := configPtr.databasePtr.GetChirpByID(request.Context(), chirpID)
	if errors.Is(err, sql.ErrNoRows) {
		respondwithError(response, 404, fmt.Sprintf("Failed to retrieve chirp at ID; %v", chirpIDString), err)
		return
	} else if err != nil {
		respondwithError(response, 500, "DB Error", err)
		return
	}

	if chirp.UserID != userID {
		respondwithError(response, 403, "Unauthorized", nil)
		return
	}

	err = configPtr.databasePtr.DeleteChirp(request.Context(), chirp.ID)
	if err != nil {
		respondwithError(response, 500, "DB error deleting chirp", err)
		return
	}

	response.WriteHeader(204)
}
