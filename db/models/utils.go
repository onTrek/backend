package models

import (
	"OnTrek/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
)

func CleanUnusedFiles(db *gorm.DB) error {
	files, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error reading gpxs directory: %w", err)
	}

	for _, file := range files {
		fileName := file.Name()

		if fileName == "ontrek.db" || fileName == "db" {
			continue
		}

		if fileName == "gpxs" || fileName == "avatars" {
			if file.IsDir() {
				fmt.Println("Checking directory:", file.Name())
				subFiles, err := os.ReadDir(filepath.Join(".", file.Name()))
				if err != nil {
					return fmt.Errorf("error reading subdirectory %s: %w", file.Name(), err)
				}

				for _, subFile := range subFiles {
					subFileName := subFile.Name()

					if subFileName == "ontrek.db" {
						continue
					}

					subFileName = subFileName[:len(subFileName)-len(filepath.Ext(subFileName))] + filepath.Ext(file.Name())

					if fileName == "gpxs" {
						_, err = GetFileByPath(db, subFileName)
						if err != nil {
							if errors.Is(err, gorm.ErrRecordNotFound) {
								err = utils.DeleteFiles(utils.Gpx{StoragePath: subFileName})
								if err != nil {
									return fmt.Errorf("error deleting file %s: %w", subFileName, err)
								}
								fmt.Println("Deleted unused file:", subFileName)
							} else {
								return fmt.Errorf("error checking file %s in database: %w", subFileName, err)
							}
						}
					} else {
						_, err = GetUserById(db, subFileName)
						if err != nil {
							if strings.Contains(err.Error(), "user not found") {
								err = os.Remove(filepath.Join(".", fileName, subFile.Name()))
								if err != nil {
									return fmt.Errorf("error deleting file %s: %w", subFileName, err)
								}
								fmt.Println("Deleted unused user file:", subFileName)
							} else {
								return fmt.Errorf("error checking user %s in database: %w", subFileName, err)
							}
						}
					}

				}
			}
		}
	}

	return nil
}
