// package api

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"

// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/github"
// )

// // This config is separate from oauth2Config in oauth.go because it uses "repo" scope.
// var repoOauth2Config *oauth2.Config

// func init() {
// 	// Only do this if you need a separate config; otherwise, reuse oauth2Config from oauth.go.
// 	repoOauth2Config = &oauth2.Config{
// 		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
// 		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
// 		Endpoint:     github.Endpoint,
// 		Scopes:       []string{"repo"},
// 	}
// }

// // GetRepositoriesHandler dynamically fetches repositories for the authenticated user.
// func GetRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get the token from the query parameter (in a real app, store it securely in session)
// 	token := r.URL.Query().Get("access_token")
// 	if token == "" {
// 		http.Error(w, "Access token is missing", http.StatusUnauthorized)
// 		return
// 	}

// 	// Use the token to create an authenticated HTTP client
// 	client := repoOauth2Config.Client(r.Context(), &oauth2.Token{AccessToken: token})

// 	// Fetch repositories from GitHub API
// 	resp, err := client.Get("https://api.github.com/user/repos")
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to fetch repositories: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Check if GitHub API returned an error
// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(resp.Body)
// 		http.Error(w, fmt.Sprintf("GitHub API error: %s", body), resp.StatusCode)
// 		return
// 	}

// 	// Decode the JSON response
// 	var repos []map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to parse GitHub response: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with the fetched repositories
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(repos)
// }
