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

	s, err := functions.CheckGroupExistsById(c.MustGet("db").(*sql.DB), groupId)
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

	groupInfo, err = functions.GetGroupInfo(c.MustGet("db").(*sql.DB), groupId)
	if err != nil {
		fmt.Println("Error getting group info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, groupInfo)
}
