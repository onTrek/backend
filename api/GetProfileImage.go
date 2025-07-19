package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
)

// GetProfileImage godoc
// @Summary Get profile image of a user
// @Description Retrieve the profile image of a user by their ID.
// @Tags users
// @Produces image/*
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path string true "User ID to retrieve profile image for"
// @Success 200 {file} file "Profile image"
// @Failure 400 {object} utils.ErrorResponse "Invalid friend ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Profile image not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to save file"
// @Router /users/:id/image [get]
func GetProfileImage(c *gin.Context) {

	// Get the friend ID from the URL parameters
	friendID := c.Param("id")

	// Check if the friend ID is valid
	if friendID == "" {
		fmt.Println("Invalid friend ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	_, err := models.GetUserById(c.MustGet("db").(*gorm.DB), friendID)
	if err != nil {
		fmt.Println("Error getting user by ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Construct the file path
	filePath, err := utils.FindFileByID(friendID)
	if err != nil {
		fmt.Println("Error retrieving profile image:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile image"})
		return
	}

	// If the file path is empty, return a default image or an error
	if filePath == "" {
		fmt.Println("Profile image not found for user ID:", friendID)
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
