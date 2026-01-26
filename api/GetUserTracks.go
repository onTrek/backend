package api

import (
	"OnTrek/db/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserTracks godoc
// @Summary      Get tracks of a specific user
// @Description  Retrieves all tracks associated with a given user ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param        id   path      string  true  "User ID"
// @Success 200 {object} utils.Url "Returns the list of tracks for the user"
// @Failure 400 {object} utils.ErrorResponse "Invalid friend ID"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to retrieve tracks"
// @Router       /users/{id}/gpxs/ [get]
func GetUserTracks(c *gin.Context) {
	friendID := c.Param("id")

	if friendID == "" {
		fmt.Println("Invalid friend ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	friend, err := models.GetUserExtension(c.MustGet("db").(*gorm.DB), friendID)
	if err != nil {
		fmt.Println("Error getting user by ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	tracks, err := models.GetFileByUserID(c.MustGet("db").(*gorm.DB), friend.ID)
	if err != nil {
		fmt.Println("Error getting tracks for user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving tracks"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}
