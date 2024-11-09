package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// In-memory user storage (for demo purposes; use a database for production)
var userStore = struct {
	sync.RWMutex
	users map[string]string
}{users: make(map[string]string)}

// Handler for user registration
func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple user creation logic
	userStore.Lock()
	defer userStore.Unlock()

	if _, exists := userStore.users[credentials.Username]; exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	userStore.users[credentials.Username] = credentials.Password
	fmt.Fprintln(w, "success")
}

// Handler for user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userStore.RLock()
	defer userStore.RUnlock()

	if storedPassword, exists := userStore.users[credentials.Username]; exists && storedPassword == credentials.Password {
		// Authentication successful
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "success")
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}

// Handler for listing repositories
func listReposHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder logic to simulate fetching repository data
	repos := []string{"Repo1", "Repo2", "Repo3"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

// Handler for running Terraform plan
func terraformPlanHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Repo      string `json:"repo"`
		Directory string `json:"directory"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate a Terraform plan execution
	fmt.Fprintf(w, "Terraform plan executed successfully for repo: %s, directory: %s", request.Repo, request.Directory)
}

// Handler for running Terraform apply
func terraformApplyHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Repo      string `json:"repo"`
		Directory string `json:"directory"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate a Terraform apply execution
	fmt.Fprintf(w, "Terraform apply executed successfully for repo: %s, directory: %s", request.Repo, request.Directory)
}

// Handler for running Terraform destroy
func terraformDestroyHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Repo      string `json:"repo"`
		Directory string `json:"directory"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate a Terraform destroy execution
	fmt.Fprintf(w, "Terraform destroy executed successfully for repo: %s, directory: %s", request.Repo, request.Directory)
}

func main() {
	// Route for the registration page (set as the default landing page)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/register.html") // Serve the registration page by default
	})

	// Route for the login page
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/login.html")
	})

	// Handle user registration
	http.HandleFunc("/auth/register", registerUserHandler)

	// Handle user login
	http.HandleFunc("/auth/login", loginHandler)

	// Route for the main operations page
	http.HandleFunc("/operations", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html")
	})

	// Handle endpoints for Terraform actions
	http.HandleFunc("/list-repos", listReposHandler)
	http.HandleFunc("/terraform-plan", terraformPlanHandler)
	http.HandleFunc("/terraform-apply", terraformApplyHandler)
	http.HandleFunc("/terraform-destroy", terraformDestroyHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./ui"))
	http.Handle("/ui/", http.StripPrefix("/ui/", fs))

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
