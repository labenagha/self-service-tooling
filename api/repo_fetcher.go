// package api

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// type RepoData struct {
// 	Name    string   `json:"name"`
// 	Folders []string `json:"folders"`
// }

// // Fetches repository details and folders from GitHub API
// func FetchRepoDetails(owner, repo string) (RepoData, error) {
// 	// GitHub API URLs
// 	repoURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
// 	contentsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", owner, repo)

// 	// Fetch repository details
// 	resp, err := http.Get(repoURL)
// 	if err != nil {
// 		return RepoData{}, fmt.Errorf("failed to fetch repo details: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	var repoDetails struct {
// 		Name string `json:"name"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&repoDetails); err != nil {
// 		return RepoData{}, fmt.Errorf("failed to decode repo details: %v", err)
// 	}

// 	// Fetch repository contents (folders only)
// 	resp, err = http.Get(contentsURL)
// 	if err != nil {
// 		return RepoData{}, fmt.Errorf("failed to fetch repo contents: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	var contents []struct {
// 		Name string `json:"name"`
// 		Type string `json:"type"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
// 		return RepoData{}, fmt.Errorf("failed to decode repo contents: %v", err)
// 	}

// 	// Filter only directories
// 	var folders []string
// 	for _, item := range contents {
// 		if item.Type == "dir" {
// 			folders = append(folders, item.Name)
// 		}
// 	}

// 	return RepoData{
// 		Name:    repoDetails.Name,
// 		Folders: folders,
// 	}, nil
// }
