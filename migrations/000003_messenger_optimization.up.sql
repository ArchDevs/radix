PRAGMA foreign_keys=OFF;

-- Recreate users table with allow_search column
CREATE TABLE users_new (
    address TEXT PRIMARY KEY,
    public_key BLOB NOT NULL,
    username TEXT UNIQUE,
    display_name TEXT,
    allow_search BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users_new (address, public_key, username, display_name, created_at)
SELECT address, public_key, username, display_name, created_at FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

CREATE INDEX idx_users_username ON users(username);

-- Recreate challenges table with ON DELETE CASCADE
CREATE TABLE challenges_new (
    address TEXT PRIMARY KEY REFERENCES users(address) ON DELETE CASCADE,
    nonce TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO challenges_new (address, nonce, created_at)
SELECT address, nonce, created_at FROM challenges;

DROP TABLE challenges;
ALTER TABLE challenges_new RENAME TO challenges;

-- Recreate messages table with ON DELETE CASCADE and better indexes
CREATE TABLE messages_new (
    id TEXT PRIMARY KEY,
    sender TEXT REFERENCES users(address) ON DELETE CASCADE,
    recipient TEXT REFERENCES users(address) ON DELETE CASCADE,
    ciphertext BLOB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivered BOOLEAN DEFAULT FALSE
);

INSERT INTO messages_new (id, sender, recipient, ciphertext, created_at, delivered)
SELECT id, sender, recipient, ciphertext, created_at, delivered FROM messages;

DROP TABLE messages;
ALTER TABLE messages_new RENAME TO messages;

-- Optimized indexes for messenger queries
CREATE INDEX idx_messages_recipient_time ON messages(recipient, created_at);
CREATE INDEX idx_messages_recipient_delivered ON messages(recipient, delivered);
CREATE INDEX idx_messages_sender ON messages(sender, created_at);

PRAGMA foreign_keys=ON;