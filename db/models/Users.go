package models

import (
	"OnTrek/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	err := db.
		Where("id = ?", userID).
		Delete(&User{}).Error

	return err
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

func SearchUsers(db *gorm.DB, query string, userId string) ([]utils.UserEssentials, error) {
	var users []utils.UserEssentials

	err := db.
		Table("users").
		Select("id, username").
		Where("LOWER(username) LIKE ? AND id != ?", "%"+strings.ToLower(query)+"%", userId).
		Limit(100).
		Scan(&users).Error

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
