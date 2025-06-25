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

// GetGroup godoc
// @Summary Get group information by group ID
// @Description Get group information by group ID
// @Tags groups
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "Group ID"
// @Success 200 {object} utils.GroupInfoResponseDoc "Group information"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 404 {object} utils.ErrorResponse "Group not found"
// @Failure 500 {object} utils.ErrorResponse "Error fetching files"
// @Router /groups/{id} [get]
func GetGroup(c *gin.Context) {

	user := c.MustGet("user").(utils.UserInfo)

	var groupInfo utils.GroupInfoResponse

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

	groupInfo, err = models.GetGroupInfo(c.MustGet("db").(*gorm.DB), groupId)
	if err != nil {
		fmt.Println("Error getting group info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, groupInfo)
}
