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
	"strconv"
)

// DeleteLeaveRemoveMember godoc
// @Summary Leave a group using group ID
// @Description Allows a user to leave a group by providing their group ID or remove a user from the group if the user is the leader.
// @Tags groups
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path string true "Group ID"
// @Param user_id query string false "If provided, the user with this ID will be removed from the group if the authenticated user is the leader"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Group not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /groups/{id}/members [delete]
func DeleteLeaveRemoveMember(c *gin.Context) {

	user := c.MustGet("user").(utils.UserInfo)

	removeUser := true

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

	// Get user ID from query parameters for remove user from group whener user id from token is the leader of the group
	userId := c.Query("user_id")
	if userId == "" {
		removeUser = false
	}

	// Check if the group exists
	s, err := models.CheckGroupExistsById(c.MustGet("db").(*gorm.DB), groupId)
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

	if removeUser {
		// Check if the user is the leader of the group
		leader, err := models.GetLeaderByGroup(c.MustGet("db").(*gorm.DB), groupId)
		if err != nil {
			fmt.Println("Error getting group leader:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if leader != user.ID {
			fmt.Println("User is not the leader of the group")
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not the leader of the group"})
			return
		}

		err = models.LeaveGroupById(c.MustGet("db").(*gorm.DB), userId, groupId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Println("User not found in group")
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found in group"})
				return
			}
			fmt.Println("Error joining group:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join group"})
			return
		}
		// Return success response
		c.JSON(http.StatusNoContent, nil)
	} else {
		// Check if the user is the leader of the group
		leader, err := models.GetLeaderByGroup(c.MustGet("db").(*gorm.DB), groupId)
		if err != nil {
			fmt.Println("Error getting group leader:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if leader == user.ID {
			fmt.Println("Leader cannot leave the group")
			c.JSON(http.StatusForbidden, gin.H{"error": "The leader cannot abandon a group he owns"})
			return
		}
		
		err = models.LeaveGroupById(c.MustGet("db").(*gorm.DB), user.ID, groupId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Println("User not found in group")
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found in group"})
				return
			}
			fmt.Println("Error joining group:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join group"})
			return
		}
		// Return success response
		c.JSON(http.StatusNoContent, nil)
	}
	return
}
