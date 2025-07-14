package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// PostGroup godoc
// @Summary Create a new group
// @Description Creates a new group for the authenticated user
// @Tags groups
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param group body utils.GroupInfoCreation true "Group information. Fields: description (required, max 64 characters), file_id (optional, must be a valid file ID)"
// @Success 201 {object} utils.GroupId "group_id"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /groups/ [post]
func PostGroup(c *gin.Context) {

	var groupInfo models.Group

	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	// Get data from the request body
	var input struct {
		Description string `json:"description" binding:"required"`
		FileId      *int   `json:"file_id" binding:"omitempty,numeric"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if input.Description == "" {
		fmt.Println("Description is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description is required"})
		return
	}

	if len(input.Description) > 64 {
		fmt.Println("Description is too long(64 characters max)")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description is too long(64 characters max)"})
		return
	}

	groupInfo.Description = input.Description
	groupInfo.CreatedBy = user.ID
	// If FileId is provided, check if it exists
	if input.FileId != nil {

		if *input.FileId < 0 {
			fmt.Println("Invalid file ID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
			return
		}

		groupInfo.FileId = input.FileId

		_, err := models.GetFileByID(c.MustGet("db").(*gorm.DB), *groupInfo.FileId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Println("File not found for user: " + user.Username)
				c.JSON(http.StatusBadRequest, gin.H{"error": "File not found for user: " + user.Username})
				return
			}
			fmt.Println("Error getting file by ID:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}

	groupId, err := models.CreateGroup(c.MustGet("db").(*gorm.DB), groupInfo)
	if err != nil {
		fmt.Println("Error creating group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the session ID
	c.JSON(http.StatusCreated, gin.H{"group_id": groupId})

}
