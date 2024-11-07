document.addEventListener('DOMContentLoaded', () => {
    // Check and apply dark mode from localStorage
    if (localStorage.getItem('darkMode') === 'enabled') {
        document.body.classList.add('dark-mode');
    }

    // Event listeners for buttons
    document.getElementById('getCodeButton').addEventListener('click', () => {
        callEndpoint('/get-terraform-code');
    });

    document.getElementById('runPlanButton').addEventListener('click', () => {
        callEndpoint('/terraform-plan');
    });

    document.getElementById('deployButton').onclick = function() {
        fetch('/terraform-apply', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.text();
        })
        .then(data => {
            document.getElementById('output').innerText = data;
        })
        .catch(error => {
            console.error('Error:', error);
            document.getElementById('output').innerText = `Error: ${error.message}`;
        });
    };    

    document.getElementById('runDestroyButton').addEventListener('click', () => {
        callEndpoint('/terraform-destroy');
    });

    document.getElementById('toggleDarkModeButton').addEventListener('click', toggleDarkMode);
});

// Toggle dark mode and save preference to localStorage
function toggleDarkMode() {
    document.body.classList.toggle('dark-mode');
    if (document.body.classList.contains('dark-mode')) {
        localStorage.setItem('darkMode', 'enabled');
    } else {
        localStorage.setItem('darkMode', 'disabled');
    }
}

// Display feedback in the output area
function displayFeedback(message, isSuccess = true) {
    const output = document.getElementById('output');
    output.innerText = message;
    output.style.color = isSuccess ? 'green' : 'red';
}

// Call an endpoint and handle response
function callEndpoint(endpoint, options = {}) {
    displayFeedback('Loading...', true);
    return fetch(endpoint, options)
        .then(response => response.text())
        .then(data => {
            displayFeedback(data, true);
        })
        .catch(error => {
            displayFeedback(`Error: ${error.message}`, false);
        });
}
