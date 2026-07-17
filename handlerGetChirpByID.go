package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (configPtr *apiConfig) handlerGetChirpByID(response http.ResponseWriter, request *http.Request) {
	chirpIDString := request.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondwithError(response, 400, "Unable to parse ID or bad/invalid ID format", err)
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

	respondWithJSON(response, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
