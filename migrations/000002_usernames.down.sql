-- Reverse the migration
PRAGMA foreign_keys=OFF;

CREATE TABLE users_new (
    address TEXT PRIMARY KEY,
    public_key BLOB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users_new (address, public_key, created_at)
SELECT address, public_key, created_at FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

PRAGMA foreign_keys=ON;