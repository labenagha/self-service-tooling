// File: api/list-repos.go

package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"self-service-tooling/config" // Replace with your actual module path
)

// Repository represents a GitHub repository.
type Repository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"` // Format: "owner/repo"
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
	// Add other fields if needed
}

// Content represents a file or directory in a GitHub repository.
type Content struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"` // "file" or "dir"
	DownloadURL string `json:"download_url"`
}

// User represents a user with an associated GitHub access token.
type User struct {
	Username          string
	GitHubAccessToken string
}

// In-memory token store (replace with persistent storage in production)
var tokenStore = map[string]User{}

// GenerateAPITokenHandler allows authenticated users to generate API tokens
func GenerateAPITokenHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve session
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		log.Println("Failed to get session in GenerateAPITokenHandler:", err)
		return
	}

	// Check if user is authenticated
	authenticated, ok := session.Values["authenticated"].(bool)
	username, uok := session.Values["username"].(string)

	if !ok || !authenticated || !uok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		log.Println("User not authenticated in GenerateAPITokenHandler")
		return
	}

	// Generate a new API token
	token, err := generateAPIToken()
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Println("Failed to generate API token:", err)
		return
	}

	// Retrieve GitHub access token from session
	githubAccessToken, tokenOk := session.Values["github_access_token"].(string)
	if !tokenOk || githubAccessToken == "" {
		http.Error(w, "GitHub OAuth not completed", http.StatusUnauthorized)
		log.Println("GitHub access token not found for user:", username)
		return
	}

	// Assign the token to the user
	AssignTokenToUser(username, token, githubAccessToken)

	// Return the token to the user
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AssignTokenToUser assigns a generated token to a user.
func AssignTokenToUser(username string, token string, githubAccessToken string) {
	user := User{
		Username:          username,
		GitHubAccessToken: githubAccessToken,
	}

	tokenStore[token] = user
	log.Printf("Assigned API Token to user '%s': %s\n", username, token)
}

// validateAPIToken validates the provided API token and returns the associated user.
func validateAPIToken(token string) (User, bool) {
	user, exists := tokenStore[token]
	if !exists {
		return User{}, false
	}
	return user, true
}

// generateAPIToken generates a secure random token.
func generateAPIToken() (string, error) {
	bytes := make([]byte, 32) // 256-bit token
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ListRepositoryDirectories handles fetching directories for a selected repository
func ListRepositoryDirectories(w http.ResponseWriter, r *http.Request) {
	// First, attempt session-based authentication
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to get session in ListRepositoryDirectories:", err)
		return
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	username, uok := session.Values["username"].(string)

	if ok && authenticated && uok {
		// Session-based authenticated user
		log.Printf("Session Authenticated: %v, Username: %s", authenticated, username)

		// Retrieve GitHub access token from session
		githubAccessToken, tokenOk := session.Values["github_access_token"].(string)
		if !tokenOk || githubAccessToken == "" {
			http.Error(w, "GitHub OAuth not completed", http.StatusUnauthorized)
			log.Println("GitHub OAuth not completed for user:", username)
			return
		}

		// Extract repository from query
		repo := r.URL.Query().Get("repo")
		repo = strings.TrimSpace(repo)
		if repo == "" {
			http.Error(w, "Repository not specified", http.StatusBadRequest)
			log.Println("Repository not specified in ListRepositoryDirectories")
			return
		}

		// Fetch directories from GitHub
		dirs, err := fetchGitHubDirectories(githubAccessToken, repo)
		if err != nil {
			http.Error(w, "Failed to fetch directories: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error fetching directories:", err)
			return
		}

		// Structure the response
		response := struct {
			Directories []Content `json:"directories"`
		}{
			Directories: dirs,
		}

		// Send JSON response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error encoding directories response:", err)
			return
		}

		log.Printf("Sent directories for repository '%s' to user '%s'\n", repo, username)
		return
	}

	// If session-based authentication fails, attempt token-based authentication
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			token := parts[1]
			user, tokenValid := validateAPIToken(token)
			if tokenValid {
				log.Printf("Token Authenticated: Username: %s", user.Username)

				// Extract repository from query
				repo := r.URL.Query().Get("repo")
				repo = strings.TrimSpace(repo)
				if repo == "" {
					http.Error(w, "Repository not specified", http.StatusBadRequest)
					log.Println("Repository not specified in ListRepositoryDirectories (Token Auth)")
					return
				}

				// Fetch directories from GitHub
				dirs, err := fetchGitHubDirectories(user.GitHubAccessToken, repo)
				if err != nil {
					http.Error(w, "Failed to fetch directories: "+err.Error(), http.StatusInternalServerError)
					log.Println("Error fetching directories:", err)
					return
				}

				// Structure the response
				response := struct {
					Directories []Content `json:"directories"`
				}{
					Directories: dirs,
				}

				// Send JSON response
				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(response)
				if err != nil {
					http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
					log.Println("Error encoding directories response:", err)
					return
				}

				log.Printf("Sent directories for repository '%s' to user '%s' via Token Auth\n", repo, user.Username)
				return
			}
		}
	}

	// If both authentication methods fail, respond with unauthorized
	http.Error(w, "User not authenticated", http.StatusUnauthorized)
	log.Println("User not authenticated in ListRepositoryDirectories")
}

