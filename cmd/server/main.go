package main

import (
	"fmt"
	"log"
	"net/http"
	"self-service-tooling/api"
)

func main() {

	// Serve static files (HTML, CSS, JS) from ./ui
	http.Handle("/", http.FileServer(http.Dir("./ui")))

	// OAuth routes (from api/oauth.go)
	http.HandleFunc("/login", api.LoginHandler)
	http.HandleFunc("/auth/github/callback", api.CallbackHandler)

	// Additional API routes
	http.HandleFunc("/api/repositories", api.GetRepositoriesHandler)
	http.HandleFunc("/api/terraform/plan", api.TerraformPlanHandler)
	http.HandleFunc("/api/terraform/apply", api.TerraformApplyHandler)

	// Start the server
	port := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
