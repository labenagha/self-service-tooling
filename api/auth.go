// File: api/auth.go

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"self-service-tooling/config" // Replace with your actual module path

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Initialize OAuth2 configuration
var githubOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Scopes:       []string{"repo"},
	Endpoint:     github.Endpoint,
	RedirectURL:  "http://localhost:8080/auth/github/callback", // Adjust as needed
}

// // GitHubOAuthHandler initiates the GitHub OAuth flow
// func GitHubOAuthHandler(w http.ResponseWriter, r *http.Request) {
// 	state := generateState()
// 	session, err := config.Store.Get(r, "session-name")
// 	if err != nil {
// 		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
// 		log.Printf("Failed to get session: %v", err)
// 		return
// 	}
// 	session.Values["oauth_state"] = state
// 	err = session.Save(r, w)
// 	if err != nil {
// 		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	url := githubOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

// GitHubCallbackHandler handles the OAuth callback from GitHub
func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found in query parameters", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	token, err := githubOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
		log.Printf("OAuth exchange error: %v", err)
		return
	}

	// Create a GitHub client with the received token
	client := githubOauthConfig.Client(r.Context(), token)

	// Retrieve user information from GitHub API
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to retrieve user info: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to retrieve user info: %v", err)
		return
	}
	defer resp.Body.Close()

	// Parse the response to get the GitHub username
	var githubUser struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to decode user info: %v", err)
		return
	}

	// Get session to store the token and authentication state
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to get session: %v", err)
		return
	}

	// Set session values for authentication
	session.Values["github_access_token"] = token.AccessToken
	session.Values["authenticated"] = true
	session.Values["username"] = githubUser.Login

	// Save the session
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to save session: %v", err)
		return
	}

	// Redirect the user to the dashboard after successful login
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
