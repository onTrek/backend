package models

import (
	"OnTrek/utils"
	"fmt"
	"mime/multipart"
	"time"

	firebaseStorage "firebase.google.com/go/v4/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gpx struct {
	ID          int       `json:"id" example:"1" gorm:"primaryKey;autoIncrement"`
	UserID      string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"type:uuid;not null"`
	Filename    string    `json:"filename" example:"MonteBianco.gpx" gorm:"not null"`
	StoragePath string    `json:"storage_path" example:"123e4567-e89b-12d3-a456-426614174000.gpx" gorm:"not null;unique"`
	UploadDate  time.Time `json:"upload_date" example:"2025-05-11T08:00:00Z" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Title       string    `json:"title" example:"Monte Faggeto" gorm:"not null"`
	KM          float64   `json:"km" example:"14.5" gorm:"not null;default:0"`
	Ascent      int       `json:"ascent" example:"1000" gorm:"not null;default:0"`
	Descent     int       `json:"descent" example:"1000" gorm:"not null;default:0"`
	Duration    string    `json:"duration" example:"06:30:00" gorm:"not null;default:''"`
	MaxAltitude int       `json:"max_altitude" example:"2500" gorm:"not null;default:0"`
	MinAltitude int       `json:"min_altitude" example:"1500" gorm:"not null;default:0"`
	Size        int64     `json:"size" example:"2048000" gorm:"not null;default:0"`
	Public      bool      `json:"public" example:"false" gorm:"not null;default:0"`

	User User `json:"user" gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE"`
}

func (Gpx) TableName() string {
	return "gpx_files"
}

func GetFiles(db *gorm.DB, userID string) ([]utils.GpxInfo, error) {
	var files []Gpx

	err := db.
		Table("gpx_files").
		Select(`id, filename, storage_path, upload_date, 
                title, km, duration, ascent, descent, max_altitude, min_altitude, size`).
		Where("user_id = ?", userID).Order("upload_date DESC").
		Scan(&files).Error
	if err != nil {
		return nil, err
	}

	var result []utils.GpxInfo
	for _, row := range files {
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

func GetFileByID(db *gorm.DB, fileID int) (utils.Gpx, error) {
	var file utils.Gpx

	err := db.Table("gpx_files").
		Select(`id, user_id, filename, storage_path, upload_date, title`).
		Where("id = ?", fileID).
		First(&file).Error

	if err != nil {
		return file, err
	}

	return file, nil
}

func GetFileByUserID(db *gorm.DB, userID string) ([]utils.GpxInfoEssential, error) {
	var gpxs []Gpx

	err := db.Table("gpx_files").
		Where("user_id = ? AND public = 1", userID).
		Find(&gpxs).Error
	if err != nil {
		return nil, err
	}

	var result []utils.GpxInfoEssential
	for _, row := range gpxs {
		info := utils.GpxInfoEssential{
			ID:    row.ID,
			Title: row.Title,
		}

		result = append(result, info)
	}

	return result, nil
}

func GetFileInfoByID(db *gorm.DB, fileID int) (utils.GpxInfoWithOwner, error) {
	var file Gpx

	err := db.Table("gpx_files").
		Where("id = ?", fileID).
		First(&file).Error

	if err != nil {
		return utils.GpxInfoWithOwner{}, err
	}

	info := utils.GpxInfoWithOwner{
		ID:         file.ID,
		Owner:      file.UserID,
		Filename:   file.Filename,
		UploadDate: file.UploadDate.Format(time.RFC3339),
		Title:      file.Title,
		Stats: utils.GPXStats{
			Km:          file.KM,
			Duration:    file.Duration,
			Ascent:      file.Ascent,
			Descent:     file.Descent,
			MaxAltitude: file.MaxAltitude,
			MinAltitude: file.MinAltitude,
		},
		FileSize: file.Size,
		Public:   file.Public,
	}

	return info, nil
}

func GetFileByIDAndUserID(db *gorm.DB, fileID int, userID string) (utils.Gpx, error) {
	var file utils.Gpx

	err := db.Table("gpx_files").Select(`id, user_id, filename, storage_path, upload_date, title`).
		Where("id = ? AND user_id = ?", fileID, userID).
		First(&file).Error

	if err != nil {
		return file, err
	}

	return file, nil
}

func DeleteFileByID(db *gorm.DB, client *firebaseStorage.Client, storage *utils.StorageConfig, fileID int, userID string, gpx utils.Gpx) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return fmt.Errorf("error enabling foreign key enforcement: %v", err)
		}

		if err := tx.Table("gpx_files").Where("id = ? AND user_id = ?", fileID, userID).Delete(&utils.Gpx{}).Error; err != nil {
			return err
		}

		if err := utils.DeleteFiles(client, storage, gpx); err != nil {
			return err
		}

		return nil
	})
}

