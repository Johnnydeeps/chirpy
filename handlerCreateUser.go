package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Johnnydeeps/chirpy/internal/auth"
	"github.com/Johnnydeeps/chirpy/internal/database"
	"github.com/google/uuid"
)

type parametersUserJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	// `json:"-"` the dash json tag flags this feild to be skipped when encoding a json response.
	HashedPassword string `json:"-"`
	Token          string `json:"token"`
	RefreshToken   string `json:"refresh_token"`
}

func (configPtr *apiConfig) handlerCreateUser(response http.ResponseWriter, request *http.Request) {
	params := parametersUserJSON{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondwithError(response, 500, "Error decoding JSON", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondwithError(response, 500, "Error hashing password", err)
		return
	}

	user, err := configPtr.databasePtr.CreateUser(request.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          params.Email,
		HashedPassword: hashedPassword,
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
