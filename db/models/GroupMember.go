package models

import (
	"OnTrek/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type GroupMember struct {
	GroupId     int       `json:"group_id" example:"1" gorm:"primaryKey;not null"`
	UserId      string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"primaryKey;type:uuid;uniqueIndex:groupColor;not null"`
	Latitude    float64   `json:"latitude" example:"45.123456" gorm:"not null;default:-1"`
	Longitude   float64   `json:"longitude" example:"7.123456" gorm:"not null;default:-1"`
	Altitude    float64   `json:"altitude" example:"1500" gorm:"not null;default:-1"`
	Accuracy    float64   `json:"accuracy" example:"5.0" gorm:"not null;default:-1"`
	HelpRequest bool      `json:"help_request" example:"false" gorm:"not null;default:false"`
	GoingTo     *string   `json:"going_to" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"default:null;type:uuid"`
	Timestamp   time.Time `json:"timestamp" example:"2025-05-11T08:00:00Z" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Color       string    `json:"color" example:"#FF5733" gorm:"uniqueIndex:groupColor;not null"`

	Group       Group `json:"group" gorm:"foreignKey:GroupId;references:id;constraint:OnDelete:CASCADE"`
	User        User  `json:"user" gorm:"foreignKey:UserId;references:id;constraint:OnDelete:CASCADE"`
	GoingToUser User  `json:"going_to_user" gorm:"foreignKey:GoingTo;references:id;constraint:OnDelete:SET NULL"`
}

func JoinGroup(tx *gorm.DB, userID string, groupID int) (utils.GroupMember, error) {
	colors := []string{
		"#e6194b", "#3cb44b", "#ffe119", "#4363d8", "#f58231", "#911eb4", "#46f0f0",
		"#f032e6", "#bcf60c", "#fabebe", "#008080", "#e6beff", "#9a6324", "#fffac8",
		"#800000", "#aaffc3", "#808000", "#ffd8b1", "#000075", "#808080",
	}

	var usedColors []string
	if err := tx.Model(&utils.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("color", &usedColors).Error; err != nil {
		return utils.GroupMember{}, err
	}

	used := make(map[string]bool)
	for _, c := range usedColors {
		used[c] = true
	}

	var selectedColor string
	for _, c := range colors {
		if !used[c] {
			selectedColor = c
			break
		}
	}

	if selectedColor == "" {
		return utils.GroupMember{}, fmt.Errorf("no available colors left for group %d", groupID)
	}

	member := GroupMember{
		GroupId: groupID,
		UserId:  userID,
		Color:   selectedColor,
	}

	if err := tx.Create(&member).Error; err != nil {
		return utils.GroupMember{}, err
	}

	res := utils.GroupMember{
		ID:    userID,
		Color: selectedColor,
	}

	return res, nil
}

func UpdateGroupMember(db *gorm.DB, userId string, group GroupMember) error {
	var updateData map[string]interface{}

	if group.GoingTo != "" {
		updateData = map[string]interface{}{
			"latitude":     group.Latitude,
			"longitude":    group.Longitude,
			"altitude":     group.Altitude,
			"accuracy":     group.Accuracy,
			"help_request": group.HelpRequest,
			"going_to":     group.GoingTo,
			"timestamp":    time.Now().UTC().Format(time.RFC3339),
		}
	} else {
		updateData = map[string]interface{}{
			"latitude":     group.Latitude,
			"longitude":    group.Longitude,
			"altitude":     group.Altitude,
			"accuracy":     group.Accuracy,
			"help_request": group.HelpRequest,
			"timestamp":    time.Now().UTC().Format(time.RFC3339),
		}
	}

	return db.Model(&utils.GroupMember{}).
		Where("group_id = ? AND user_id = ?", group.GroupId, userId).
		Updates(updateData).Error
}

func GetMembersInfoByGroupId(db *gorm.DB, groupId int) ([]utils.MemberInfo, error) {
	var groupMembers []GroupMember
	err := db.Preload("User").
		Where("group_id = ?", groupId).
		Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}

	var members []utils.MemberInfo
	for _, gm := range groupMembers {
		member := utils.MemberInfo{
			User: utils.UserEssentials{
				ID:       gm.User.ID,
				Username: gm.User.Username,
			},
			Latitude:      gm.Latitude,
			Longitude:     gm.Longitude,
			Altitude:      gm.Altitude,
			Accuracy:      gm.Accuracy,
			HelpRequested: gm.HelpRequest,
			TimeStamp:     gm.Timestamp.Format(time.RFC3339),
		}

		if gm.GoingTo != nil {
			member.GoingTo = *gm.GoingTo
		} else {
			member.GoingTo = ""
		}

		members = append(members, member)
	}

	return members, nil
}

func GetMemberByUserIdAndGroupId(tx *gorm.DB, userId string, groupId int) (utils.GroupMember, error) {
	var member utils.GroupMember
	err := tx.
		Select("gm.user_id, u.username, gm.group_id, gm.color").
		Where("user_id = ? AND group_id = ?", userId, groupId).
		Joins("JOIN users u ON gm.user_id = u.id").
		First(&member).Error
	if err != nil {
		return utils.GroupMember{}, err
	}
	return member, nil
}

func LeaveGroupById(db *gorm.DB, userId string, groupId int) error {
	err := db.
		Where("user_id = ? AND group_id = ?", userId, groupId).
		Delete(&GroupMember{}).Error
	return err
}
