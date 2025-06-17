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

// GetFile godoc
// @Summary Retrieve a file by ID
// @Description Retrieves a file based on the provided file ID and authorization token
// @Tags files
// @Accept json
// @Produce octet-stream
// @Param Authorization header string true "Bearer token for user authentication"
// @Param id path int true "File ID"
// @Success 200 {file} file "Returns the requested file"
// @Failure 400 {object} utils.ErrorResponse "Invalid file ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "File not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /gpx/{id} [get]
func GetFile(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	_, err := db.GetUserByToken(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		if err.Error() == "token expired" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}
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
	gpx, err := db.GetFileByID(c.MustGet("db").(*sql.DB), fileID)
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
