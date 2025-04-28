package db

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
)

func RegisterUser(db *sql.DB, user utils.User) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO users (id, email, password_hash, name, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(user.ID, user.Email, user.Password, user.Name, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
