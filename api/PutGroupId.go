package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
// @Param user path string true "User ID"
// @Success 201 {object} utils.GroupMember "User successfully added to the group"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 400 {object} utils.ErrorResponse "Group ID is required"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Only group leader can add members"
// @Failure 404 {object} utils.ErrorResponse "Group not found"
// @Failure 409 {object} utils.ErrorResponse "User already joined the group"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /groups/{id}/members/{user} [put]
func PutGroupId(c *gin.Context) {
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

	// Chekc if the group exists
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

	userId := c.Param("user")
	if userId == "" {
		fmt.Println("User ID is required to add a member to the group")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required to add a member to the group"})
		return
	}

	user2, err := models.GetUserById(c.MustGet("db").(*gorm.DB), userId)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			fmt.Println("User not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	leader, err := models.GetLeaderByGroup(c.MustGet("db").(*gorm.DB), groupId)
	if err != nil {
		fmt.Println("Error getting group leader:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if user.ID != leader {
		fmt.Println("Only the group leader can add members")
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the group leader can add members"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var member utils.GroupMember

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return fmt.Errorf("error enabling foreign key enforcement: %v", err)
		}

		member, err = models.JoinGroup(tx, user2.ID, groupId)
		if err != nil {
			return fmt.Errorf("error joining group: %v", err)
		}

		return nil
	})
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

	member.Username = user2.Username

	// Return success response
	c.JSON(http.StatusCreated, member)
}
