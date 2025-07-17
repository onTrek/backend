package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// PostUpload godoc
// @Summary Upload a GPX file
// @Description Uploads a GPX file for the authenticated user. The file must be sent as form-data with the key 'file'.
// @Tags gpx
// @Accept multipart/form-data
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param file formData file true "GPX file to upload"
// @Param title formData string true "Title for the GPX file(max 64 characters)"
// @Success 201 {object} utils.GpxID "File id of the uploaded GPX file"
// @Failure 400 {object} utils.ErrorResponse "Invalid file"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to save file"
// @Router /gpx/ [post]
func PostUpload(c *gin.Context) {
	var title string

	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	// Get the gpx file from the form data
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error getting file from form data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// Get the title from the form data
	title = c.PostForm("title")

	if title == "" {
		fmt.Println("Title is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if len(title) > 64 {
		fmt.Println("Title is too long, maximum 64 characters allowed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is too long, maximum 64 characters allowed"})
		return
	}

	// Check if the file is a GPX file
	if file.Header.Get("Content-Type") != "application/gpx+xml" {
		fmt.Println("Invalid file type:", file.Header.Get("Content-Type"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}
	// Save the file to the server
	var gpx models.Gpx
	gpx.UserID = user.ID
	gpx.Filename = file.Filename
	gpx.Title = title
	gpx.ID, err = models.SaveFile(c.MustGet("db").(*gorm.DB), gpx, file)
	if err != nil {
		if strings.Contains(err.Error(), "invalid GPX file: no coordinates found") {
			fmt.Println("Invalid GPX file: no coordinates found")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GPX file: no coordinates found"})
			return
		}
		fmt.Println("Error saving file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	if gpx.ID == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"file_id": gpx.ID})
}
