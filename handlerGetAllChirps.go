package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (configPtr *apiConfig) handlerGetAllChirps(response http.ResponseWriter, request *http.Request) {
	authorIDString := request.URL.Query().Get("author_id")
	sortIDString := request.URL.Query().Get("sort")

	if authorIDString == "" {

		chirps, err := configPtr.databasePtr.GetAllChirps(request.Context())
		if err != nil {
			respondwithError(response, 500, "Error retrieving chirps from DB.", err)
			return
		}
		chirpsJSON := []Chirp{}

		for _, chirp := range chirps {
			chirpsJSON = append(chirpsJSON, Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}
		// sort to DESC/descending order using slice in memory rather than writing a new SQL query.
		// db retrieval queires (ie. GetAllChirps or GetChirpsByUser) return ASC/ascending sorted
		// result by default.
		if sortIDString == "desc" {
			sort.Slice(chirpsJSON,
				func(i, j int) bool { return chirpsJSON[i].CreatedAt.After(chirpsJSON[j].CreatedAt) })
			respondWithJSON(response, 200, chirpsJSON)
			return
		}
		respondWithJSON(response, 200, chirpsJSON)
		return
	}

	requestUserID, err := uuid.Parse(authorIDString)
	if err != nil {
		respondwithError(response, 400, "Bad Request: Error parsing user id or malformed string", err)
		return
	}

	chirps, err := configPtr.databasePtr.GetChirpsByUser(request.Context(), requestUserID)
	if errors.Is(err, sql.ErrNoRows) {
		respondwithError(response, 404, "No chirps matching user id", err)
		return
	} else if err != nil {
		respondwithError(response, 500, "DB error", err)
		return
	}

	chirpsJSON := []Chirp{}
	for _, chirp := range chirps {
		chirpsJSON = append(chirpsJSON, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	// sort to DESC/descending order using slice in memory rather than writing a new SQL query.
	// db retrieval queires (ie. GetAllChirps or GetChirpsByUser) return ASC/ascending sorted
	// result by default.
	if sortIDString == "desc" {
		sort.Slice(chirpsJSON,
			func(i, j int) bool { return chirpsJSON[i].CreatedAt.After(chirpsJSON[j].CreatedAt) })
		respondWithJSON(response, 200, chirpsJSON)
		return
	}
	respondWithJSON(response, 200, chirpsJSON)

}
