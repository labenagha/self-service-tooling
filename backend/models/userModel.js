// /backend/models/userModel.js

const db = require('../database/database');

const UserModel = {
    createUser: (username, hashedPassword, callback) => {
        const query = `INSERT INTO users (username, password) VALUES (?, ?)`;
        db.run(query, [username, hashedPassword], function(err) {
            callback(err, this.lastID);
        });
    },

    findUserByUsername: (username, callback) => {
        const query = `SELECT * FROM users WHERE username = ?`;
        db.get(query, [username], (err, row) => {
            callback(err, row);
        });
    },

    // Add more model methods as needed
};

module.exports = UserModel;
