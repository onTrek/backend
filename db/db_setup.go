package db

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func SetupDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("./db/ontrek.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Errore apertura DB:", err)
	}

	if err := DB.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
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
		var count int64
		err := DB.Raw("SELECT count(name) FROM sqlite_master WHERE type='table' AND name = ?", tableName).Scan(&count).Error
		if err != nil {
			log.Fatalf("Error verifying table %s: %v", tableName, err)
		}
		if count == 0 {
			log.Fatalf("Missing table: %s", tableName)
		} else {
			log.Printf("Table %s is present in the database.", tableName)
		}
	}

	fmt.Println("All required tables are present in the database.")
}

func DatabaseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if DB == nil {
			log.Fatal("Database not initialized")
		}

		c.Set("db", DB)
		c.Next()
	}
}
