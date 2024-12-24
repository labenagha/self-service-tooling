// Grab references to DOM elements
const terraformOutput = document.getElementById("terraform-output");
const approveBtn = document.getElementById("approve-btn");
const cancelBtn = document.getElementById("cancel-btn");

// Simulate Terraform output
function simulateTerraformOutput() {
    const messages = [
        "Initializing modules...",
        "Downloading providers...",
        "Creating infrastructure...",
        "Applying changes...",
        "Terraform Apply complete!"
    ];

    let index = 0;

    const interval = setInterval(() => {
        if (index < messages.length) {
            terraformOutput.textContent += `\n${messages[index]}`;
            index++;
            terraformOutput.scrollTop = terraformOutput.scrollHeight; // Auto-scroll to bottom
        } else {
            clearInterval(interval);
        }
    }, 2000); // Simulate output every 2 seconds
}

// Handle Approve button click
approveBtn.addEventListener("click", () => {
    alert("Terraform changes approved!");
    terraformOutput.textContent += "\nApproval received. Proceeding with changes...";
});

// Handle Cancel button click
cancelBtn.addEventListener("click", () => {
    alert("Terraform apply canceled.");
    terraformOutput.textContent += "\nOperation canceled by user.";
});

// // Start simulating Terraform output when the page loads
// simulateTerraformOutput();
