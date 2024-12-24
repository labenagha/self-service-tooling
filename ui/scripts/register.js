document.getElementById('registerForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const usernameInput = document.getElementById('username').value.trim();
    const passwordInput = document.getElementById('password').value.trim();
    const errorMessage = document.getElementById('error-message');

    errorMessage.textContent = '';

    if (!usernameInput || !passwordInput) {
        errorMessage.textContent = 'Please enter both username and password.';
        return;
    }

    try {
        const response = await fetch('http://localhost:5000/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username: usernameInput, password: passwordInput }),
        });

        const data = await response.json();

        if (response.status === 201) {
            alert('Account created successfully! You can now log in.');
            window.location.href = 'login.html';
        } else {
            errorMessage.textContent = data.message || 'An error occurred.';
        }
    } catch (error) {
        console.error('Error:', error);
        errorMessage.textContent = 'An error occurred while creating the account.';
    }
});

document.getElementById('returnLoginBtn').addEventListener('click', function () {
    window.location.href = 'login.html';
});
