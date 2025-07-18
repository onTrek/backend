package api

import (
	"OnTrek/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

// GetProfileImage godoc
// @Summary Get profile image
// @Description Retrieves the profile image for the authenticated user.
// @Tags profile
// @Produces image/*
// @Param Bearer header string true "Bearer token for user authentication"
// @Success 200 {file} file "Profile image"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Profile image not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to save file"
// @Router /profile/image [get]
func GetProfileImage(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	// Construct the file path
	filePath, err := utils.FindFileByID(user.ID)
	if err != nil {
		fmt.Println("Error retrieving profile image:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile image"})
		return
	}

	// If the file path is empty, return a default image or an error
	if filePath == "" {
		fmt.Println("Profile image not found for user ID:", user.ID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile image not found"})
		return
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Profile image file does not exist:", filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile image not found"})
		return
	}

	// Get file extension
	fileExt := filePath[strings.LastIndex(filePath, "."):]

	// Set the content type and attachment header
	c.Header("Content-Type", "image/"+fileExt[1:])
	c.Header("Content-Disposition", "attachment; filename="+filePath[strings.LastIndex(filePath, "/")+1:])
	// Serve the profile image
	c.File(filePath)
}
