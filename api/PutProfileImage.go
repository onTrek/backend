package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"net/http"
	"strings"

	firebaseStorage "firebase.google.com/go/v4/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PutProfileImage godoc
// @Summary Upload a profile image
// @Description Uploads a profile image for the authenticated user. The image must be sent as form-data with the key 'file'. It must be an image file (jpg, jpeg, png) and should not exceed 5MB in size.
// @Tags profile
// @Accept image/*
// @Param Bearer header string true "Bearer token for user authentication"
// @Param file formData file true "Profile image to upload"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse "Invalid file"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to save file"
// @Router /profile/image [put]
func PutProfileImage(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	// Get the file from the form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// Check if the file is an image
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// Check if the file size is within the limit (5MB)
	if file.Size > 5*1024*1024 { // 5MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5MB limit"})
		return
	}

	// Get the file extension
	extension := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported image format"})
		return
	}

	err = models.UpdateExtension(c.MustGet("db").(*gorm.DB), user.ID, extension)
	if err != nil {
		fmt.Println("Error updating user extension:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	// Save the profile image
	_, err = utils.SaveFile(c.MustGet("firebaseStorage").(*firebaseStorage.Client), c.MustGet("storageConfig").(*utils.StorageConfig), file, "avatars", user.ID, extension)
	if err != nil {
		fmt.Println("Error saving profile image:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
