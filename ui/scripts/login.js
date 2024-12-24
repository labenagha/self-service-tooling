document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('token');

    // If token exists, redirect to home page
    if (token) {
        window.location.href = 'index.html';
    }
});

document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const usernameInput = document.getElementById('username').value.trim();
    const passwordInput = document.getElementById('password').value.trim();
    const errorMessage = document.getElementById('error-message');

    // Log the input values for debugging
    console.log('Username:', usernameInput);
    console.log('Password:', passwordInput);

    // Clear previous error message
    errorMessage.textContent = '';

    // Basic validation
    if (!usernameInput || !passwordInput) {
        errorMessage.textContent = 'Please enter both username and password.';
        return;
    }

    try {
        const response = await fetch('http://localhost:5000/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username: usernameInput, password: passwordInput })
        });

        const data = await response.json();

        // Log the response status and data for debugging
        console.log('Response Status:', response.status);
        console.log('Response Data:', data);

        if (response.status === 200) {
            // Save JWT token in localStorage
            localStorage.setItem('token', data.token);

            // Optionally, store user information
            const user = {
                username: usernameInput,
                token: data.token
            };
            localStorage.setItem('user', JSON.stringify(user));

            // Redirect to Home page
            window.location.href = 'index.html';
        } else if (response.status === 400 && data.errors) {
            // Display validation errors from express-validator
            const messages = data.errors.map(error => error.msg).join(' ');
            errorMessage.textContent = messages;
        } else {
            // Display generic error message
            errorMessage.textContent = data.message || 'An error occurred.';
        }
    } catch (error) {
        console.error('Error:', error);
        errorMessage.textContent = 'An error occurred while logging in.';
    }
});
