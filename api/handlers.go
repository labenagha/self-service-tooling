package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"self-service-tooling/config" // Replace with your actual module path
)

// Struct to parse JSON payload from the frontend
type TerraformRequest struct {
	Repo      string `json:"repo"`
	Directory string `json:"directory"`
}

// GetTerraformCode handles fetching Terraform code from the selected repository and directory.
func GetTerraformCode(w http.ResponseWriter, r *http.Request) {
	// Get session
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Session values in GetTerraformCode:", session.Values)

	// Check if the user is authenticated
	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		log.Println("GitHub OAuth not completed for user:", session.Values["username"])
		http.Error(w, "GitHub OAuth not completed", http.StatusUnauthorized)
		return
	}

	// Get the GitHub access token from the session
	token, ok := session.Values["github_access_token"].(string)
	if !ok || token == "" {
		http.Error(w, "GitHub token is not set", http.StatusInternalServerError)
		return
	}

	repo := r.URL.Query().Get("repo")
	dirPath := r.URL.Query().Get("dir")

	if repo == "" || dirPath == "" {
		http.Error(w, "Repository or directory not specified", http.StatusBadRequest)
		return
	}

	// Call the existing function from utils.go
	err = DownloadRepositoryContent(repo, dirPath, token)
	if err != nil {
		http.Error(w, "Error fetching Terraform code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Terraform code fetched successfully from repository:", repo)
}

// DeployTerraformCode handles the deployment request.
func DeployTerraformCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req TerraformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Run Terraform apply in the specified directory
	output, err := RunTerraformCommand("terraform apply -auto-approve", req.Directory)
	if err != nil {
		http.Error(w, "Error running Terraform apply: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, output)
}

// RunTerraformPlan handles running `terraform plan`.
func RunTerraformPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TerraformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	output, err := RunTerraformCommand("terraform plan", req.Directory)
	if err != nil {
		http.Error(w, "Error running Terraform plan: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, output)
}

// DestroyTerraformCode handles running `terraform destroy`.
func DestroyTerraformCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req TerraformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Run Terraform destroy in the specified directory
	output, err := RunTerraformCommand("terraform destroy -auto-approve", req.Directory)
	if err != nil {
		http.Error(w, "Error running Terraform destroy: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, output)
}
