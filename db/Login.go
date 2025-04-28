package db

import (
	"OnTrek/utils"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB, user utils.User) (utils.User, error) {
	// Check if the user exists
	var hashedPassword string
	err := db.QueryRow("SELECT password_hash FROM users WHERE email = ?", user.Email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.User{}, fmt.Errorf("user not found")
		}
		return utils.User{}, fmt.Errorf("failed to query user: %w", err)
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return utils.User{}, fmt.Errorf("invalid password: %w", err)
	}

	err = db.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.User{}, fmt.Errorf("user not found")
		}
		return utils.User{}, fmt.Errorf("failed to query user ID: %w", err)
	}

	return user, nil
}
