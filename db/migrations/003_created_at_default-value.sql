CREATE TABLE users_new (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

INSERT INTO users_new (id, email, username, password_hash, created_at)
SELECT id, email, username, password_hash, created_at FROM users;

DROP TABLE users;

ALTER TABLE users_new RENAME TO users;

CREATE TABLE tokens_new (
    user_id TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    PRIMARY KEY (user_id, token),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO tokens_new (user_id, token, created_at)
SELECT user_id, token, created_at FROM tokens;

DROP TABLE tokens;

ALTER TABLE tokens_new RENAME TO tokens;

CREATE TABLE IF NOT EXISTS gpx_files_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    filename TEXT NOT NULL,
    storage_path TEXT NOT NULL UNIQUE,
    upload_date DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    title TEXT NOT NULL,
    km FLOAT NOT NULL DEFAULT 0,
    ascent INTEGER NOT NULL DEFAULT 0,
    descent INTEGER NOT NULL DEFAULT 0,
    duration TEXT NOT NULL DEFAULT 0,
    max_altitude INTEGER NOT NULL DEFAULT 0,
    min_altitude INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO gpx_files_new (id, user_id, filename, storage_path, upload_date, title, km, ascent, descent, duration, max_altitude, min_altitude)
SELECT id, user_id, filename, storage_path, upload_date, title, km, ascent, descent, duration, max_altitude, min_altitude FROM gpx_files;

DROP TABLE gpx_files;

ALTER TABLE gpx_files_new RENAME TO gpx_files;