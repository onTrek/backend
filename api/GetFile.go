package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func GetFile(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the file ID from the URL parameter
	file := c.Param("id")
	// Validate the file ID
	fileID, err := strconv.Atoi(file)
	if err != nil {
		fmt.Println("Error converting file ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Fetch the file from the database
	gpx, err := db.GetFileByID(c.MustGet("db").(*sql.DB), fileID, user.ID)
	if err != nil {
		fmt.Println("Error fetching file from database:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	gpxFile, err := os.OpenFile(gpx.StoragePath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
		return
	}
	defer gpxFile.Close()

	// Set the content type and attachment header
	c.Header("Content-Type", "application/gpx+xml")
	c.Header("Content-Disposition", "attachment; filename="+gpx.Filename)

	// Send the file content as a response
	c.File(gpx.StoragePath)
}
