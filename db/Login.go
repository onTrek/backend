package db

import (
	"OnTrek/utils"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var tokenExpiry = 7 * 24 * time.Hour // Token expiry duration

func Login(db *sql.DB, email string, password string) (utils.UserToken, error) {

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return utils.UserToken{}, err
	}

	// Rollback on error unless committed
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var createdAt string
	var user utils.UserToken
	var userID string

	// Check if the user exists
	var hashedPassword string
	err = tx.QueryRow("SELECT password_hash FROM users WHERE email = ?", email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.UserToken{}, fmt.Errorf("user not found")
		}
		return utils.UserToken{}, fmt.Errorf("failed to query user: %w", err)
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return utils.UserToken{}, fmt.Errorf("invalid password: %w", err)
	}

	err = tx.QueryRow("SELECT users.id, tokens.token, tokens.created_at FROM tokens JOIN users ON users.id = tokens.user_id WHERE users.email = ?", email).Scan(&userID, &user.Token, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.UserToken{}, fmt.Errorf("user not found")
		}
		return utils.UserToken{}, fmt.Errorf("failed to query user ID: %w", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return utils.UserToken{}, fmt.Errorf("failed to parse created_at time: %w", err)
	}

	if time.Since(parsedTime) > tokenExpiry {
		user, err = UpdateToken(tx, userID)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return utils.UserToken{}, err
	}

	return user, nil
}
