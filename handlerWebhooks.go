package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Johnnydeeps/chirpy/internal/auth"
	"github.com/google/uuid"
)

type parametersUpgradeUserEvent struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (configPtr *apiConfig) handlerWebhooks(response http.ResponseWriter, request *http.Request) {
	params := parametersUpgradeUserEvent{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondwithError(response, 400, "Bad Resquest: Error decoding JSON", err)
		return
	}

	//verify polka apikey from webhook request
	apiKey, err := auth.GetAPIKey(request.Header)
	if err != nil {
		respondwithError(response, 401, "Apikey missing or malformed in request header", err)
		return
	}
	if apiKey != configPtr.polkaKey {
		response.WriteHeader(401)
	}

	if params.Event != "user.upgraded" {
		response.WriteHeader(204)
		return
	}

	_, err = configPtr.databasePtr.UgradeUserRedChirp(request.Context(), params.Data.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		respondwithError(response, 404, "No matching user id", err)
		return
	} else if err != nil {
		respondwithError(response, 500, "DB Error updating user table column", err)
		return
	}

	response.WriteHeader(204)
}