// ListUserRepositories handles the /list-repos route.
// It fetches the user's GitHub repositories and renders them in the template.
func ListUserRepositories(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if access_token exists
	accessToken, ok := session.Values["access_token"].(string)
	if !ok || accessToken == "" {
		// Redirect to GitHub OAuth login
		http.Redirect(w, r, "/auth/github/login", http.StatusSeeOther)
		return
	}

	// Handle GET and POST methods
	switch r.Method {
	case http.MethodGet:
		// Fetch repositories from GitHub
		repos, err := fetchGitHubRepositories(accessToken)
		if err != nil {
			http.Error(w, "Failed to fetch repositories: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error fetching repositories:", err)
			return
		}

		// Render the list-repos.html template with repositories data
		tmpl, err := template.ParseFiles("ui/list-repos.html")
		if err != nil {
			http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error parsing template:", err)
			return
		}

		data := struct {
			Repositories []Repository
		}{
			Repositories: repos,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error executing template:", err)
			return
		}

	case http.MethodPost:
		// Handle repository and directory selection
		var payload struct {
			Repo      string `json:"repo"`
			Directory string `json:"directory"`
		}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Println("Invalid request body in ListUserRepositories POST:", err)
			return
		}

		repo := strings.TrimSpace(payload.Repo)
		directory := strings.TrimSpace(payload.Directory)

		if repo == "" || directory == "" {
			http.Error(w, "Repository and directory cannot be empty", http.StatusBadRequest)
			log.Println("Empty repository or directory in ListUserRepositories POST")
			return
		}

		// Validate repository
		repos, err := fetchGitHubRepositories(accessToken)
		if err != nil {
			http.Error(w, "Failed to fetch repositories: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error fetching repositories:", err)
			return
		}

		validRepo := false
		for _, r := range repos {
			if r.FullName == repo {
				validRepo = true
				break
			}
		}
		if !validRepo {
			http.Error(w, "Selected repository does not exist", http.StatusBadRequest)
			log.Println("Selected repository does not exist:", repo)
			return
		}

		// Fetch directories from the selected repository
		dirs, err := fetchGitHubDirectories(accessToken, repo)
		if err != nil {
			http.Error(w, "Failed to fetch directories: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error fetching directories:", err)
			return
		}

		validDir := false
		for _, d := range dirs {
			if d.Path == directory && d.Type == "dir" {
				validDir = true
				break
			}
		}
		if !validDir {
			http.Error(w, "Selected directory does not exist in the repository", http.StatusBadRequest)
			log.Println("Selected directory does not exist in repository:", directory)
			return
		}

		// Save selected repository and directory in session
		session.Values["selected_repo"] = repo
		session.Values["selected_directory"] = directory
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
			log.Println("Failed to save session in ListUserRepositories POST:", err)
			return
		}

		// Respond with success message
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Repository and directory selected successfully. Redirecting to dashboard...")
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// fetchGitHubDirectories fetches directories from a specific GitHub repository using the provided access token.
func fetchGitHubDirectories(token, fullRepoName string) ([]Content, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/contents", fullRepoName)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository not found: %s", fullRepoName)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var contents []Content
	err = json.NewDecoder(resp.Body).Decode(&contents)
	if err != nil {
		return nil, err
	}

	// Filter directories
	var dirs []Content
	for _, item := range contents {
		if item.Type == "dir" {
			dirs = append(dirs, item)
		}
	}

	return dirs, nil
}

// fetchGitHubRepositories fetches the user's repositories from GitHub
func fetchGitHubRepositories(token string) ([]Repository, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/repos", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var repos []Repository
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
