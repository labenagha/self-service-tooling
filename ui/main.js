document.addEventListener('DOMContentLoaded', () => {
    if (sessionStorage.getItem('darkMode') === 'enabled') { // Changed from localStorage to sessionStorage
        document.body.classList.add('dark-mode');
    }

    const getCodeButton = document.getElementById('getCodeButton');
    const runPlanButton = document.getElementById('runPlanButton');
    const deployButton = document.getElementById('deployButton');
    const runDestroyButton = document.getElementById('runDestroyButton');
    const toggleDarkModeButton = document.getElementById('toggleDarkModeButton');
    const output = document.getElementById('output');

    if (!output) {
        console.error("Element with ID 'output' not found in the DOM.");
        return;
    }

    if (getCodeButton) {
        getCodeButton.addEventListener('click', () => {
            window.location.href = '/list-repos';
        });
    }

    if (runPlanButton) {
        runPlanButton.addEventListener('click', () => {
            runTerraformPlan();
        });
    }

    if (deployButton) {
        deployButton.addEventListener('click', () => {
            runTerraformAction('/terraform-apply');
        });
    }

    if (runDestroyButton) {
        runDestroyButton.addEventListener('click', () => {
            runTerraformAction('/terraform-destroy');
        });
    }

    if (toggleDarkModeButton) {
        toggleDarkModeButton.addEventListener('click', toggleDarkMode);
    }
});

function toggleDarkMode() {
    document.body.classList.toggle('dark-mode');
    if (document.body.classList.contains('dark-mode')) {
        sessionStorage.setItem('darkMode', 'enabled'); // Changed from localStorage to sessionStorage
    } else {
        sessionStorage.setItem('darkMode', 'disabled'); // Changed from localStorage to sessionStorage
    }
}

function displayFeedback(message, isSuccess = true) {
    const output = document.getElementById('output');
    if (output) {
        output.innerText = message;
        output.style.color = isSuccess ? 'green' : 'red';
    } else {
        console.error("Element with ID 'output' not found in the DOM.");
    }
}

function runTerraformPlan() {
    const repo = sessionStorage.getItem('selectedRepo');
    const directory = sessionStorage.getItem('terraformDirectory');

    if (!repo || !directory) {
        console.error('selectedRepo or terraformDirectory is not set in sessionStorage.');
        displayFeedback('Please select a repository and specify a directory first.', false);
        return;
    }

    try {
        const payload = JSON.stringify({ repo, directory });

        callEndpoint('/terraform-plan', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: payload
        });
    } catch (error) {
        console.error('Error in runTerraformPlan:', error);
        displayFeedback(`Error: ${error.message}`, false);
    }
}

function runTerraformAction(endpoint) {
    const repo = sessionStorage.getItem('selectedRepo');
    const directory = sessionStorage.getItem('terraformDirectory');

    if (!repo || !directory) {
        console.error('selectedRepo or terraformDirectory is not set in sessionStorage.');
        displayFeedback('Please select a repository and specify a directory first.', false);
        return;
    }

    try {
        const payload = JSON.stringify({ repo, directory });

        callEndpoint(endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: payload
        });
    } catch (error) {
        console.error('Error in runTerraformAction:', error);
        displayFeedback(`Error: ${error.message}`, false);
    }
}

function callEndpoint(endpoint, options = {}) {
    showSpinner();
    displayFeedback('Loading...', true);

    return fetch(endpoint, options)
        .then(response => {
            hideSpinner();
            if (!response.ok) {
                throw new Error(`Network response was not ok: ${response.status} ${response.statusText}`);
            }
            return response.text();
        })
        .then(data => {
            displayFeedback(data, true);
        })
        .catch(error => {
            hideSpinner();
            displayFeedback(`Error: ${error.message}`, false);
        });
}

function showSpinner() {
    const spinner = document.getElementById('spinner');
    if (spinner) {
        spinner.style.display = 'block';
    } else {
        console.error("Element with ID 'spinner' not found in the DOM.");
    }
}

function hideSpinner() {
    const spinner = document.getElementById('spinner');
    if (spinner) {
        spinner.style.display = 'none';
    } else {
        console.error("Element with ID 'spinner' not found in the DOM.");
    }
}
