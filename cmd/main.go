// File: cmd/main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"self-service-tooling/api"    // Replace with your actual module path
	"self-service-tooling/config" // Replace with your actual module path

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Global in-memory user store (for demonstration purposes)
// In production, use a persistent database
var userStore = map[string]User{}

// User represents a registered user.
type User struct {
	Username     string
	PasswordHash string
}

// GitHubUser represents the authenticated GitHub user.
type GitHubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

// GitHub OAuth Configurations
var (
	githubClientID     = os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	githubRedirectURL  = "http://localhost:8080/auth/github/callback" // Update if different
)

// Initialize Templates and GitHub OAuth configurations
func init() {
	// Ensure GitHub OAuth credentials are set
	if githubClientID == "" || githubClientSecret == "" {
		log.Fatal("GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET must be set as environment variables")
	}
}

// Function to generate a random state string for CSRF protection (implement this function)
func generateRandomState() string {
	// Your implementation to generate a random string
	return "randomStateString"
}

// Step to initiate the GitHub OAuth flow
func initiateOAuthFlowHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState() // Function to generate a random state string
	session, _ := config.Store.Get(r, "session-name")
	session.Values["state"] = state
	session.Save(r, w)

	// Redirect the user to GitHub with the state parameter
	http.Redirect(w, r, "https://github.com/login/oauth/authorize?client_id=YOUR_CLIENT_ID&state="+state, http.StatusFound)
}

// renderTemplate renders an HTML template with provided data.
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("ui/" + tmpl)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
	}
}

// isAuthenticated checks if the user is authenticated.
func isAuthenticated(r *http.Request) bool {
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		log.Println("Failed to get session:", err)
		return false
	}

	// Check if the 'authenticated' flag is set in the session
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

// homeHandler redirects users based on authentication status.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
}

// loginPageHandler serves the login page.
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, "login.html", nil)
}

// registerPageHandler serves the registration page.
func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, "register.html", nil)
}

// registerHandler handles user registration.
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON body
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Invalid request body in registerHandler:", err)
		return
	}

	username := strings.TrimSpace(payload.Username)
	password := strings.TrimSpace(payload.Password)

	if username == "" || password == "" {
		http.Error(w, "Username and password cannot be empty", http.StatusBadRequest)
		log.Println("Empty username or password in registerHandler")
		return
	}

	// Check if user already exists
	if _, exists := userStore[username]; exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		log.Println("Username already exists:", username)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		log.Println("Error hashing password:", err)
		return
	}

	// Store the user
	userStore[username] = User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	log.Println("User registered successfully:", username)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Registration successful. Redirecting to login...")
}

// loginHandler handles user login.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON body
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Invalid request body in loginHandler:", err)
		return
	}

	username := strings.TrimSpace(payload.Username)
	password := strings.TrimSpace(payload.Password)

	if username == "" || password == "" {
		http.Error(w, "Username and password cannot be empty", http.StatusBadRequest)
		log.Println("Empty username or password in loginHandler")
		return
	}

	// Retrieve user
	user, exists := userStore[username]
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Println("User not found in loginHandler:", username)
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Println("Invalid password for user:", username)
		return
	}

	// Save user info in session
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to get session in loginHandler:", err)
		return
	}

	session.Values["authenticated"] = true
	session.Values["username"] = user.Username
	// Initialize selections as empty
	session.Values["selected_repo"] = ""
	session.Values["selected_directory"] = ""
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to save session in loginHandler:", err)
		return
	}

	log.Println("User logged in successfully via username/password:", username)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Login successful. Redirecting to dashboard...")
}

