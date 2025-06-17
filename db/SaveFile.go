package db

import (
	"OnTrek/utils"
	"database/sql"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"time"
)

func SaveFile(db *sql.DB, gpx utils.Gpx, file *multipart.FileHeader) error {

	// Save the file to the server
	gpx.StoragePath = "gpxs/" + uuid.New().String()

	err := saveUploadFile(file, gpx.StoragePath)
	if err != nil {
		return err
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO gpx_files (user_id, filename, storage_path, upload_date, title) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(gpx.UserID, gpx.Filename, gpx.StoragePath, time.Now().Format(time.RFC3339), gpx.Title)
	if err != nil {
		return err
	}

	return nil
}

func saveUploadFile(file *multipart.FileHeader, storagePath string) error {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(storagePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the file
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
