package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProfileImage godoc
// @Summary Get profile image of a user
// @Description Retrieve the profile image of a user by their ID.
// @Tags users
// @Produces image/*
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path string true "User ID to retrieve profile image for"
// @Success 200 {object} utils.Url "Returns the signed URL for the profile image"
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

	friend, err := models.GetUserExtension(c.MustGet("db").(*gorm.DB), friendID)
	if err != nil {
		fmt.Println("Error getting user by ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if friend.Extension == nil || *friend.Extension == "" {
		fmt.Println("Profile image not found for user:", friendID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile image not found"})
		return
	}

	downloadURL, err := utils.GenerateSignedURL(c.MustGet("storageConfig").(*utils.StorageConfig), friend.ID, utils.FileTypeAvatar, *friend.Extension)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error during file retrieval"})
		return
	}

	c.JSON(200, gin.H{
		"url": downloadURL,
	})

}
