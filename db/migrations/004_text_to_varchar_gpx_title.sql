CREATE TABLE IF NOT EXISTS gpx_files_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    filename TEXT NOT NULL,
    storage_path TEXT NOT NULL UNIQUE,
    upload_date DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    title VARCHAR(64) NOT NULL,
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

CREATE TABLE IF NOT EXISTS friends_new (
    user_id1 TEXT NOT NULL,
    user_id2 TEXT NOT NULL,
    pending BOOLEAN DEFAULT TRUE,
    date DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    PRIMARY KEY (user_id1, user_id2),
    FOREIGN KEY (user_id1) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id2) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO friends_new (user_id1, user_id2, pending, date)
SELECT user_id1, user_id2, pending, strftime('%Y-%m-%dT%H:%M:%SZ', 'now') FROM friends;

DROP TABLE friends;

ALTER TABLE friends_new RENAME TO friends;

CREATE TABLE IF NOT EXISTS groups_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description VARCHAR(64) NOT NULL,
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