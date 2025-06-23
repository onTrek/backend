package api

import (
	"OnTrek/db/functions"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

// GetFile godoc
// @Summary Download a GPX file by ID
// @Description Retrieves a file based on the provided file ID and authorization token
// @Tags gpx
// @Accept json
// @Produce octet-stream
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "File ID"
// @Success 200 {file} file "Returns the requested file as .gpx"
// @Failure 400 {object} utils.ErrorResponse "Invalid file ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "File not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /gpx/{id}/download [get]
func GetFile(c *gin.Context) {

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
	gpx, err := functions.GetFileByID(c.MustGet("db").(*sql.DB), fileID)
	if err != nil {
		fmt.Println("Error fetching file from database:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	path := "gpxs/" + gpx.StoragePath
	gpxFile, err := os.OpenFile(path, os.O_RDONLY, 0644)
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
	c.File(path)
}
