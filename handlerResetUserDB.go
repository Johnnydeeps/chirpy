package main

import "net/http"

func (configPtr *apiConfig) reset(response http.ResponseWriter, request *http.Request) {
	if configPtr.platform != "dev" {
		respondwithError(response, 403, "Access Denied: do not have dev access", nil)
		return
	}
	configPtr.fileserverHits.Store(0)

	err := configPtr.databasePtr.ResetAllUsers(request.Context())
	if err != nil {
		respondwithError(response, 500, "Error Reseting DB", err)
		return
	}

	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Count reset to 0"))
}
