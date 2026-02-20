CREATE TABLE users (
    address TEXT PRIMARY KEY,
    public_key BLOB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE challenges (
    address TEXT PRIMARY KEY REFERENCES users(address),
    nonce TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    sender TEXT REFERENCES users(address),
    recipient TEXT REFERENCES users(address),
    ciphertext BLOB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivered BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_messages_recipient ON messages(recipient, created_at);