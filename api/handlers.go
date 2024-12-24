// File: api/handlers.go

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
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Repo represents a GitHub repository
type Repo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	HTMLURL         string `json:"html_url"`
	Language        string `json:"language"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
}

// ConfigurePageData holds data to populate the configure form
type ConfigurePageData struct {
	RepoName   string
	FolderName string
	SubFolders string
}

// Global OAuth2 configuration
var oauth2Config *oauth2.Config

// Session store
var store = sessions.NewCookieStore([]byte("super-secret-key")) // Replace with a secure key in production

// Initialize the OAuth2 configuration
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
		Scopes:       []string{"repo"}, // Adjust scopes as needed
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/github/callback",
	}

	// Validate environment variables
	if oauth2Config.ClientID == "" || oauth2Config.ClientSecret == "" {
		log.Fatalf("Missing GITHUB_CLIENT_ID or GITHUB_CLIENT_SECRET in .env file")
	}

	// Configure session options
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure:   true, // Uncomment when using HTTPS
	}
}

// generateState creates a random state string for CSRF protection
func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Error generating random state: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

// LoginHandler redirects the user to GitHub for login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Login handler triggered")

	// Generate state and store in session
	state := generateState()
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to GitHub OAuth URL
	authURL := oauth2Config.AuthCodeURL(state)
	log.Printf("Redirecting to GitHub OAuth URL: %s", authURL)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// CallbackHandler processes the GitHub OAuth callback
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback handler triggered")

	// Retrieve session
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Validate state to prevent CSRF
	oauthState, ok := session.Values["state"].(string)
	if !ok || r.FormValue("state") != oauthState {
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

	// Log full token details for debugging
	log.Printf("Full token: %+v", token)

	// Log token expiry
	log.Printf("Token expiry time: %v", token.Expiry)

	// Store token fields in session
	session.Values["access_token"] = token.AccessToken
	session.Values["token_type"] = token.TokenType
	session.Values["refresh_token"] = token.RefreshToken
	session.Values["expiry"] = token.Expiry.Format(time.RFC3339)
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Stored token fields in session")

	// Redirect to the get-code.html page
	http.Redirect(w, r, "/get-code.html", http.StatusSeeOther)
}

// GetReposHandler handles the /api/repos endpoint to fetch user repositories
func GetReposHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetReposHandler called")

	// Retrieve session
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve token fields from session
	accessToken, ok := session.Values["access_token"].(string)
	if !ok || accessToken == "" {
		log.Println("Access token not found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenType, ok := session.Values["token_type"].(string)
	if !ok {
		tokenType = "Bearer" // Default token type
	}

	refreshToken, _ := session.Values["refresh_token"].(string)

	expiryStr, ok := session.Values["expiry"].(string)
	var expiry time.Time
	if ok && expiryStr != "" {
		expiry, err = time.Parse(time.RFC3339, expiryStr)
		if err != nil {
			log.Printf("Error parsing expiry time: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the token is expired only if Expiry is set
		if !expiry.IsZero() && expiry.Before(time.Now()) {
			log.Printf("Token expired at %v", expiry)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("Token expiry time: %v", expiry)
	} else {
		log.Println("Expiry not found in session or zero value; assuming token does not expire")
	}

	log.Println("Access token retrieved successfully")

	// Reconstruct the oauth2.Token
	token := &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}

	// Create an OAuth2 client with the token
	client := oauth2Config.Client(context.Background(), token)

	// Fetch repositories from GitHub API
	repos, err := fetchAllRepos(client)
	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
		http.Error(w, "Failed to fetch repositories", http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched %d repositories", len(repos))

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode repositories to JSON
	if err := json.NewEncoder(w).Encode(repos); err != nil {
		log.Printf("Error encoding repositories to JSON: %v", err)
		http.Error(w, "Failed to encode repositories", http.StatusInternalServerError)
		return
	}
}

// fetchAllRepos retrieves all repositories, handling pagination
func fetchAllRepos(client *http.Client) ([]Repo, error) {
	var allRepos []Repo
	url := "https://api.github.com/user/repos?per_page=100"

	for url != "" {
		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
		}

		var repos []Repo
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		// Check for pagination
		links := resp.Header.Get("Link")
		url = parseNextLink(links)
	}

	return allRepos, nil
}

// parseNextLink parses the 'Link' header to find the next page URL
func parseNextLink(links string) string {
	if links == "" {
		return ""
	}
	parts := strings.Split(links, ",")
	for _, part := range parts {
		section := strings.Split(strings.TrimSpace(part), ";")
		if len(section) < 2 {
			continue
		}
		url := strings.Trim(section[0], "<>")
		rel := strings.TrimSpace(section[1])
		if rel == `rel="next"` {
			return url
		}
	}
	return ""
}

// ConfigureHandler processes the form submission from get-code.html
func ConfigureHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ConfigureHandler called")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	data := ConfigurePageData{
		RepoName:   r.FormValue("repoName"),
		FolderName: r.FormValue("folderName"),
		SubFolders: r.FormValue("subFolders"),
	}

	// TODO: Implement your logic to handle the configuration
	// For example, cloning the repository, setting up folders, etc.
	log.Printf("Configuration received: %+v", data)

	// Redirect to a success page or display a success message
	http.Redirect(w, r, "/success.html", http.StatusSeeOther)
}

// DebugHandler lists session data (for debugging purposes only)
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DebugHandler called")

	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// For security, restrict access to localhost
	remoteAddr := r.RemoteAddr
	allowedAddrs := []string{"127.0.0.1:8080", "[::1]:8080"}
	allowed := false
	for _, addr := range allowedAddrs {
		if remoteAddr == addr {
			allowed = true
			break
		}
	}
	if !allowed {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Serialize session values to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.Values)
}
