package api

import (
	"OnTrek/db/functions"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// PutGroupId godoc
// @Summary Join a group using group ID
// @Description Allows a user to join a group by providing their group ID
// @Tags groups
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path string true "Group ID"
// @Success 201 {object} utils.SuccessResponse "Successfully joined group"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 400 {object} utils.ErrorResponse "Group ID is required"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Group not found"
// @Failure 409 {object} utils.ErrorResponse "User already joined the group"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /groups/{id}/members/ [put]
func PutGroupId(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.User)

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

	// Chekc if the group exists
	s, err := functions.CheckGroupExistsById(c.MustGet("db").(*sql.DB), groupId)
	if err != nil {
		fmt.Println("Error checking group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// Check if the group is valid
	if !s {
		fmt.Println("Group not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	err = functions.JoinGroupById(c.MustGet("db").(*sql.DB), user.ID, groupId)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			fmt.Println("User already joined the group")
			c.JSON(http.StatusConflict, gin.H{"error": "User already joined the group"})
			return
		}
		fmt.Println("Error joining group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join group"})
		return
	}
	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully joined group"})
}
