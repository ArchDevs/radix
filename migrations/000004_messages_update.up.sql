PRAGMA foreign_keys=OFF;

-- Recreate messages table with content TEXT and read column
CREATE TABLE messages_new (
    id TEXT PRIMARY KEY,
    sender TEXT REFERENCES users(address) ON DELETE CASCADE,
    recipient TEXT REFERENCES users(address) ON DELETE CASCADE,
    content TEXT NOT NULL, -- encrypted base64
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivered BOOLEAN DEFAULT FALSE,
    read BOOLEAN DEFAULT FALSE
);

INSERT INTO messages_new (id, sender, recipient, content, created_at, delivered, read)
SELECT id, sender, recipient, hex(ciphertext), created_at, delivered, FALSE FROM messages;

DROP TABLE messages;
ALTER TABLE messages_new RENAME TO messages;

-- SQLite-optimized indexes for common query patterns
CREATE INDEX idx_messages_recipient_created ON messages(recipient, created_at DESC);
CREATE INDEX idx_messages_sender_created ON messages(sender, created_at DESC);
CREATE INDEX idx_messages_recipient_read ON messages(recipient, read) WHERE read = FALSE;
CREATE INDEX idx_messages_recipient_delivered ON messages(recipient, delivered);

PRAGMA foreign_keys=ON;
