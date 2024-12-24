// /backend/index.js

require('dotenv').config(); // Load environment variables first

const express = require('express');
const authRoutes = require('./routes/authRoutes');
const jwt = require('jsonwebtoken');
const cors = require('cors');
const bodyParser = require('body-parser');

const app = express();
const PORT = process.env.PORT || 5000;

// Middleware
app.use(cors());
app.use(bodyParser.json());

// Secret key for JWT
const JWT_SECRET = process.env.JWT_SECRET;

if (!JWT_SECRET) {
    console.error('FATAL ERROR: JWT_SECRET is not defined.');
    process.exit(1);
}

console.log('JWT_SECRET loaded:', JWT_SECRET);

// Mount Authentication Routes
app.use('/api', authRoutes);

// Protected Route Example
app.get('/api/protected', verifyToken, (req, res) => {
    jwt.verify(req.token, JWT_SECRET, (err, data) => {
        if (err) {
            return res.status(403).json({ message: 'Forbidden.' });
        }
        res.json({ message: 'This is protected data.', data });
    });
});

// Middleware to verify JWT Token
function verifyToken(req, res, next) {
    const bearerHeader = req.headers['authorization'];
    if (typeof bearerHeader !== 'undefined') {
        const bearer = bearerHeader.split(' ');
        const bearerToken = bearer[1];
        req.token = bearerToken;
        next();
    } else {
        res.status(403).json({ message: 'Forbidden.' });
    }
}

// Start the Server
app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
});
