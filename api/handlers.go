package api

import (
	"fmt"
	"net/http"
)

func GetTerraformCode(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This will fetch and return Terraform code.")
}

// DeployTerraformCode handles the deployment request
func DeployTerraformCode(w http.ResponseWriter, r *http.Request) {
	// Ensure this is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	output, err := RunTerraformCommand("terraform apply -auto-approve") // Adjust this as needed
	if err != nil {
		http.Error(w, "Error running Terraform apply: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, output)
}

func RunTerraformPlan(w http.ResponseWriter, r *http.Request) {
	output, err := RunTerraformCommand("terraform plan") // Direct call
	if err != nil {
		http.Error(w, "Error running Terraform plan: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, output) // Ensure the cleaned output is sent back
}

func DestroyTerraformCode(w http.ResponseWriter, r *http.Request) {
	output, err := RunTerraformCommand("terraform destroy -auto-approve") // Direct call
	if err != nil {
		http.Error(w, "Error running Terraform destroy: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, output)
}
