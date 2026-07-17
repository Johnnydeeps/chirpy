package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Johnnydeeps/chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (configPtr *apiConfig) handlerCreateUser(response http.ResponseWriter, request *http.Request) {
	params := parameters{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondwithError(response, 500, "Error decoding JSON", err)
		return
	}

	user, err := configPtr.databasePtr.CreateUser(request.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     params.Email,
	})
	if err != nil {
		respondwithError(response, 500, "Error creating user in database", err)
		return
	}

	respondWithJSON(response, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
