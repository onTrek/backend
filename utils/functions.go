package utils

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func IsLogged(c *gin.Context, token string) (User, error) {
	db := c.MustGet("db").(*sql.DB)
	var user User
	err := db.QueryRow("SELECT id, email, name FROM users WHERE id = ?", token).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("failed to query user: %w", err)
	}
	return user, nil
}

func DeleteFile(path any) error {
	err := os.Remove(path.(string))
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
