package models

import (
	"OnTrek/utils"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"type:uuid;primaryKey"`
	Email        string    `json:"email" example:"user@example.com" gorm:"unique;not null"`
	Username     string    `json:"username" example:"John Doe" gorm:"not null"`
	PasswordHash string    `json:"password_hash" example:"strongPassword123" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" example:"2025-05-11T08:00:00Z" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Extension    *string   `json:"extension" example:".png" gorm:"default:null"`
}

var tokenExpiry = 365 * 24 * time.Hour

func RegisterUser(db *gorm.DB, user User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := UpdateToken(tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func DeleteUser(db *gorm.DB, userID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return fmt.Errorf("error enabling foreign key enforcement: %v", err)
		}

		if err := tx.Where("id = ?", userID).Delete(&User{}).Error; err != nil {
			return fmt.Errorf("failed to delete user token: %w", err)
		}

		filepath, err := utils.FindFileByID(userID)
		if err != nil {
			return fmt.Errorf("failed to find user file: %w", err)
		}
		fmt.Println("File path to delete:", filepath)

		if filepath == "" {
			return nil
		} else {
			if _, err := os.Stat(filepath); err == nil {
				err := os.Remove(filepath)
				if err != nil {
					return fmt.Errorf("failed to delete file: %w", err)
				}
			} else if os.IsNotExist(err) {
				return nil
			} else {
				return err
			}
		}
		return nil
	})
}

func Login(db *gorm.DB, email string, password string) (Token, error) {
	var user User
	var token Token

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Token{}, fmt.Errorf("user not found")
		}
		return Token{}, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return Token{}, fmt.Errorf("invalid password")
	}

	tx := db.Begin()
	if tx.Error != nil {
		return Token{}, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Where("user_id = ?", user.ID).First(&token).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return Token{}, fmt.Errorf("failed to get token: %w", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || time.Since(token.CreatedAt) > tokenExpiry {
		token = Token{}

		err = UpdateToken(tx, user.ID)
		if err != nil {
			tx.Rollback()
			return Token{}, fmt.Errorf("failed to update token: %w", err)
		}

		err = tx.Table("tokens").Where("user_id = ?", user.ID).First(&token).Error
		if err != nil {
			tx.Rollback()
			return Token{}, fmt.Errorf("failed to retrieve new token: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return Token{}, err
	}

	return token, nil
}

func SearchUsers(db *gorm.DB, query string, friends bool, userId string) ([]utils.UserSearchResponse, error) {
	var users []utils.UserSearchResponse

	baseQuery := db.
		Table("users").
		Select(`
		users.id,
		users.username,
		CASE
			WHEN f.pending = FALSE THEN 1
			WHEN f.pending = TRUE AND f.user_id1 = ? THEN 0
			ELSE -1
		END AS state
	`, userId).
		Joins(`
		LEFT JOIN friends f ON (
			(f.user_id1 = users.id AND f.user_id2 = ?) OR
			(f.user_id2 = users.id AND f.user_id1 = ?)
		)
	`, userId, userId).
		Where("LOWER(users.username) LIKE ? AND users.id != ?", "%"+strings.ToLower(query)+"%", userId)

	if friends {
		baseQuery = baseQuery.Where("f.pending = FALSE")
	}

	err := baseQuery.Order("users.username").Scan(&users).Error
	return users, err
}

func GetUserById(db *gorm.DB, userId string) (utils.UserEssentials, error) {
	var user utils.UserEssentials
	err := db.Table("users").Select("id, username").First(&user, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.UserEssentials{}, fmt.Errorf("user not found")
		}
		return utils.UserEssentials{}, fmt.Errorf("failed to query user: %w", err)
	}
	return user, nil
}

func GetUserExtension(db *gorm.DB, userId string) (utils.UserExtension, error) {
	var extension utils.UserExtension
	err := db.Table("users").Select("id, extension").First(&extension, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.UserExtension{}, fmt.Errorf("user not found")
		}
		return utils.UserExtension{}, fmt.Errorf("failed to query user: %w", err)
	}
	return extension, nil
}

func UpdateExtension(db *gorm.DB, userId string, extension string) error {
	return db.Model(&User{}).Where("id = ?", userId).Update("extension", extension).Error
}
