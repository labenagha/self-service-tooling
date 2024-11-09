document.addEventListener('DOMContentLoaded', () => {
    fetch('/api/list-repos')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            const repoList = document.getElementById('repo-list');
            if (Array.isArray(data)) {
                data.forEach(repo => {
                    const repoItem = document.createElement('div');
                    repoItem.classList.add('card', 'm-2', 'p-3');
                    repoItem.style.width = '300px';
                    repoItem.innerHTML = `
                        <h5 class="card-title">${repo.full_name}</h5>
                        <p class="card-text">${repo.description ? `Description: ${repo.description}` : 'No description available'}</p>
                        <a href="${repo.html_url}" target="_blank" class="btn btn-primary mb-2">View on GitHub</a>
                        <button class="btn btn-success select-repo" data-repo="${repo.full_name}">Select</button>
                    `;
                    repoList.appendChild(repoItem);
                });

                // Handle repository selection
                document.querySelectorAll('.select-repo').forEach(button => {
                    button.addEventListener('click', (event) => {
                        const selectedRepo = event.target.getAttribute('data-repo');
                        sessionStorage.setItem('selectedRepo', selectedRepo); // Changed to sessionStorage
                        
                        // Prompt for Terraform directory
                        const directoryInput = document.createElement('input');
                        directoryInput.placeholder = 'Enter Terraform directory';
                        directoryInput.classList.add('form-control', 'mt-2');
                        const confirmButton = document.createElement('button');
                        confirmButton.textContent = 'Confirm Directory';
                        confirmButton.classList.add('btn', 'btn-primary', 'mt-2');

                        // Append input and button to repo item
                        repoList.innerHTML = '';  // Clear repo list to show only directory input
                        repoList.appendChild(directoryInput);
                        repoList.appendChild(confirmButton);

                        // Handle directory confirmation
                        confirmButton.addEventListener('click', () => {
                            const terraformDirectory = directoryInput.value;
                            if (terraformDirectory) {
                                sessionStorage.setItem('terraformDirectory', terraformDirectory); // Changed to sessionStorage
                                alert(`Directory "${terraformDirectory}" selected for ${selectedRepo}`);
                                window.location.href = '/'; // Redirect to main page
                            } else {
                                alert('Please enter a valid directory.');
                            }
                        });
                    });
                });
            } else {
                repoList.innerText = 'Failed to fetch repositories.';
            }
        })
        .catch(error => {
            console.error('Error fetching repositories:', error);
            document.getElementById('repo-list').innerText = `An error occurred: ${error.message}`;
        });
});
