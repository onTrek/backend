CREATE TABLE IF NOT EXISTS new_users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    extension VARCHAR(10) DEFAULT NULL
);

INSERT INTO new_users (id, email, password_hash, username, created_at)
SELECT id, email, password_hash, username, created_at FROM users;

DROP TABLE users;

ALTER TABLE new_users RENAME TO users;