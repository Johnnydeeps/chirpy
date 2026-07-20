package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Johnnydeeps/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	databasePtr    *database.Queries
	platform       string
	secretKey      string
	polkaKey       string
}

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
	response.Write(fmt.Appendf(nil, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", count)) // ignore
}

// custom handler function (not object like for fileserver) that writes an http response as
// a function to by called by a handler below
func healthz(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("OK"))
}

func main() {
	// load and assign .env values to be stored in api config struct for the module.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretKey := os.Getenv("JWT_SECRET_KEY")
	polkakey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	//***************************************************************************************

	// config struct pointer for module reference.
	apiCfg := &apiConfig{
		databasePtr: dbQueries,
		platform:    platform,
		secretKey:   secretKey,
		polkaKey:    polkakey,
	}
	//***************************************************************************************

	// Route requests by path.
	serverMux := http.NewServeMux()

	//call the custom handler function if a http request is made to /healthz
	serverMux.HandleFunc("GET /api/healthz", healthz)

	//call the config method functions if an http request is made to /metrics or /reset etc.
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.metrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.reset)
	serverMux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	serverMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serverMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	serverMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshJWT)
	serverMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)
	serverMux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	serverMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	serverMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerWebhooks)

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
