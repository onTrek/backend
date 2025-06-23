package functions

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"time"
)

func GetUserByToken(db *sql.DB, token string) (utils.User, error) {
	var user utils.User
	var createdAt string

	err := db.QueryRow("SELECT users.id, users.email, users.username, tokens.created_at FROM users JOIN tokens ON users.id = tokens.user_id WHERE tokens.token = ?", token).Scan(&user.ID, &user.Email, &user.Username, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.User{}, fmt.Errorf("user not found")
		}
		return utils.User{}, fmt.Errorf("failed to query user: %w", err)
	}

	// Parse the created_at time
	parsedTime, err := time.Parse(time.RFC3339, createdAt)
	if time.Since(parsedTime) > tokenExpiry {
		return utils.User{}, fmt.Errorf("token expired")
	}

	return user, nil
}
