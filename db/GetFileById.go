package db

import (
	"OnTrek/utils"
	"database/sql"
)

func GetFileByID(db *sql.DB, fileID int, userID string) (*utils.Gpx, error) {
	// Create a new file instance
	file := &utils.Gpx{}

	// Query the database for the file with the given ID
	err := db.QueryRow("SELECT id, activity_id, user_id, filename, storage_path, upload_date, stats FROM gpx_files WHERE id = ? AND user_id = ?", fileID, userID).Scan(&(file.ID), &(file.ActivityID), &(file.UserID), &(file.Filename), &(file.StoragePath), &(file.UploadDate), &(file.Stats))
	if err != nil {
		return nil, err
	}

	return file, nil
}
