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

// DeleteSavedTrack godoc
// @Summary      Delete a saved track
// @Description  Remove a track from the user's saved tracks
// @Tags         gpx
// @Produce      json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "File ID"
// @Success 200 {object} utils.SuccessResponse "Track saved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid file ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Error saving track"
// @Router       /gpx/{id}/unsave [delete]
func DeleteSavedTrack(c *gin.Context) {
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

	err = models.UnsaveTrack(c.MustGet("db").(*gorm.DB), user.ID, fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Saved track not found"})
			return
		}
		fmt.Println("Error deleting saved track:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete saved track"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Saved track deleted successfully"})
}
