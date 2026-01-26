package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetFile godoc
// @Summary Download a GPX file by ID
// @Description Retrieves a file based on the provided file ID and authorization token
// @Tags gpx
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "File ID"
// @Success 200 {object} utils.Url "Returns the signed URL for the GPX file"
// @Failure 400 {object} utils.ErrorResponse "Invalid file ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Unauthorized access to file"
// @Failure 404 {object} utils.ErrorResponse "File not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /gpx/{id}/download [get]
func GetFile(c *gin.Context) {

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

	gpx, err := models.GetFileByID(c.MustGet("db").(*gorm.DB), fileID)
	if err != nil {
		fmt.Println("Error fetching file from database:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	downloadURL, err := utils.GenerateSignedURL(c.MustGet("storageConfig").(*utils.StorageConfig), gpx.StoragePath, utils.FileTypeGPX, "")
	if err != nil {
		c.JSON(500, gin.H{"error": "Error during file retrieval"})
		return
	}

	c.JSON(200, gin.H{
		"filename": gpx.Filename,
		"url":      downloadURL,
	})

}
