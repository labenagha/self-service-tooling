package api

import (
	"net/http"
	"os/exec"
)

// TerraformPlanHandler runs "terraform plan"
func TerraformPlanHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("terraform", "plan")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, string(output), http.StatusInternalServerError)
		return
	}
	w.Write(output)
}
