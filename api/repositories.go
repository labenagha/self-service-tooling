package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// ListUserRepositories lists the repositories of the authenticated user.
func ListUserRepositories(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		http.Error(w, "GitHub token not set", http.StatusInternalServerError)
		return
	}

	client := oauth2.NewClient(r.Context(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	apiUrl := "https://api.github.com/user/repos"

	resp, err := client.Get(apiUrl)
	if err != nil {
		http.Error(w, "Failed to fetch repositories from GitHub: "+err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error fetching repositories: %v\n", err) // Debugging log
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "GitHub API error: "+resp.Status, http.StatusInternalServerError)
		fmt.Printf("GitHub API error: %s\n", resp.Status) // Debugging log
		return
	}

	var repos []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		http.Error(w, "Failed to decode response body: "+err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error decoding response body: %v\n", err) // Debugging log
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(repos); err != nil {
		fmt.Printf("Error encoding response: %v\n", err) // Debugging log
		http.Error(w, "Failed to send JSON response", http.StatusInternalServerError)
	}
}
