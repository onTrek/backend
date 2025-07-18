package models

import (
	"OnTrek/utils"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
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

func GetFileByPath(db *gorm.DB, storagePath string) (utils.Gpx, error) {
	var file utils.Gpx

	err := db.Table("gpx_files").
		Select(`id, user_id, filename, storage_path, upload_date, title`).
		Where("storage_path = ?", storagePath).
		First(&file).Error

	if err != nil {
		return utils.Gpx{}, err
	}

	return file, nil
}

func GetFileInfoByID(db *gorm.DB, fileID int) (utils.GpxInfo, error) {
	var file Gpx

	err := db.Table("gpx_files").
		Where("id = ?", fileID).
		First(&file).Error

	if err != nil {
		return utils.GpxInfo{}, err
	}

	info := utils.GpxInfo{
		ID:         file.ID,
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

func DeleteFileByID(db *gorm.DB, fileID int, userID string, gpx utils.Gpx) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return fmt.Errorf("error enabling foreign key enforcement: %v", err)
		}

		if err := tx.Table("gpx_files").Where("id = ? AND user_id = ?", fileID, userID).Delete(&utils.Gpx{}).Error; err != nil {
			return err
		}

		if err := utils.DeleteFiles(gpx); err != nil {
			return err
		}

		return nil
	})
}

func SaveFile(db *gorm.DB, gpx Gpx, file *multipart.FileHeader) (int, error) {
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

		if err := utils.SaveFile(file, "gpxs", gpx.StoragePath, ".gpx"); err != nil {
			return err
		}

		if err := utils.CreateMap(file, gpx.StoragePath); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error during transaction:", err)
		fmt.Println("Attempting to delete files due to error...")
		err = utils.DeleteFiles(utils.Gpx{StoragePath: gpx.StoragePath})
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

func CleanUnusedFiles(db *gorm.DB) error {
	files, err := os.ReadDir("./root")
	if err != nil {
		return fmt.Errorf("error reading gpxs directory: %w", err)
	}

	for _, file := range files {
		fileName := file.Name()

		// Skip db and png files
		if fileName == "ontrek.db" || file.IsDir() {
			continue
		}

		// Get file without any extension
		fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))] + filepath.Ext(file.Name())
		fmt.Println("Checking file:", fileName)

	}

	return nil
}
