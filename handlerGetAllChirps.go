package main

import "net/http"

func (configPtr *apiConfig) handlerGetAllChirps(response http.ResponseWriter, request *http.Request) {
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
	respondWithJSON(response, 200, chirpsJSON)
}
