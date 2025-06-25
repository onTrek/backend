package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// PatchSessionGpx godoc
// @Summary Update group GPX file
// @Description Updates the GPX file for a group session.
// @Tags groups
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "Group ID"
// @Param file_id body utils.FileBody true "File ID of the GPX file to be used for the group"
// @Success 200 {object} utils.SuccessResponse "Group GPX updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden"
// @Failure 404 {object} utils.ErrorResponse "Group not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /groups/{id}/gpx [patch]
func PatchSessionGpx(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	// Get group ID from the URL
	group := c.Param("id")
	if group == "" {
		fmt.Println("Group ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group ID is required"})
		return
	}

	groupId, err := strconv.Atoi(group)
	if err != nil {
		fmt.Println("Error converting group ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get data from the request body
	var input struct {
		FileId int `json:"file_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if input.FileId <= 0 {
		fmt.Println("Invalid file ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Check if the group exists
	s, err := models.CheckGroupExistsByIdAndUserId(c.MustGet("db").(*gorm.DB), groupId, user.ID)
	if err != nil {
		fmt.Println("Error checking group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !s {
		fmt.Println("Group not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	leader, err := models.GetLeaderByGroup(c.MustGet("db").(*gorm.DB), groupId)
	if err != nil {
		fmt.Println("Error getting group leader:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if leader != user.ID {
		fmt.Println("User is not the group leader")
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the group leader"})
		return
	}

	file, err := models.GetFileByIDAndUserID(c.MustGet("db").(*gorm.DB), input.FileId, user.ID)
	if err != nil {
		fmt.Println("Error getting file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if file.UserID != user.ID {
		fmt.Println("File does not belong to the user")
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to use this file"})
		return
	}

	err = models.UpdateFileForTheGroup(c.MustGet("db").(*gorm.DB), groupId, input.FileId)
	if err != nil {
		fmt.Println("Error updating group file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group GPX updated successfully"})
}
