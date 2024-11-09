document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const repoName = urlParams.get('repo');

    if (!repoName) {
        alert('Repository not specified.');
        window.location.href = '/list-repos.html';
        return;
    }

    // Display the selected repository name in the UI
    document.getElementById('repo-name').textContent = repoName;

    const dirSelect = document.getElementById('directory');
    const messageDiv = document.getElementById('message');

    // Fetch the GitHub access token from sessionStorage
    const token = sessionStorage.getItem('github_access_token');

    if (!token) {
        messageDiv.innerHTML = `<div class="alert alert-danger">GitHub authentication token is missing.</div>`;
        return;
    }

    // Fetch directories for the selected repository
    fetch(`/api/directories?repo=${encodeURIComponent(repoName)}`, {
        method: 'GET',
        headers: {
            'Authorization': `token ${token}`,
            'Content-Type': 'application/json'
        }
    })
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

        // Populate the directories dropdown
        data.directories.forEach(dir => {
            const option = document.createElement('option');
            option.value = dir.path; // Assuming `dir.path` is the directory path
            option.textContent = dir.name; // Assuming `dir.name` is the directory name
            dirSelect.appendChild(option);
        });
        dirSelect.disabled = false;
    })
    .catch(error => {
        console.error('Error fetching directories:', error);
        messageDiv.innerHTML = `<div class="alert alert-danger">Failed to load directories.</div>`;
        dirSelect.disabled = true;
    });

    // Form submission for selecting a directory
    const form = document.getElementById('directory-form');
    form.addEventListener('submit', event => {
        event.preventDefault();

        const directory = dirSelect.value;
        if (!directory) {
            messageDiv.innerHTML = `<div class="alert alert-warning">Please select a directory.</div>`;
            return;
        }

        // Store the selected repository and directory in sessionStorage
        sessionStorage.setItem('selectedRepo', repoName);
        sessionStorage.setItem('terraformDirectory', directory);

        // Redirect to the operations page
        window.location.href = '/operations';
    });
});
