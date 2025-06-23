-- Attiva vincoli di chiave esterna
PRAGMA foreign_keys = ON;

-- Crea tabella users
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TEXT NOT NULL
);

-- Crea tabella tokens
CREATE TABLE IF NOT EXISTS tokens (
    user_id TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL,
    PRIMARY KEY (user_id, token),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Crea tabella gpx_files
CREATE TABLE IF NOT EXISTS gpx_files (
     id INTEGER PRIMARY KEY AUTOINCREMENT,
     user_id TEXT NOT NULL,
     filename TEXT NOT NULL,
     storage_path TEXT NOT NULL,
     upload_date TEXT NOT NULL,
     title TEXT NOT NULL,
     km FLOAT NOT NULL DEFAULT 0,
     ascent FLOAT NOT NULL DEFAULT 0,
     descent FLOAT NOT NULL DEFAULT 0,
     duration INTEGER NOT NULL DEFAULT 0,
     max_altitude FLOAT NOT NULL DEFAULT 0,
     min_altitude FLOAT NOT NULL DEFAULT 0,
     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Crea tabella sessions
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    created_by TEXT NOT NULL,
    file_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    closed_at TEXT,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES gpx_files(id) ON DELETE SET NULL
);

-- Crea tabella session_members
CREATE TABLE IF NOT EXISTS session_members (
    session_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    latitude FLOAT DEFAULT -1,
    longitude FLOAT DEFAULT -1,
    altitude FLOAT DEFAULT -1,
    accuracy FLOAT DEFAULT -1,
    help_request BOOLEAN DEFAULT FALSE,
    going_to TEXT DEFAULT '',
    timestamp TEXT NOT NULL,
    color TEXT NOT NULL,
    PRIMARY KEY (session_id, user_id),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (going_to) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(session_id, color)
);

-- Crea tabella friends
CREATE TABLE IF NOT EXISTS friends (
    user_id1 TEXT NOT NULL,
    user_id2 TEXT NOT NULL,
    pending BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (user_id1, user_id2),
    FOREIGN KEY (user_id1) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id2) REFERENCES users(id) ON DELETE CASCADE
);
