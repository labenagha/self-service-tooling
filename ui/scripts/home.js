// For the "Last Successful Run" & "Last Failed Run" buttons
const getPlusButtons = document.querySelectorAll('.plus-btn, .second-plus-btn');

getPlusButtons.forEach(button => {
    button.addEventListener('click', function () {
        // Find the parent container of the button (either .box1 or .box2)
        const parentBox = this.closest('.box1, .box2');
        
        // Find the hidden box within this parent container
        const hiddenBox = parentBox.querySelector('.hidden-box1, .hidden-box2');
        
        if (hiddenBox) {
            hiddenBox.classList.toggle('show');
        } else {
            console.error('Hidden box not found for button:', this);
        }
    });
});

// Add click event to the Get Source Code button
document.querySelector('.get-plus-btn').addEventListener('click', () => {
    fetchRepositoryDetails();
});

// Function to display repository details
function displayRepositoryDetails(data) {
    // Update repository name
    document.getElementById('repo-name').textContent = `Repository Name: ${data.name}`;

    // Update folder list
    const folderList = document.getElementById('folder-list');
    folderList.innerHTML = ""; // Clear previous content
    data.folders.forEach(folder => {
        const li = document.createElement('li');
        li.textContent = folder;
        folderList.appendChild(li);
    });

    // Show the hidden-box with repository details
    document.querySelector('.hidden-box').style.display = 'block';
}

