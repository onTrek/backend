package db

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func UpdateToken(tx *sql.Tx, userID string) (utils.UserToken, error) {

	var tokenStruct utils.UserToken

	// Generate token for the user
	token := uuid.New().String()

	// Delete any existing token for the user
	_, err := tx.Exec("DELETE FROM tokens WHERE user_id = ?", userID)
	if err != nil {
		return utils.UserToken{}, fmt.Errorf("failed to delete existing token: %w", err)
	}

	// Store the token in the database (assuming a tokens table exists)
	stmt, err := tx.Prepare("INSERT INTO tokens (user_id, token, created_at) VALUES (?, ?, ?)")
	if err != nil {
		return utils.UserToken{}, fmt.Errorf("failed to prepare token statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, token, time.Now().Format(time.RFC3339))
	if err != nil {
		return utils.UserToken{}, fmt.Errorf("failed to execute token statement: %w", err)
	}

	tokenStruct.Token = token

	return tokenStruct, nil
}
