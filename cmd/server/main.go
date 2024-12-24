// File: cmd/server/main.go

package main

import (
	"log"
	"net/http"
	"os"

	"self-service-tooling/api"
)

func main() {
	// OAuth routes
	http.HandleFunc("/auth/github/login", api.LoginHandler)
	http.HandleFunc("/auth/github/callback", api.CallbackHandler)

	// API routes
	http.HandleFunc("/api/repos", api.GetReposHandler)

	// Form submission route
	http.HandleFunc("/configure", api.ConfigureHandler)

	// Debug endpoint (remove in production)
	http.HandleFunc("/debug/session", api.DebugHandler)

	// Serve static files from the 'ui' directory
	// Ensure this is registered after API routes to prevent overriding
	fs := http.FileServer(http.Dir("./ui"))
	http.Handle("/", fs)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