func SaveFile(db *gorm.DB, client *firebaseStorage.Client, storage *utils.StorageConfig, gpx Gpx, file *multipart.FileHeader) (int, error) {
	var createdID int

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return fmt.Errorf("error enabling foreign key enforcement: %v", err)
		}

		stats, err := utils.CalculateStats(file)
		if err != nil {
			return err
		}

		gpx.StoragePath = uuid.New().String()
		gpx.KM = stats.Km
		gpx.Ascent = stats.Ascent
		gpx.Descent = stats.Descent
		gpx.Duration = stats.Duration
		gpx.MaxAltitude = stats.MaxAltitude
		gpx.MinAltitude = stats.MinAltitude
		gpx.Size = file.Size

		if err := tx.Create(&gpx).Error; err != nil {
			return err
		}

		createdID = gpx.ID

		if _, err := utils.SaveFile(client, storage, file, "gpxs", gpx.StoragePath, ".gpx"); err != nil {
			return err
		}

		if _, err := utils.CreateMap(file, client, storage, gpx.StoragePath); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error during transaction:", err)
		fmt.Println("Attempting to delete files due to error...")
		err = utils.DeleteFiles(client, storage, utils.Gpx{StoragePath: gpx.StoragePath})
		if err != nil {
			fmt.Println("Error deleting files after transaction failure:", err)
			return -1, err
		} else {
			fmt.Println("Files deleted successfully after transaction failure.")
		}
		return -1, err
	}

	return createdID, err
}

func CheckFilePermissions(db *gorm.DB, fileID int, userID string) (bool, error) {
	var count int64

	err := db.Table("gpx_files gf").
		Joins("LEFT JOIN groups g ON g.file_id = gf.id").
		Joins("LEFT JOIN group_members gm ON gm.group_id = g.id").
		Where("gf.id = ? AND (gf.public = 1 OR gf.user_id = ? OR gm.user_id = ?)", fileID, userID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func UpdateFilePrivacy(db *gorm.DB, fileID int, isPublic bool) error {
	publicValue := 0
	if isPublic {
		publicValue = 1
	}

	result := db.Model(&Gpx{}).Where("id = ?", fileID).Update("public", publicValue)
	if result.Error != nil {
		fmt.Println("Error updating file privacy:", result.Error)
		return result.Error
	}

	return nil
}

func SearchGpxs(db *gorm.DB, query string, userID string) ([]utils.GpxInfoEssential, error) {
	var gpxs []Gpx

	err := db.Table("gpx_files").
		Where("(title LIKE ?) AND (public = 1) AND (user_id != ?)", "%"+query+"%", userID).
		Find(&gpxs).Error
	if err != nil {
		return nil, err
	}

	var result []utils.GpxInfoEssential
	for _, row := range gpxs {
		info := utils.GpxInfoEssential{
			ID:    row.ID,
			Title: row.Title,
		}

		result = append(result, info)
	}

	return result, nil
}
