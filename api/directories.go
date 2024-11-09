package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// ListRepositoryDirectories lists directories within a selected repository.
func ListRepositoryDirectories(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "Repository not specified", http.StatusBadRequest)
		return
	}

	token := os.Getenv("GITHUB_TOKEN") // Ensure the token is securely set
	if token == "" {
		http.Error(w, "GitHub token not set", http.StatusInternalServerError)
		return
	}

	client := oauth2.NewClient(r.Context(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/contents", repo)
	resp, err := client.Get(apiUrl)
	if err != nil {
		http.Error(w, "Failed to fetch directories: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "GitHub API error: "+resp.Status, http.StatusInternalServerError)
		return
	}

	var contents []struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Path string `json:"path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		http.Error(w, "Failed to decode response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter directories only
	directories := []string{}
	for _, item := range contents {
		if item.Type == "dir" {
			directories = append(directories, item.Path)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(directories)
}
