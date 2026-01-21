package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSavedTracks godoc
// @Summary Retrieve user's saved GPX files
// @Description Returns a list of GPX files that the authenticated user has saved
// @Tags gpx
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Success 200 {object} []utils.GpxInfo "gpx_files"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Error fetching files"
// @Router /gpx/ [get]
func GetSavedTracks(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	files, err := models.GetSavedTracks(c.MustGet("db").(*gorm.DB), user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching saved tracks"})
		return
	}

	if len(files) == 0 {
		c.JSON(200, []utils.GpxInfo{})
		return
	}

	c.JSON(http.StatusOK, files)
}