// operationsHandler serves the main operations dashboard.
func operationsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve session
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to get session in operationsHandler:", err)
		return
	}

	// Extract session values
	username, _ := session.Values["username"].(string)
	repo, repoExists := session.Values["selected_repo"].(string)
	directory, dirExists := session.Values["selected_directory"].(string)

	log.Printf("Rendering operations page for user '%s': Repo='%s', Directory='%s'\n", username, repo, directory)

	// If repo or directory not set, display "Not Selected"
	if !repoExists || repo == "" {
		repo = "Not Selected"
	}
	if !dirExists || directory == "" {
		directory = "Not Selected"
	}

	data := struct {
		Username  string
		Repo      string
		Directory string
		IsRepoSet bool
		IsDirSet  bool
	}{
		Username:  username,
		Repo:      repo,
		Directory: directory,
		IsRepoSet: repo != "Not Selected",
		IsDirSet:  directory != "Not Selected",
	}

	renderTemplate(w, "index.html", data)
}

// logoutHandler handles user logout.
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to get session in logoutHandler:", err)
		return
	}

	// Revoke user's authentication
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Values["selected_repo"] = ""
	session.Values["selected_directory"] = ""
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to save session in logoutHandler:", err)
		return
	}

	log.Println("User logged out successfully")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// GitHubOAuthHandler initiates GitHub OAuth login.
func GitHubOAuthHandler(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to get session in GitHubOAuthHandler:", err)
		return
	}

	session.Values["oauth_state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to save session in GitHubOAuthHandler:", err)
		return
	}

	authURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=repo",
		githubClientID, githubRedirectURL, state)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GitHubCallbackHandler handles the OAuth callback from GitHub.
func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to get session in GitHubCallbackHandler:", err)
		return
	}

	// Verify state
	queryState := r.URL.Query().Get("state")
	storedState, ok := session.Values["oauth_state"].(string)
	if !ok || queryState != storedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		log.Println("Invalid state parameter in GitHubCallbackHandler")
		return
	}

	// Get code
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found in callback", http.StatusBadRequest)
		log.Println("Code not found in GitHubCallbackHandler")
		return
	}

	// Exchange code for access token
	token, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error exchanging code for token:", err)
		return
	}

	// Save access token in session
	session.Values["access_token"] = token
	// Clear oauth_state
	delete(session.Values, "oauth_state")
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to save session in GitHubCallbackHandler:", err)
		return
	}

	// Optionally, fetch GitHub user info
	user, err := fetchGitHubUser(token)
	if err != nil {
		http.Error(w, "Failed to fetch GitHub user: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error fetching GitHub user:", err)
		return
	}

	// Save GitHub username in session
	session.Values["github_username"] = user.Login
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to save session in GitHubCallbackHandler:", err)
		return
	}

	// Redirect to repository selection page
	http.Redirect(w, r, "/list-repos", http.StatusSeeOther)
}

// generateState generates a random string for OAuth state parameter.
// For production, use a secure random generator.
func generateState() string {
	return fmt.Sprintf("state-%d", time.Now().UnixNano())
}

// exchangeCodeForToken exchanges the authorization code for an access token.
func exchangeCodeForToken(code string) (string, error) {
	data := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&state=%s",
		githubClientID, githubClientSecret, code, githubRedirectURL, generateState())

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respData struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", err
	}

	if respData.AccessToken == "" {
		return "", fmt.Errorf("no access token received")
	}

	return respData.AccessToken, nil
}

// fetchGitHubUser fetches the authenticated user's GitHub username.
func fetchGitHubUser(token string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
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

	var user GitHubUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// handleTerraformAction returns a handler function for the specified Terraform action.
func handleTerraformAction(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse JSON body
		var payload struct {
			Repo      string `json:"repo"`
			Directory string `json:"directory"`
		}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Printf("Invalid request body in handleTerraformAction (%s): %v", action, err)
			return
		}

		repo := strings.TrimSpace(payload.Repo)
		directory := strings.TrimSpace(payload.Directory)

		if repo == "" || directory == "" {
			http.Error(w, "Repository and directory cannot be empty", http.StatusBadRequest)
			log.Printf("Empty repository or directory in handleTerraformAction (%s)", action)
			return
		}

		// Execute the Terraform command
		output, err := executeTerraformCommand(repo, directory, action)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing Terraform command (%s): %v", action, err)
			return
		}

		// Return the output
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, output)
	}
}

