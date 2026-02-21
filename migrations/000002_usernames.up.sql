PRAGMA foreign_keys=OFF;

CREATE TABLE users_new (
    address TEXT PRIMARY KEY,
    public_key BLOB NOT NULL,
    username TEXT UNIQUE,
    display_name TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users_new (address, public_key, created_at)
SELECT address, public_key, created_at FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

CREATE INDEX idx_users_username ON users(username);

PRAGMA foreign_keys=ON;