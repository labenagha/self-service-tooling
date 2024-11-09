package api

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOauthConfig = &oauth2.Config{
	ClientID:     "", // Replace with your Client ID
	ClientSecret: "", // Replace with your Client Secret
	RedirectURL:  "http://localhost:8080/auth/callback",
	Scopes:       []string{"repo"},
	Endpoint:     github.Endpoint,
}

// HandleGitHubLogin redirects the user to GitHub for login
func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleGitHubCallback processes the GitHub callback
func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Access Token: %s", token.AccessToken)
}
