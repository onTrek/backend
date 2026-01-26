package models

import (
	"OnTrek/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SavedTrack struct {
	UserId  string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"primaryKey;type:uuid;not null"`
	FileId  int       `json:"track_id" gorm:"primaryKey;column:track_id;not null"`
	SavedAt time.Time `json:"created_at" example:"2025-05-11T08:00:00Z" gorm:"not null;default:CURRENT_TIMESTAMP"`

	User User `json:"user" gorm:"foreignKey:UserId;references:id;constraint:OnDelete:CASCADE"`
	File Gpx  `json:"track" gorm:"foreignKey:FileId;references:id;constraint:OnDelete:CASCADE"`
}

func SaveTrack(db *gorm.DB, userID string, fileID int) error {
	var track Gpx

	if err := db.Select("id, user_id, public").First(&track, fileID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	if track.UserID == userID {
		return errors.New("Cannot save your own track")
	}

	if track.Public == false {
		return errors.New("Cannot save a private track")
	}

	savedTrack := SavedTrack{
		UserId: userID,
		FileId: fileID,
	}
	result := db.Create(&savedTrack)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UnsaveTrack(db *gorm.DB, userID string, fileID int) error {
	return db.Delete(&SavedTrack{}, "user_id = ? AND track_id = ?", userID, fileID).Error
}

func GetSavedTracks(db *gorm.DB, userID string) ([]utils.GpxInfo, error) {
	var savedTracks []Gpx
	err := db.Table("gpx_files").
		Select("gpx_files.*").
		Joins("join saved_tracks on saved_tracks.track_id = gpx_files.id").
		Where("saved_tracks.user_id = ?", userID).
		Order("saved_tracks.saved_at DESC").
		Find(&savedTracks).Error
	if err != nil {
		return nil, err
	}

	var result []utils.GpxInfo
	for _, row := range savedTracks {
		info := utils.GpxInfo{
			ID:         row.ID,
			Filename:   row.Filename,
			UploadDate: row.UploadDate.Format(time.RFC3339),
			Title:      row.Title,
			Stats: utils.GPXStats{
				Km:          row.KM,
				Duration:    row.Duration,
				Ascent:      row.Ascent,
				Descent:     row.Descent,
				MaxAltitude: row.MaxAltitude,
				MinAltitude: row.MinAltitude,
			},
			FileSize: row.Size, // Convert to KB
		}

		result = append(result, info)
	}

	return result, nil
}
