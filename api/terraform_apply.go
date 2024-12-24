// package api

// import (
// 	"net/http"
// 	"os/exec"
// )

// // TerraformApplyHandler runs "terraform apply -auto-approve"
// func TerraformApplyHandler(w http.ResponseWriter, r *http.Request) {
// 	cmd := exec.Command("terraform", "apply", "-auto-approve")
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		http.Error(w, string(output), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(output)
// }
