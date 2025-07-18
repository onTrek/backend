package models

import (
	"OnTrek/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type User struct {
	ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"type:uuid;primaryKey"`
	Email        string    `json:"email" example:"user@example.com" gorm:"unique;not null"`
	Username     string    `json:"username" example:"John Doe" gorm:"not null"`
	PasswordHash string    `json:"password_hash" example:"strongPassword123" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" example:"2025-05-11T08:00:00Z" gorm:"not null;default:CURRENT_TIMESTAMP"`
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

func Login(db *gorm.DB, email string, password string) (utils.UserToken, error) {
	var user User
	var token Token

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.UserToken{}, fmt.Errorf("user not found")
		}
		return utils.UserToken{}, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return utils.UserToken{}, fmt.Errorf("invalid password")
	}

	tx := db.Begin()
	if tx.Error != nil {
		return utils.UserToken{}, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Where("user_id = ?", user.ID).First(&token).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return utils.UserToken{}, fmt.Errorf("failed to get token: %w", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || time.Since(token.CreatedAt) > tokenExpiry {
		token = Token{}

		err = UpdateToken(tx, user.ID)
		if err != nil {
			tx.Rollback()
			return utils.UserToken{}, fmt.Errorf("failed to update token: %w", err)
		}

		err = tx.Table("tokens").Where("user_id = ?", user.ID).First(&token).Error
		if err != nil {
			tx.Rollback()
			return utils.UserToken{}, fmt.Errorf("failed to retrieve new token: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return utils.UserToken{}, err
	}

	return utils.UserToken{
		Token: token.Token,
	}, nil
}

func SearchUsers(db *gorm.DB, query string, friends bool, userId string) ([]utils.UserEssentials, error) {
	var users []utils.UserEssentials

	baseQuery := db.
		Table("users").
		Select("users.id, users.username").
		Where("LOWER(users.username) LIKE ? AND users.id != ?", "%"+strings.ToLower(query)+"%", userId)

	if friends {
		baseQuery = baseQuery.Joins(`
			JOIN friends ON (
				(friends.user_id1 = users.id AND friends.user_id2 = ?) OR
				(friends.user_id2 = users.id AND friends.user_id1 = ?)
			)
		`, userId, userId).Where("friends.pending = FALSE")
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

func CleanUnusedProfileImages(db *gorm.DB) error {
	files, err := os.ReadDir("profile")
	if err != nil {
		return fmt.Errorf("error reading gpxs directory: %w", err)
	}

	for _, file := range files {
		fileName := file.Name()

		// Skip db file and directories
		if fileName == "ontrek.db" || file.IsDir() {
			continue
		}

		// Get file without any extension
		file := fileName[:len(fileName)-len(filepath.Ext(fileName))]
		_, err = GetUserById(db, file)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				err = os.Remove(filepath.Join("profile", fileName))
				if err != nil {
					return fmt.Errorf("error deleting unused profile image %s: %w", file, err)
				}
				fmt.Println("Deleted unused profile image:", fileName)
			} else {
				return fmt.Errorf("error checking user %s in database: %w", file, err)
			}
		}
	}

	return nil
}
