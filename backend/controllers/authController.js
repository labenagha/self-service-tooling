const { validationResult } = require('express-validator');
const UserModel = require('../models/userModel');
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');

const JWT_SECRET = process.env.JWT_SECRET;

if (!JWT_SECRET) {
    console.error('FATAL ERROR: JWT_SECRET is not defined.');
    process.exit(1);
}

const authController = {
    register: (req, res) => {
        const errors = validationResult(req);
        if (!errors.isEmpty()) {
            return res.status(400).json({ errors: errors.array() });
        }

        const { username, password } = req.body;

        // Check if user already exists
        UserModel.findUserByUsername(username, (err, user) => {
            if (err) {
                console.error('Database error during user lookup:', err);
                return res.status(500).json({ message: 'Internal server error.' });
            }

            if (user) {
                return res.status(400).json({ message: 'Username already exists.' });
            }

            // Hash the password
            bcrypt.hash(password, 10, (err, hash) => {
                if (err) {
                    console.error('Error hashing password:', err);
                    return res.status(500).json({ message: 'Internal server error.' });
                }

                // Insert the new user into the database
                UserModel.createUser(username, hash, (err, userId) => {
                    if (err) {
                        console.error('Database insert error:', err);
                        return res.status(500).json({ message: 'Internal server error.' });
                    }

                    console.log(`User ${username} registered successfully with ID ${userId}.`);
                    return res.status(201).json({ message: 'User registered successfully.', userId });
                });
            });
        });
    },

    login: (req, res) => {
        const errors = validationResult(req);
        if (!errors.isEmpty()) {
            return res.status(400).json({ errors: errors.array() });
        }

        const { username, password } = req.body;

        // Retrieve user from the database
        UserModel.findUserByUsername(username, (err, user) => {
            if (err) {
                console.error('Database error during user lookup:', err);
                return res.status(500).json({ message: 'Internal server error.' });
            }

            if (!user) {
                return res.status(400).json({ message: 'Invalid username or password.' });
            }

            // Compare passwords
            bcrypt.compare(password, user.password, (err, isMatch) => {
                if (err) {
                    console.error('Error comparing passwords:', err);
                    return res.status(500).json({ message: 'Internal server error.' });
                }

                if (!isMatch) {
                    return res.status(400).json({ message: 'Invalid username or password.' });
                }

                // Create JWT payload
                const payload = {
                    id: user.id,
                    username: user.username
                };

                // Sign token
                jwt.sign(payload, JWT_SECRET, { expiresIn: '1h' }, (err, token) => {
                    if (err) {
                        console.error('Error signing JWT:', err);
                        return res.status(500).json({ message: 'Internal server error.' });
                    }

                    console.log(`User ${username} logged in successfully.`);
                    return res.status(200).json({ message: 'Login successful.', token });
                });
            });
        });
    },
};

module.exports = authController;
