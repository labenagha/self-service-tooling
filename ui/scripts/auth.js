// /frontend/scripts/auth.js

function checkAuthentication() {
    const token = localStorage.getItem('token');
    const currentPage = window.location.pathname.split('/').pop();

    if (currentPage === 'login.html') {
        if (token) {
            // User is already logged in, redirect to Home page
            window.location.href = 'index.html';
        }
        return;
    }

    if (!token) {
        // User is not logged in, redirect to login page
        window.location.href = 'login.html';
    }
}

checkAuthentication();
