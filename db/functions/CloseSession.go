package functions

import (
	"database/sql"
	"fmt"
	"time"
)

func CloseSession(db *sql.DB, sessionId int) error {
	// Close the session in the database
	_, err := db.Exec("UPDATE sessions SET closed_at = ? WHERE id = ?", time.Now().Format(time.RFC3339), sessionId)
	if err != nil {
		return fmt.Errorf("failed to close session: %v", err)
	}
	return nil
}
