PRAGMA foreign_keys=OFF;

-- Restore messages table
CREATE TABLE messages_new (
    id TEXT PRIMARY KEY,
    sender TEXT REFERENCES users(address),
    recipient TEXT REFERENCES users(address),
    ciphertext BLOB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivered BOOLEAN DEFAULT FALSE
);

INSERT INTO messages_new (id, sender, recipient, ciphertext, created_at, delivered)
SELECT id, sender, recipient, ciphertext, created_at, delivered FROM messages;

DROP TABLE messages;
ALTER TABLE messages_new RENAME TO messages;

CREATE INDEX idx_messages_recipient ON messages(recipient, created_at);

-- Restore challenges table
CREATE TABLE challenges_new (
    address TEXT PRIMARY KEY REFERENCES users(address),
    nonce TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO challenges_new (address, nonce, created_at)
SELECT address, nonce, created_at FROM challenges;

DROP TABLE challenges;
ALTER TABLE challenges_new RENAME TO challenges;

-- Restore users table
CREATE TABLE users_new (
    address TEXT PRIMARY KEY,
    public_key BLOB NOT NULL,
    username TEXT UNIQUE,
    display_name TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users_new (address, public_key, username, display_name, created_at)
SELECT address, public_key, username, display_name, created_at FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

CREATE INDEX idx_users_username ON users(username);

PRAGMA foreign_keys=ON;