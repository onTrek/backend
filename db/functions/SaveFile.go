package functions

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

func SaveFile(db *sql.DB, gpx utils.Gpx, file *multipart.FileHeader) error {

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Rollback on error unless committed
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Enable foreign key enforcement
	_, err = tx.Exec("PRAGMA foreign_keys = ON") // Enable foreign key enforcement
	if err != nil {
		return fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

	stats, err := utils.CalculateStats(file)
	if err != nil {
		return err
	}

	// Save the file to the server

	gpx.StoragePath = uuid.New().String()

	// Prepare the SQL statement
	stmt, err := tx.Prepare("INSERT INTO gpx_files (user_id, filename, storage_path, upload_date, title, km, ascent, descent, duration, max_altitude, min_altitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(gpx.UserID, gpx.Filename, gpx.StoragePath, time.Now().Format(time.RFC3339), gpx.Title, stats.Km, stats.Ascent, stats.Descent, stats.Duration, stats.MaxAltitude, stats.MinAltitude)
	if err != nil {
		return err
	}

	err = utils.SaveFile(file, gpx.StoragePath)
	if err != nil {
		return err
	}

	err = utils.CreateMap(file, gpx.StoragePath)
	if err != nil {
		return fmt.Errorf("error creating map: %v", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
