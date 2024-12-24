
document.getElementById('codeForm').addEventListener('submit', function(e) {
    e.preventDefault();

    const repoName = document.getElementById('repoName').value.trim();
    const folderName = document.getElementById('folderName').value.trim();
    const subFoldersInput = document.getElementById('subFolders').value.trim();
    const subFolders = subFoldersInput ? subFoldersInput.split(',').map(folder => folder.trim()) : [];

    // Create a settings object
    const settings = {
        repoName,
        folderName,
        subFolders
    };

    // Save settings to localStorage
    localStorage.setItem('repoSettings', JSON.stringify(settings));

    // Redirect back to the main page
    window.location.href = 'index.html';
});
