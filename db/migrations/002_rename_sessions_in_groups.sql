CREATE TABLE groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    created_by TEXT NOT NULL,
    created_at TEXT NOT NULL,
    last_update TEXT NOT NULL,
    file_id INTEGER,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES gpx_files(id) ON DELETE SET NULL
);

INSERT INTO groups (
    id, description, created_by, file_id, created_at, last_update
)
SELECT
    id, description, created_by, file_id, created_at, ''
FROM sessions;

CREATE TABLE group_members (
    group_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    latitude FLOAT DEFAULT -1,
    longitude FLOAT DEFAULT -1,
    altitude FLOAT DEFAULT -1,
    accuracy FLOAT DEFAULT -1,
    help_request BOOLEAN DEFAULT FALSE,
    going_to TEXT DEFAULT NULL,
    timestamp TEXT NOT NULL,
    color TEXT NOT NULL,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (going_to) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(group_id, color)
);

INSERT INTO group_members (
    group_id, user_id, latitude, longitude, altitude, accuracy,
    help_request, going_to, timestamp, color
)
SELECT
    session_id, user_id, latitude, longitude, altitude, accuracy,
    help_request, going_to, timestamp, color
FROM session_members;

DROP TABLE session_members;

DROP TABLE sessions;
