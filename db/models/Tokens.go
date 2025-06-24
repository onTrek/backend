package models

import (
	"OnTrek/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Token struct {
	UserId    string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"type:uuid;primaryKey;"`
	Token     string    `json:"token" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"created_at" example:"2025-05-11T08:00:00Z" gorm:"not null;default:CURRENT_TIMESTAMP"`

	User User `json:"user" gorm:"foreignKey:UserId;references:id;constraint:,OnDelete:CASCADE;"`
}

func UpdateToken(tx *gorm.DB, userID string) error {
	token := Token{
		UserId: userID,
		Token:  uuid.New().String(),
	}

	// Delete existing token for the user
	if err := tx.Where("user_id = ?", userID).Delete(&Token{}).Error; err != nil {
		return err
	}

	// Create a new token for the user
	return tx.Create(&token).Error
}

func GetUserByToken(db *gorm.DB, tokenStr string) (utils.UserInfo, error) {
	var token Token

	err := db.Preload("User").Where("token = ?", tokenStr).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.UserInfo{}, fmt.Errorf("user not found")
		}
		return utils.UserInfo{}, fmt.Errorf("failed to query token: %w", err)
	}

	if time.Since(token.CreatedAt) > tokenExpiry {
		return utils.UserInfo{}, fmt.Errorf("token expired")
	}

	return utils.UserInfo{
		ID:       token.UserId,
		Email:    token.User.Email,
		Username: token.User.Username,
	}, nil
}
