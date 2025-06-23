package functions

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
)

func GetUserById(db *sql.DB, userId string) (utils.User, error) {
	var user utils.User
	err := db.QueryRow("SELECT id, email, username FROM users WHERE id = ?", userId).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.User{}, fmt.Errorf("user not found")
		}
		return utils.User{}, fmt.Errorf("failed to query user: %w", err)
	}
	return user, nil
}
