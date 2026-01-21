package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PatchFilePrivacy godoc
// @Summary Change GPX visibility
// @Description Set is_public to true or false
// @Tags gpx
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "Gpx ID"
// @Param input body utils.PrivacyUpdateInput true "Privacy update input"
// @Success 200 {object} utils.SuccessResponse "File privacy updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden"
// @Failure 404 {object} utils.ErrorResponse "Group not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /gpx/{id}/privacy [patch]
func PatchFilePrivacy(c *gin.Context) {
	// 1. Setup e Auth
	user := c.MustGet("user").(utils.UserInfo)

	// Get the file ID from the URL parameter
	file := c.Param("id")
	// Validate the file ID
	fileID, err := strconv.Atoi(file)
	if err != nil {
		fmt.Println("Error converting file ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// 3. Parsing del Body (JSON)
	var input utils.PrivacyUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, expected JSON with 'is_public' boolean field"})
		return
	}

	permission, err := models.CheckFilePermissions(c.MustGet("db").(*gorm.DB), fileID, user.ID)
	if err != nil {
		fmt.Println("Error checking file permissions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking file permissions"})
		return
	}
	if !permission {
		fmt.Println("Unauthorized access to file:", fileID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	if input.IsPublic != nil {
		err = models.UpdateFilePrivacy(c.MustGet("db").(*gorm.DB), fileID, *input.IsPublic)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println("File not found:", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
				return
			}
			fmt.Println("Error updating file privacy:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating file privacy"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File privacy updated successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, 'is_public' field is required"})
	}
}
