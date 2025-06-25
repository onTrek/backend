package models

import (
	"OnTrek/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Group struct {
	ID          int       `json:"id" example:"1" gorm:"primaryKey;autoIncrement"`
	Description string    `json:"description" example:"Group hike to Monte Bianco" gorm:"not null"`
	CreatedBy   string    `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"type:uuid;not null"`
	CreatedAt   time.Time `json:"created_at" example:"2025-05-11T08:00:00Z" gorm:"not null;default:CURRENT_TIMESTAMP"`
	LastUpdate  time.Time `json:"last_update" example:"2025-05-11T08:00:00Z" gorm:"not null;default:CURRENT_TIMESTAMP"`
	FileId      *int      `json:"file_id" example:"1" gorm:"default:null"`

	User User `json:"user" gorm:"foreignKey:CreatedBy;references:id;constraint:OnDelete:CASCADE"`
	Gpx  Gpx  `json:"gpx" gorm:"foreignKey:FileId;references:id;constraint:OnDelete:SET NULL"`
}

type GroupWithFile struct {
	ID          int       `json:"group_id"`
	CreatedBy   string    `json:"created_by"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdate  time.Time `json:"last_update"`
	FileID      *int      `json:"file_id"`
	FileName    *string   `json:"file_name"`
}

func CreateGroup(db *gorm.DB, group Group) (int, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return fmt.Errorf("error enabling foreign key enforcement: %v", err)
		}

		if err := tx.Create(&group).Error; err != nil {
			return err
		}

		if err := JoinGroup(tx, group.CreatedBy, group.ID); err != nil {
			return fmt.Errorf("error joining group: %v", err)
		}

		return nil
	})

	if err != nil {
		return -1, err
	}

	return group.ID, nil
}

func CheckGroupExistsByIdAndUserId(db *gorm.DB, groupId int, userId string) (bool, error) {
	var count int64
	err := db.Model(&GroupMember{}).
		Where("user_id = ? AND group_id = ?", userId, groupId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckGroupExistsById(db *gorm.DB, groupId int) (bool, error) {
	var count int64
	err := db.Model(&Group{}).
		Where("id = ?", groupId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetGroupsByUserId(db *gorm.DB, userId string) ([]utils.Group, error) {
	var results []GroupWithFile

	err := db.
		Table("groups AS g").
		Select(`g.id, g.created_by, g.description, g.created_at, gf.id AS file_id, gf.filename AS file_name`).
		Joins("JOIN group_members gm ON g.id = gm.group_id").
		Joins("LEFT JOIN gpx_files gf ON g.file_id = gf.id").
		Where("gm.user_id = ?", userId).
		Order("g.last_update DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var groups []utils.Group
	for _, r := range results {
		g := utils.Group{
			ID:          r.ID,
			CreatedBy:   r.CreatedBy,
			Description: r.Description,
			CreatedAt:   r.CreatedAt.UTC().Format(time.RFC3339),
		}

		if r.FileID != nil {
			g.File.ID = *r.FileID
		} else {
			g.File.ID = -1
		}

		if r.FileName != nil {
			g.File.Filename = *r.FileName
		}

		groups = append(groups, g)
	}

	return groups, nil
}

func GetLeaderByGroup(db *gorm.DB, groupId int) (string, error) {
	var leaderId string
	err := db.
		Table("groups").
		Select("created_by").
		Where("id = ?", groupId).
		Scan(&leaderId).Error
	if err != nil {
		return "", err
	}
	return leaderId, nil
}

func DeleteGroupById(db *gorm.DB, userId string, groupId int) error {
	result := db.
		Where("id = ? AND created_by = ?", groupId, userId).
		Delete(&Group{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("nessun gruppo trovato con id %d e created_by %s", groupId, userId)
	}

	return nil
}

func GetGroupInfo(db *gorm.DB, groupId int) (utils.GroupInfoResponse, error) {
	var groupInfo utils.GroupInfoResponse

	type groupRow struct {
		UserID      string
		Username    string
		Description string
		CreatedAt   time.Time
	}

	var result groupRow
	err := db.Table("users AS u").
		Select("u.id as user_id, u.username, g.description, g.created_at").
		Joins("JOIN groups g ON u.id = g.created_by").
		Where("g.id = ?", groupId).
		Scan(&result).Error

	if err != nil {
		return groupInfo, err
	}

	groupInfo.CreatedBy.ID = result.UserID
	groupInfo.CreatedBy.Username = result.Username
	groupInfo.Description = result.Description
	groupInfo.CreatedAt = result.CreatedAt.Format(time.RFC3339)

	type memberRow struct {
		ID       string
		Username string
		Color    string
	}

	var members []memberRow
	err = db.Table("users AS u").
		Select("u.id, u.username, gm.color").
		Joins("JOIN group_members gm ON u.id = gm.user_id").
		Where("gm.group_id = ?", groupId).
		Scan(&members).Error

	if err != nil {
		return utils.GroupInfoResponse{}, err
	}

	for _, m := range members {
		groupInfo.Members = append(groupInfo.Members, utils.GroupMember{
			ID:       m.ID,
			Username: m.Username,
			Color:    m.Color,
		})
	}

	return groupInfo, nil
}

func UpdateFileForTheGroup(db *gorm.DB, groupId int, fileId int) error {
	err := db.Model(&Group{}).
		Where("id = ?", groupId).
		Updates(map[string]interface{}{
			"file_id":     fileId,
			"last_update": time.Now().UTC().Format(time.RFC3339),
		}).Error
	return err
}
