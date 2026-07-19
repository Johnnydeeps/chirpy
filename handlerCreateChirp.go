package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Johnnydeeps/chirpy/internal/auth"
	"github.com/Johnnydeeps/chirpy/internal/database"
	"github.com/google/uuid"
)

type parametersChirp struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (configPtr *apiConfig) handlerCreateChirp(response http.ResponseWriter, request *http.Request) {
	paramsChirp := parametersChirp{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&paramsChirp)
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
		respondwithError(response, 401, "Unauthorized", err)
		return
	}

	chirp, err := configPtr.databasePtr.CreateChirp(request.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      paramsChirp.Body,
		UserID:    userID,
	})
	if err != nil {
		respondwithError(response, 500, "Error creating Chirp", err)
		return
	}

	respondWithJSON(response, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
