package models

import (
	"OnTrek/utils"
	"fmt"
	"gorm.io/gorm"
)

type Friend struct {
	UserId1 string `json:"user_id1" gorm:"type:uuid;primaryKey"`
	UserId2 string `json:"user_id2" gorm:"type:uuid;primaryKey"`
	Pending bool   `json:"pending" gorm:"default:true"`

	User1 User `gorm:"foreignKey:UserId1;references:ID"`
	User2 User `gorm:"foreignKey:UserId2;references:ID"`
}

func GetFriends(db *gorm.DB, userID string) ([]utils.UserEssentials, error) {
	var friends []utils.UserEssentials

	subQuery1 := db.Table("friends").
		Select("user_id1").
		Where("user_id2 = ? AND pending = FALSE", userID)

	subQuery2 := db.Table("friends").
		Select("user_id2").
		Where("user_id1 = ? AND pending = FALSE", userID)

	err := db.Table("users").
		Where("id IN (?) OR id IN (?)", subQuery1, subQuery2).
		Find(&friends).Error

	if err != nil {
		return nil, err
	}

	return friends, nil
}

func DeleteFriend(db *gorm.DB, userID, friendID string) error {
	res := db.Where(
		"(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)",
		userID, friendID, friendID, userID,
	).Delete(&Friend{})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func AddFriend(db *gorm.DB, userID string, friendID string) error {

	var count int64
	err := db.Model(&Friend{}).
		Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)", userID, friendID, friendID, userID).
		Where("pending = FALSE").
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("users are already friends")
	}

	// Add the friend request with pending = TRUE
	friend := Friend{
		UserId1: userID,
		UserId2: friendID,
		Pending: true,
	}
	if err := db.Create(&friend).Error; err != nil {
		return err
	}

	return nil
}

func GetFriendRequestsByUserId(db *gorm.DB, userId string) ([]utils.UserEssentials, error) {
	var friendRequests []utils.UserEssentials

	err := db.Table("users u").
		Select("u.id, u.username").
		Joins("JOIN friends f ON u.id = f.user_id1").
		Where("f.user_id2 = ? AND f.pending = TRUE", userId).
		Scan(&friendRequests).Error

	if err != nil {
		return nil, err
	}

	return friendRequests, nil
}

func AcceptFriendRequest(db *gorm.DB, userID, friendID string) error {
	result := db.Model(&Friend{}).
		Where("user_id1 = ? AND user_id2 = ? AND pending = TRUE", friendID, userID).
		Update("pending", false)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeleteFriendRequest(db *gorm.DB, userID, friendID string) error {
	err := db.
		Where("user_id1 = ? AND user_id2 = ? AND pending = TRUE", friendID, userID).
		Delete(&Friend{}).Error
	return err
}
