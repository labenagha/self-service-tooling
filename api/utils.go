package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"golang.org/x/oauth2"
)

// FileInfo represents the structure of a file or directory returned by the GitHub API.
type FileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

// DownloadRepositoryContent downloads the contents of a directory from a GitHub repository.
func DownloadRepositoryContent(repo, path, token string) error {
	ctx := context.Background() // Create a context
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo, path)
	resp, err := client.Get(apiUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch repository content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var files []FileInfo
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	// Create local directory for storing downloaded files
	localDir := "./terraform-code"
	if err := os.MkdirAll(localDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Download each file
	for _, file := range files {
		if file.Type == "file" {
			if err := downloadFile(file.DownloadURL, filepath.Join(localDir, file.Name), token); err != nil {
				return fmt.Errorf("failed to download file %s: %w", file.Name, err)
			}
		}
	}

	return nil
}

// downloadFile downloads a file from the specified URL and saves it to the local path.
func downloadFile(url, localPath, token string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := os.WriteFile(localPath, body, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// RunTerraformCommand runs the provided Terraform command and returns the output and errors.
func RunTerraformCommand(command, workingDir string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = workingDir

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("failed to run command: %s, error: %s", command, err)
	}
	return StripColorCodes(out.String()), nil
}

// StripColorCodes removes ANSI color codes from the output.
func StripColorCodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}
