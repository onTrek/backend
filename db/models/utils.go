package models

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

func CleanUnusedFiles(db *gorm.DB) error {
	files, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error reading gpxs directory: %w", err)
	}

	for _, file := range files {
		fileName := file.Name()

		if fileName == "gpxs" || fileName == "maps" || fileName == "profile" {
			if file.IsDir() {
				fmt.Println("Checking directory:", file.Name())
				subFiles, err := os.ReadDir(filepath.Join(".", file.Name()))
				if err != nil {
					return fmt.Errorf("error reading subdirectory %s: %w", file.Name(), err)
				}
				for _, subFile := range subFiles {
					subFileName := subFile.Name()

					// Skip db
					if subFileName == "ontrek.db" {
						continue
					}
					// Get file without any extension
					subFileName = subFileName[:len(subFileName)-len(filepath.Ext(subFileName))] + filepath.Ext(file.Name())
					fmt.Println("Checking file:", subFileName)
				}
			}
		}
	}

	return nil
}
