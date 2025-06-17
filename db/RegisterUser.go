package db

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
)

func RegisterUser(db *sql.DB, user utils.User) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Rollback on error unless committed
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Prepare the SQL statement
	stmt, err := tx.Prepare("INSERT INTO users (id, email, password_hash, username, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(user.ID, user.Email, user.Password, user.Username, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	// Insert tokens table
	_, err = UpdateToken(tx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
