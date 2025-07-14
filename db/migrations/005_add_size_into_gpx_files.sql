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
    size INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO gpx_files_new (id, user_id, filename, storage_path, upload_date, title, km, ascent, descent, duration, max_altitude, min_altitude, size)
SELECT id, user_id, filename, storage_path, upload_date, title, km, ascent, descent, duration, max_altitude, min_altitude, 0 FROM gpx_files;

DROP TABLE gpx_files;

ALTER TABLE gpx_files_new RENAME TO gpx_files;