package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSearchTrack godoc
// @Summary Search for GPX tracks by title
// @Description Searches for GPX tracks by title
// @Tags search
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param track query string true "Search for track title"
// @Success 200 {array} utils.GpxInfoEssential "Returns a list of GPX tracks matching the search query ordered by upload date."
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "No tracks found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /search/tracks [get]
func GetSearchTrack(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	query := c.Query("track")
	if query == "" {
		fmt.Println("Error: track parameter is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "track parameter is required"})
		return
	}

	tracks, err := models.SearchGpxs(c.MustGet("db").(*gorm.DB), query, user.ID)
	if err != nil {
		fmt.Println("Error searching tracks:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tracks"})
		return
	}

	if len(tracks) == 0 {
		c.JSON(200, []utils.GpxInfoEssential{})
		return
	}

	c.JSON(http.StatusOK, tracks)
}
