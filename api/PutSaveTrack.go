package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PutSaveTrack godoc
// @Summary Save Track
// @Description Save a GPX track to the user's saved tracks
// @Tags gpx
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "File ID"
// @Success 200 {object} utils.SuccessResponse "Track saved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid file ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Error saving track"
// @Router /gpx/{id}/save [put]
func PutSaveTrack(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	// Get the file ID from the URL parameter
	file := c.Param("id")
	// Validate the file ID
	fileID, err := strconv.Atoi(file)
	if err != nil {
		fmt.Println("Error converting file ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	err = models.SaveTrack(c.MustGet("db").(*gorm.DB), user.ID, fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			fmt.Println("Track already saved:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Track already saved"})
			return
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("File not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		} else if err.Error() == "Cannot save your own track" {
			fmt.Println("Error saving track:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "Cannot save a private track" {
			fmt.Println("Error saving track:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Error saving track:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving track"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Track saved successfully"})
}
