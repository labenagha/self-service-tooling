package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
)

// RunTerraformCommand runs the provided Terraform command and returns the output and errors.
func RunTerraformCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command) // Use bash to handle the command
	cmd.Dir = "./terraform"                    // Adjust path as needed

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

// StripColorCodes removes ANSI color codes from the output
func StripColorCodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}

// DeployHandler handles the /deploy endpoint for Terraform deployment
func DeployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	output, err := RunTerraformCommand("terraform apply -auto-approve")
	if err != nil {
		log.Printf("Error running Terraform apply: %v\n", err) // Log the error for debugging
		http.Error(w, "Error running Terraform apply: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, output)
}
