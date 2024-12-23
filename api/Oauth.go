package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Global OAuth2 configuration
var oauth2Config *oauth2.Config

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize OAuth2 configuration
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"read:user"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/github/callback",
	}

	// Validate environment variables
	if oauth2Config.ClientID == "" || oauth2Config.ClientSecret == "" {
		log.Fatalf("Missing GITHUB_CLIENT_ID or GITHUB_CLIENT_SECRET in .env file")
	}
}

// LoginHandler redirects the user to GitHub for login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Login handler triggered")

	// Generate state and set cookie for CSRF protection
	state := randomState()
	http.SetCookie(w, &http.Cookie{
		Name:  "oauthstate",
		Value: state,
		Path:  "/",
	})
	authURL := oauth2Config.AuthCodeURL(state)
	log.Printf("Redirecting to GitHub OAuth URL: %s", authURL)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// CallbackHandler processes the GitHub OAuth callback
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback handler triggered")

	// Validate state to prevent CSRF
	oauthState, err := r.Cookie("oauthstate")
	if err != nil || r.FormValue("state") != oauthState.Value {
		log.Println("Invalid OAuth state")
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	// Get the authorization code from the query parameters
	code := r.FormValue("code")
	log.Printf("Authorization code: %s", code)

	// Exchange the code for an access token
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging token: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Use the token to fetch user information
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Printf("Error fetching user info: %v", err)
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the user's information
	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Printf("Error decoding user info: %v", err)
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Construct the URL to the user's repositories page
	githubRepositoriesURL := fmt.Sprintf("https://github.com/%s?tab=repositories", user.Login)
	log.Printf("Redirecting user to: %s", githubRepositoriesURL)
	http.Redirect(w, r, githubRepositoriesURL, http.StatusFound)
}

// randomState generates a random string for CSRF protection
func randomState() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Error generating random state: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
