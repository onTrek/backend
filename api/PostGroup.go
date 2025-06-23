package api

import (
	"OnTrek/db/functions"
	"OnTrek/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PostGroup godoc
// @Summary Create a new group
// @Description Creates a new group for the authenticated user
// @Tags groups
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param group body utils.GroupInfoCreation true "Group information"
// @Success 201 {object} utils.GroupId "group_id"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /groups/ [post]
func PostGroup(c *gin.Context) {

	var groupInfo utils.GroupInfo

	// Get the user from the context
	user := c.MustGet("user").(utils.User)

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

	groupInfo.Description = input.Description
	// If FileId is provided, check if it exists
	if input.FileId != nil {
		groupInfo.FileId = *input.FileId

		if groupInfo.FileId < 0 {
			fmt.Println("Invalid file ID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
			return
		}

		_, err := functions.GetFileByID(c.MustGet("db").(*sql.DB), groupInfo.FileId)
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
	} else {
		groupInfo.FileId = -1
	}

	// Create a new group
	group, err := functions.CreateGroup(c.MustGet("db").(*sql.DB), user, groupInfo)
	if err != nil {
		fmt.Println("Error creating group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the session ID
	c.JSON(http.StatusCreated, gin.H{"group_id": group.ID})

}
