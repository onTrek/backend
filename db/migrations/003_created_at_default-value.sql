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

CREATE TABLE IF NOT EXISTS groups_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    created_by TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    last_update DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    file_id INTEGER,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES gpx_files(id) ON DELETE SET NULL
);

INSERT INTO groups_new (id, description, created_by, file_id, created_at, last_update)
SELECT id, description, created_by, file_id, created_at, last_update FROM groups;

DROP TABLE groups;

ALTER TABLE groups_new RENAME TO groups;

CREATE TABLE IF NOT EXISTS group_members_new (
    group_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    latitude FLOAT DEFAULT -1,
    longitude FLOAT DEFAULT -1,
    altitude FLOAT DEFAULT -1,
    accuracy FLOAT DEFAULT -1,
    help_request BOOLEAN DEFAULT FALSE,
    going_to TEXT DEFAULT NULL,
    timestamp DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    color TEXT NOT NULL,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (going_to) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(group_id, color)
);

INSERT INTO group_members_new (group_id, user_id, latitude, longitude, altitude, accuracy, help_request, going_to, timestamp, color)
SELECT group_id, user_id, latitude, longitude, altitude, accuracy, help_request, going_to, timestamp, color FROM group_members;

DROP TABLE group_members;

ALTER TABLE group_members_new RENAME TO group_members;