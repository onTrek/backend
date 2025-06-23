package db

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabase() {
	db, err := sql.Open("sqlite3", "./db/ontrek.db")
	if err != nil {
		log.Fatal("Errore apertura DB:", err)
	}
	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Errore attivazione foreign_keys:", err)
	}

	requiredTables := []string{
		"users",
		"tokens",
		"gpx_files",
		"groups",
		"group_members",
		"friends",
	}

	for _, tableName := range requiredTables {
		var name string
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
		err := db.QueryRow(query, tableName).Scan(&name)
		if err == sql.ErrNoRows {
			log.Fatalf("Tabella mancante: %s", tableName)
		} else if err != nil {
			log.Fatalf("Errore verifica tabella %s: %v", tableName, err)
		} else {
			log.Printf("Tabella presente: %s", tableName)
		}
	}

	fmt.Println("Tutte le tabelle richieste sono presenti nel database.")
}

func DatabaseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Open a connection to the SQLite database
		db, err := sql.Open("sqlite3", "./db/ontrek.db")
		if err != nil {
			log.Fatal(err)
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(db)

		_, err = db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			log.Fatal("Impossibile attivare PRAGMA foreign_keys:", err)
		}

		c.Set("db", db)
		c.Next()
	}
}
