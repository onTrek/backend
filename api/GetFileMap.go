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

// GetFileMap godoc
// @Summary      Get a map file by ID
// @Description  Get a map file by its ID
// @Tags         gpx
// @Accept       json
// @Produce      image/png
// @Param Bearer header string true "Bearer token for user authentication"
// @Param        id   path      int  true  "File ID"
// @Success      200 {file} string "Returns the map file as a PNG image"
// @Failure      400 {object} utils.ErrorResponse "Invalid file ID"
// @Failure      401 {object} utils.ErrorResponse "Unauthorized"
// @Failure      404 {object} utils.ErrorResponse "File not found"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /gpx/{id}/map [get]
func GetFileMap(c *gin.Context) {

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

	path := "maps/" + gpx.StoragePath + ".png"
	gpxFile, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
		return
	}
	defer gpxFile.Close()

	// Set the content type and attachment header
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "attachment; filename="+gpx.Title+".png")

	// Send the file content as a response
	c.File(path)
}
