package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

var apiCfg = &apiConfig{}

func (configPtr *apiConfig) middlewareMetricsINC(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		configPtr.fileserverHits.Add(1)
		next.ServeHTTP(response, request)
	})
}

func (configPtr *apiConfig) metrics(response http.ResponseWriter, request *http.Request) {
	count := configPtr.fileserverHits.Load()
	response.Header().Set("Content-Type", "text/html")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", count))) // ignore
}

func (configPtr *apiConfig) reset(response http.ResponseWriter, request *http.Request) {
	configPtr.fileserverHits.Store(0)
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Count reset to 0"))
}

// custom handler function (not object like for fileserver) that writes an http response as
// a function to by called by a handler below
func healthz(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("OK"))
}

func main() {
	// Route requests by path.
	serverMux := http.NewServeMux()

	//call the custom handler function if a http request is made to /healthz
	serverMux.HandleFunc("GET /api/healthz", healthz)

	//call the config method functions if an http request is made to /metrics or /reset
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.metrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.reset)
	// Serve files from current directory.
	fileServer := http.FileServer(http.Dir("."))
	// Remove "/app" from URL before file lookup.
	appHandler := http.StripPrefix("/app", fileServer)
	// Route /app/... requests to the file server.
	serverMux.Handle("/app/", apiCfg.middlewareMetricsINC(appHandler))

	// Server listens on port 8080 and uses the mux.
	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	// Start listening for HTTP requests.
	server.ListenAndServe()
}
