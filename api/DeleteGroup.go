package api

import (
	"OnTrek/db/functions"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DeleteGroup godoc
// @Summary Delete a group
// @Description Deletes a group based on the provided group ID from the URL and the user's token from the header
// @Tags groups
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path string true "Group ID to be deleted" example:"12345"
// @Success 200 {object} utils.SuccessResponse "Group deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid group ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - User is not the leader of the group"
// @Failure 404 {object} utils.ErrorResponse "Group not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to delete group"
// @Router /groups/{id} [delete]
func DeleteGroup(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.User)

	// Get group ID from the URL
	group := c.Param("id")
	groupId, err := strconv.Atoi(group)
	if err != nil {
		fmt.Println("Error converting group ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	if groupId < 0 {
		fmt.Println("Group ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group ID is required"})
		return
	}

	// Check if the group exists
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

	leader, err := functions.GetLeaderByGroup(c.MustGet("db").(*sql.DB), groupId)
	if err != nil {
		fmt.Println("Error getting leader by group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve group leader"})
		return
	}

	// Check if the user is the leader of the group
	if user.ID != leader {
		fmt.Println("User is not the leader of the group")
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this group"})
		return
	}

	// Call the database function to delete the group
	err = functions.DeleteGroupById(c.MustGet("db").(*sql.DB), user.ID, groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}
