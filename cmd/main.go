package main

import (
	"fmt"
	"net/http"
	"self-service-tooling/api" // Adjust the import path based on your project structure
)

func main() {
	http.HandleFunc("/get-terraform-code", getTerraformCodeHandler)
	http.HandleFunc("/terraform-apply", api.DeployHandler)
	http.HandleFunc("/terraform-plan", api.RunTerraformPlan)        // Call directly if functions are exported
	http.HandleFunc("/terraform-destroy", api.DestroyTerraformCode) // Call directly

	// Serve static files (e.g., index.html, main.js, styles.css)
	fs := http.FileServer(http.Dir("./ui"))
	http.Handle("/", fs)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func getTerraformCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Fetching Terraform code...")
}

func DeployHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Deploying...")
}
