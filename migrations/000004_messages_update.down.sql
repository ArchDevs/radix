PRAGMA foreign_keys=OFF;

-- Revert messages table to previous schema
CREATE TABLE messages_old (
    id TEXT PRIMARY KEY,
    sender TEXT REFERENCES users(address) ON DELETE CASCADE,
    recipient TEXT REFERENCES users(address) ON DELETE CASCADE,
    ciphertext BLOB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivered BOOLEAN DEFAULT FALSE
);

INSERT INTO messages_old (id, sender, recipient, ciphertext, created_at, delivered)
SELECT id, sender, recipient, unhex(content), created_at, delivered FROM messages;

DROP TABLE messages;
ALTER TABLE messages_old RENAME TO messages;

-- Restore previous indexes
CREATE INDEX idx_messages_recipient_time ON messages(recipient, created_at);
CREATE INDEX idx_messages_recipient_delivered ON messages(recipient, delivered);
CREATE INDEX idx_messages_sender ON messages(sender, created_at);

PRAGMA foreign_keys=ON;
