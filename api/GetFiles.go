package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetFiles godoc
// @Summary Retrieve user's GPX files
// @Description Returns a list of GPX files associated with the authenticated user
// @Tags gpx
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Success 200 {object} []utils.GpxInfo "gpx_files"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Error fetching files"
// @Router /gpx/ [get]
func GetFiles(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.User)

	// Get files from the database
	files, err := db.GetFiles(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		fmt.Println("Error fetching files:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching files"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusOK, []utils.GpxInfo{})
		return
	}

	// Return the files as JSON
	c.JSON(http.StatusOK, files)
}
