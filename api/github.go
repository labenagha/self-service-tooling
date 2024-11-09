package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// FetchRepoFile handles the HTTP request for fetching the content of a file from a GitHub repository.
func FetchRepoFile(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "Repository not specified", http.StatusBadRequest)
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		http.Error(w, "GitHub token not set", http.StatusInternalServerError)
		return
	}

	client := oauth2.NewClient(r.Context(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/contents/main.tf", repo)

	resp, err := client.Get(apiUrl)
	if err != nil {
		http.Error(w, "Failed to fetch file from GitHub: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "GitHub API error: "+resp.Status, http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
