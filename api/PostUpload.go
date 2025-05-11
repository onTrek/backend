package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PostUpload godoc
// @Summary Upload a GPX file
// @Description Uploads a GPX file for the authenticated user. The file must be sent as form-data with the key 'file'.
// @Tags files
// @Accept application/gpx+xml
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Param file formData file true "GPX file to upload"
// @Success 200 {object} utils.SuccessResponse "File uploaded successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid file"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to save file"
// @Router /gpx/ [post]
func PostUpload(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the gpx file from the form data
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error getting file from form data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	// Check if the file is a GPX file
	if file.Header.Get("Content-Type") != "application/gpx+xml" {
		fmt.Println("Invalid file type:", file.Header.Get("Content-Type"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}
	// Save the file to the server
	var gpx utils.Gpx
	gpx.UserID = user.ID
	gpx.Filename = file.Filename
	gpx.Stats = ""
	err = db.SaveFile(c.MustGet("db").(*sql.DB), gpx, file)
	if err != nil {
		fmt.Println("Error saving file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
