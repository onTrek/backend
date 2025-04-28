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

	// Crea tabella users
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY, -- UUID
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TEXT NOT NULL -- ISO8601 format
	);
	`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Errore creazione tabella users: %v", err)
	}

	// Crea tabella activities
	createActivitiesTable := `
	CREATE TABLE IF NOT EXISTS activities (
		id INTEGER PRIMARY KEY AUTOINCREMENT, -- INTEGER ID autoincrement
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		start_time TEXT,
		end_time TEXT,
		created_at TEXT NOT NULL,
		km_percorsi FLOAT,
		dislivello_positivo FLOAT,
		dislivello_negativo FLOAT,
		altezza_partenza FLOAT,
		altezza_massima FLOAT,
		traccia TEXT, -- JSON array di punti GPS
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(createActivitiesTable)
	if err != nil {
		log.Fatalf("Errore creazione tabella activities: %v", err)
	}

	// Crea tabella gpx_files
	createGpxFilesTable := `
	CREATE TABLE IF NOT EXISTS gpx_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT, -- INTEGER ID autoincrement
		activity_id TEXT NOT NULL,
		filename TEXT NOT NULL,
		storage_path TEXT NOT NULL,
		upload_date TEXT NOT NULL,
		stats TEXT, -- JSON string
		FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(createGpxFilesTable)
	if err != nil {
		log.Fatalf("Errore creazione tabella gpx_files: %v", err)
	}

	fmt.Println("Database e tabelle create correttamente!")
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