// executeTerraformCommand runs a Terraform command in the specified repository and directory.
// It returns the combined stdout and stderr output or an error.
func executeTerraformCommand(repo, directory, action string) (string, error) {
	// Define the path to the repository and directory
	// Replace "/path/to/repos" with your actual path
	repoPath := fmt.Sprintf("/path/to/repos/%s/%s", repo, directory)

	// Clone the repository if it doesn't exist
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		// Extract owner and repo name from full repo name (owner/repo)
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repository format: %s", repo)
		}
		owner := parts[0]
		repoName := parts[1]

		cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repoName)

		// Clone the repository
		cmd := exec.Command("git", "clone", cloneURL, repoPath)
		var outBuf, errBuf strings.Builder
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf

		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("failed to clone repository: %v\nStderr: %s", err, errBuf.String())
		}
	}

	// Ensure the directory exists within the repository
	targetDir := repoPath // Assuming the repository root is the target
	if directory != "" && directory != "/" {
		targetDir = fmt.Sprintf("%s/%s", repoPath, directory)
	}
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return "", fmt.Errorf("directory does not exist: %s", targetDir)
	}

	// Create a context with a timeout to prevent long-running operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Initialize the command based on the action
	var cmd *exec.Cmd
	switch action {
	case "plan":
		cmd = exec.CommandContext(ctx, "terraform", "plan")
	case "apply":
		// To automate 'apply' without manual confirmation, use the '-auto-approve' flag
		cmd = exec.CommandContext(ctx, "terraform", "apply", "-auto-approve")
	case "destroy":
		// To automate 'destroy' without manual confirmation, use the '-auto-approve' flag
		cmd = exec.CommandContext(ctx, "terraform", "destroy", "-auto-approve")
	case "getcode":
		// Example: Initialize the directory
		cmd = exec.CommandContext(ctx, "terraform", "init")
	default:
		return "", fmt.Errorf("invalid Terraform action: %s", action)
	}

	// Set the working directory
	cmd.Dir = targetDir

	// Optional: Set environment variables if needed
	// cmd.Env = append(os.Environ(), "TF_VAR_example=value")

	// Capture the output
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %v\nStderr: %s", err, errBuf.String())
	}

	// Combine stdout and stderr
	output := outBuf.String() + errBuf.String()
	return output, nil
}

func main() {
	// Initialize router
	r := mux.NewRouter()

	// Set up routes
	http.HandleFunc("/auth/github/start", initiateOAuthFlowHandler)

	// Routes for pages
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/login", loginPageHandler).Methods("GET")
	r.HandleFunc("/register", registerPageHandler).Methods("GET")
	r.HandleFunc("/auth/register", registerHandler).Methods("POST")
	r.HandleFunc("/auth/login", loginHandler).Methods("POST")
	r.HandleFunc("/operations", operationsHandler).Methods("GET")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

	// Routes for GitHub OAuth
	r.HandleFunc("/auth/github/login", GitHubOAuthHandler).Methods("GET")
	r.HandleFunc("/auth/github/callback", GitHubCallbackHandler).Methods("GET")

	// Routes for repository and directory selection
	r.HandleFunc("/list-repos", api.ListUserRepositories).Methods("GET", "POST")
	r.HandleFunc("/api/directories", api.ListRepositoryDirectories).Methods("GET")

	// Routes for Terraform actions
	r.HandleFunc("/terraform-plan", handleTerraformAction("plan")).Methods("POST")
	r.HandleFunc("/terraform-apply", handleTerraformAction("apply")).Methods("POST")
	r.HandleFunc("/terraform-destroy", handleTerraformAction("destroy")).Methods("POST")
	r.HandleFunc("/terraform-getcode", handleTerraformAction("getcode")).Methods("POST")

	// Route for API token generation
	r.HandleFunc("/api/generate-token", api.GenerateAPITokenHandler).Methods("POST")

	// Serve static files from /ui/ directory
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
