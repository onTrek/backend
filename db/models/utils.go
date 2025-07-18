package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

func CleanUnusedFiles(db *gorm.DB) error {
	files, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error reading gpxs directory: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", cwd)

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
