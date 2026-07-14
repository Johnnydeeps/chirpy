package main

import "net/http"

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
	serverMux.HandleFunc("/healthz", healthz)

	// Serve files from current directory.
	fileServer := http.FileServer(http.Dir("."))
	// Remove "/app" from URL before file lookup.
	appHandler := http.StripPrefix("/app", fileServer)
	// Route /app/... requests to the file server.
	serverMux.Handle("/app/", appHandler)

	// Server listens on port 8080 and uses the mux.
	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	// Start listening for HTTP requests.
	server.ListenAndServe()
}
