// File: ui/scripts/get-code.js

document.addEventListener('DOMContentLoaded', () => {
    const repoSelect = document.getElementById('repoName');

    // Function to fetch repositories from the server
    const fetchRepos = async () => {
        try {
            console.log('Fetching repositories from /api/repos');
            const response = await fetch('/api/repos', {
                method: 'GET',
                credentials: 'include' // Include cookies for authentication
            });

            console.log(`Received response: ${response.status}`);

            if (!response.ok) {
                throw new Error(`Failed to fetch repositories: ${response.statusText}`);
            }

            const repos = await response.json();
            console.log(`Fetched ${repos.length} repositories`);

            // Populate the select dropdown with repositories
            repos.forEach(repo => {
                const option = document.createElement('option');
                option.value = repo.name;
                option.textContent = repo.name;
                repoSelect.appendChild(option);
            });

            if (repos.length === 0) {
                const option = document.createElement('option');
                option.value = "";
                option.textContent = "No repositories found";
                repoSelect.appendChild(option);
                repoSelect.disabled = true;
            }
        } catch (error) {
            console.error('Error fetching repositories:', error);
            const option = document.createElement('option');
            option.value = "";
            option.textContent = "Error fetching repositories";
            repoSelect.appendChild(option);
            repoSelect.disabled = true;
        }
    };

    fetchRepos();

    // Optional: Handle form submission
    const form = document.getElementById('codeForm');
    form.addEventListener('submit', (e) => {
        // Add any client-side validation or processing here
        // For now, it will submit normally
    });
});
