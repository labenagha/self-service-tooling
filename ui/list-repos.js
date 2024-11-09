// File: ui/list-repos.js

// JavaScript specific to the list-repos page
console.log("List Repositories JavaScript loaded.");

// Fetch and display repositories
function fetchAndDisplayRepos() {
    fetch('/list-repos-data')
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to fetch repositories.');
        }
        return response.json();
    })
    .then(data => {
        const reposList = document.getElementById('reposList');
        if (reposList) {
            reposList.innerHTML = ''; // Clear any existing content
            if (data.length === 0) {
                reposList.innerHTML = `<li class="list-group-item">No repositories found.</li>`;
                return;
            }
            data.forEach(repo => {
                const listItem = document.createElement('li');
                listItem.className = 'list-group-item repo-item';
                listItem.textContent = repo.FullName;

                const selectButton = document.createElement('a');
                selectButton.href = `/select-directory?repo=${encodeURIComponent(repo.FullName)}`;
                selectButton.className = 'btn btn-success btn-sm';
                selectButton.innerHTML = '<i class="fas fa-check"></i> Select';

                listItem.appendChild(selectButton);
                reposList.appendChild(listItem);
            });
        }
    })
    .catch(error => {
        console.error('Error fetching repositories:', error);
        const reposList = document.getElementById('reposList');
        if (reposList) {
            reposList.innerHTML = `<li class="list-group-item text-danger">Failed to load repositories.</li>`;
        }
    });
}

document.addEventListener('DOMContentLoaded', fetchAndDisplayRepos);
