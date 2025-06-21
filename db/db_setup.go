package db

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func SetupDatabase() {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", ".root/db/ontrek.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	// Crea tabella users
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY, -- UUID
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		username TEXT NOT NULL,
		created_at TEXT NOT NULL -- ISO8601 format
	);
	`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Errore creazione tabella users: %v", err)
	}

	// Crea tabella tokens
	createTokensTable := `
	CREATE TABLE IF NOT EXISTS tokens (
		user_id TEXT NOT NULL, -- UUID user ID
		token TEXT NOT NULL UNIQUE, --UUID token
		created_at TEXT NOT NULL, -- ISO8601 format
		PRIMARY KEY (user_id, token),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(createTokensTable)
	if err != nil {
		log.Fatal("Errore creazione tabella tokens:", err)
	}

	// Crea tabella gpx_files
	createGpxFilesTable := `
	CREATE TABLE IF NOT EXISTS gpx_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT, -- INTEGER ID autoincrement
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
	`
	_, err = db.Exec(createGpxFilesTable)
	if err != nil {
		log.Fatalf("Errore creazione tabella gpx_files: %v", err)
	}

	// Crea tabella sessions
	createSessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		created_by TEXT NOT NULL,
		file_id INTEGER NOT NULL,
		created_at TEXT NOT NULL,
	    closed_at TEXT,
	    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE    
	    FOREIGN KEY (file_id) REFERENCES gpx_files(id) ON DELETE SET NULL
	);
	`
	_, err = db.Exec(createSessionsTable)
	if err != nil {
		log.Fatalf("Errore creazione tabella sessions: %v", err)
	}

	// Crea tabella membri
	createMembriTable := `
	CREATE TABLE IF NOT EXISTS session_members (
		session_id INTEGER NOT NULL,
		user_id TEXT NOT NULL,
		latitude FLOAT,
		longitude FLOAT,
		altitude FLOAT,
		accuracy FLOAT,
		help_request BOOLEAN DEFAULT FALSE,
		going_to TEXT,
		timestamp TEXT NOT NULL,
		PRIMARY KEY (session_id, user_id),
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	    FOREIGN KEY (going_to) REFERENCES users(id) ON DELETE SET NULL
	);
	`
	_, err = db.Exec(createMembriTable)
	if err != nil {
		log.Fatalf("Errore creazione tabella membri: %v", err)
	}

	// Crea tabella amici
	createAmiciTable := `
	CREATE TABLE IF NOT EXISTS friends (
		user_id1 TEXT NOT NULL,
		user_id2 TEXT NOT NULL,
		PRIMARY KEY (user_id1, user_id2),
		FOREIGN KEY (user_id1) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id2) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(createAmiciTable)
	if err != nil {
		log.Fatalf("Errore creazione tabella amici: %v", err)
	}

	fmt.Println("Database e tabelle create correttamente!")

	// Attiva i vincoli di chiave esterna
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Impossibile attivare PRAGMA foreign_keys:", err)
	}
}

func DatabaseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Open a connection to the SQLite database
		db, err := sql.Open("sqlite3", "./ontrek.db")
		if err != nil {
			log.Fatal(err)
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(db)

		c.Set("db", db)
		c.Next()
	}
}
