// Function to clear saved settings and logout
function performLogout() {
    if (typeof(Storage) !== "undefined") {
        // Clear user information and repository settings from localStorage
        localStorage.removeItem('user');
        localStorage.removeItem('repoSettings');

        // Optionally, clear all localStorage
        // localStorage.clear();
    } else {
        console.warn('LocalStorage is not supported by this browser.');
    }
}

// Perform logout when the page loads
window.onload = performLogout;

// Function to redirect to Login page
function redirectToLogin() {
    window.location.href = 'login.html';
}

// Function to clear saved settings and logout
function performLogout() {
    if (typeof(Storage) !== "undefined") {
        // Clear user information and repository settings from localStorage
        localStorage.removeItem('user');
        localStorage.removeItem('repoSettings');
        localStorage.removeItem('token'); // Clear JWT token
    } else {
        console.warn('LocalStorage is not supported by this browser.');
    }
}

// Perform logout when the page loads
window.onload = performLogout;

// Function to redirect to Login page
function redirectToLogin() {
    window.location.href = 'login.html';
}

document.getElementById('createAccountBtn').addEventListener('click', function() {
    // Redirect to a registration page
    window.location.href = 'register.html';
});
