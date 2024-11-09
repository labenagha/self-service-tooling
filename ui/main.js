// File: ui/main.js

document.addEventListener('DOMContentLoaded', () => {
    // Handle Registration Form
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        const messageDiv = document.getElementById('message');

        registerForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const username = document.getElementById('regUsername').value.trim();
            const password = document.getElementById('regPassword').value.trim();

            fetch('/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            })
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => { throw new Error(text) });
                }
                return response.text();
            })
            .then(data => {
                messageDiv.innerHTML = `<div class="alert alert-success">${data}</div>`;
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
            })
            .catch(error => {
                messageDiv.innerHTML = `<div class="alert alert-danger">${error.message}</div>`;
            });
        });
    }

    // Handle Login Form
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        const messageDiv = document.getElementById('message');

        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const username = document.getElementById('username').value.trim();
            const password = document.getElementById('password').value.trim();

            fetch('/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            })
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => { throw new Error(text) });
                }
                return response.text();
            })
            .then(data => {
                messageDiv.innerHTML = `<div class="alert alert-success">${data}</div>`;
                setTimeout(() => {
                    window.location.href = '/operations';
                }, 2000);
            })
            .catch(error => {
                messageDiv.innerHTML = `<div class="alert alert-danger">${error.message}</div>`;
            });
        });
    }

    // Handle Repository and Directory Selection Form
    const selectionForm = document.getElementById('selectionForm');
    if (selectionForm) {
        const repoSelect = document.getElementById('repoSelect');
        const dirSelect = document.getElementById('dirSelect');
        const messageDiv = document.getElementById('message');

        repoSelect.addEventListener('change', () => {
            const selectedRepo = repoSelect.value;
            dirSelect.innerHTML = '<option value="">-- Select Directory --</option>';
            if (selectedRepo) {
                fetch(`/api/directories?repo=${encodeURIComponent(selectedRepo)}`)
                    .then(response => {
                        if (!response.ok) {
                            throw new Error('Failed to fetch directories.');
                        }
                        return response.json();
                    })
                    .then(data => {
                        if (!data.directories || !Array.isArray(data.directories)) {
                            throw new Error('Invalid directories data.');
                        }
                        data.directories.forEach(dir => {
                            const option = document.createElement('option');
                            option.value = dir.path; // Corrected from dir.Path to dir.path
                            option.textContent = dir.name; // Corrected from dir.Name to dir.name
                            dirSelect.appendChild(option);
                        });
                        dirSelect.disabled = false;
                    })
                    .catch(error => {
                        console.error('Error fetching directories:', error);
                        messageDiv.innerHTML = `<div class="alert alert-danger">Failed to load directories.</div>`;
                        dirSelect.disabled = true;
                    });
            } else {
                dirSelect.disabled = true;
            }
        });

        selectionForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const selectedRepo = repoSelect.value;
            const selectedDir = dirSelect.value;

            if (!selectedRepo || !selectedDir) {
                messageDiv.innerHTML = `<div class="alert alert-warning">Please select both repository and directory.</div>`;
                return;
            }

            fetch('/list-repos', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ repo: selectedRepo, directory: selectedDir })
            })
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => { throw new Error(text) });
                }
                return response.text();
            })
            .then(data => {
                messageDiv.innerHTML = `<div class="alert alert-success">${data}</div>`;
                setTimeout(() => {
                    window.location.href = '/operations';
                }, 2000);
            })
            .catch(error => {
                messageDiv.innerHTML = `<div class="alert alert-danger">${error.message}</div>`;
            });
        });
    }

    // Handle Terraform Actions in Operations Dashboard
    const terraformPlanButton = document.getElementById('runPlanButton');
    const terraformApplyButton = document.getElementById('deployButton');
    const terraformDestroyButton = document.getElementById('runDestroyButton');
    const getCodeButton = document.getElementById('getCodeButton');
    const outputDiv = document.getElementById('output');

    const handleTerraformAction = (action) => {
        // Fetch selected repo and directory from DOM
        const repo = document.getElementById('selectedRepo').textContent.trim();
        const directory = document.getElementById('selectedDirectory').textContent.trim();

        if (!repo || !directory || repo === 'Not Selected' || directory === 'Not Selected') {
            if (outputDiv) {
                outputDiv.innerHTML = `<div class="alert alert-warning">Please select a repository and directory first.</div>`;
            }
            return;
        }

        // Display loading message
        if (outputDiv) {
            outputDiv.innerHTML = `<div class="alert alert-info">Executing Terraform ${action}...</div>`;
        }

        // Send POST request to perform Terraform action
        fetch(`/terraform-${action}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ repo, directory })
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.text();
        })
        .then(data => {
            if (outputDiv) {
                outputDiv.innerHTML = `<pre>${escapeHtml(data)}</pre>`;
            }
        })
        .catch(error => {
            if (outputDiv) {
                outputDiv.innerHTML = `<div class="alert alert-danger">${escapeHtml(error.message)}</div>`;
            }
        });
    };

    // Function to escape HTML to prevent XSS
    function escapeHtml(text) {
        var map = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#039;'
        };
        return text.replace(/[&<>"']/g, function(m) { return map[m]; });
    }

    if (terraformPlanButton) {
        terraformPlanButton.addEventListener('click', (e) => {
            e.preventDefault();
            handleTerraformAction('plan');
        });
    }

    if (terraformApplyButton) {
        terraformApplyButton.addEventListener('click', (e) => {
            e.preventDefault();
            handleTerraformAction('apply');
        });
    }

    if (terraformDestroyButton) {
        terraformDestroyButton.addEventListener('click', (e) => {
            e.preventDefault();
            handleTerraformAction('destroy');
        });
    }

    if (getCodeButton) {
        getCodeButton.addEventListener('click', (e) => {
            e.preventDefault();
            handleTerraformAction('getcode');
        });
    }

    // Handle GitHub OAuth Login Button (if present)
    const githubLoginButton = document.getElementById('githubLoginButton');
    if (githubLoginButton) {
        githubLoginButton.addEventListener('click', (e) => {
            e.preventDefault();
            window.location.href = '/auth/github/login';
        });
    }
});
