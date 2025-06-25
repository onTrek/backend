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

// GetMembersInfo godoc
// @Summary Get members information of a group
// @Description Get members information by group ID
// @Tags groups
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "Group ID"
// @Success 200 {array} utils.MemberInfo "List of members in the group"
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Group not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /groups/{id}/members/ [get]
func GetMembersInfo(c *gin.Context) {

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

	members, err := models.GetMembersInfoByGroupId(c.MustGet("db").(*gorm.DB), groupId)
	if err != nil {
		fmt.Println("Error getting groups:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if len(members) == 0 {
		c.JSON(http.StatusOK, []utils.MemberInfo{})
		return
	}

	c.JSON(http.StatusOK, members)
}
