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

// PutGroupLocation godoc
// @Summary Update group information
// @Description Updates the user location(latitude, longitude, altitude, accuracy) and help request status in a group.
// @Tags groups
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "Group ID"
// @Param location body utils.GroupInfoUpdate true "Location data for the user in the group. GoingTo is optional."
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Group not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /groups/{id}/members/location [put]
func PutGroupLocation(c *gin.Context) {

	var groupInfo models.GroupMember

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
		Latitude    float64  `json:"latitude" binding:"required"`
		Longitude   float64  `json:"longitude" binding:"required"`
		Altitude    *float64 `json:"altitude" binding:"required"`
		Accuracy    float64  `json:"accuracy" binding:"required"`
		HelpRequest *bool    `json:"help_request" binding:"required"`
		GoingTo     *string  `json:"going_to" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Create a new group object
	groupInfo.GroupId = groupId
	groupInfo.Latitude = input.Latitude
	groupInfo.Longitude = input.Longitude
	groupInfo.Altitude = *input.Altitude
	groupInfo.HelpRequest = *input.HelpRequest
	groupInfo.Accuracy = input.Accuracy

	if input.GoingTo != nil {
		groupInfo.GoingTo = input.GoingTo
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

	err = models.UpdateGroupMember(c.MustGet("db").(*gorm.DB), user.ID, groupInfo)
	if err != nil {
		fmt.Println("Error updating group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)

}
