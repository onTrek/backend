package api

import (
	"OnTrek/db/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GetFileInfo godoc
// @Summary Get File Info
// @Description Retrieve information about a specific GPX file by its ID
// @Tags gpx
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "File ID"
// @Success 200 {object} utils.GpxInfoWithOwner "Gpx file information"
// @Failure 400 {object} utils.ErrorResponse "Invalid file ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "File not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /gpx/{id} [get]
func GetFileInfo(c *gin.Context) {

	// Get the file ID from the URL parameter
	file := c.Param("id")
	// Validate the file ID
	fileID, err := strconv.Atoi(file)
	if err != nil {
		fmt.Println("Error converting file ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	gpx, err := models.GetFileInfoByID(c.MustGet("db").(*gorm.DB), fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("File not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		fmt.Println("Error fetching file from database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching file from database"})
		return
	}

	c.JSON(http.StatusOK, gpx)

}
