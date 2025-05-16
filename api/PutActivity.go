package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// PutActivity godoc
// @Summary Update an existing activity
// @Description Update an existing activity by ID
// @Tags activity
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Param id path int true "Activity ID"
// @Success 200 {object} utils.SuccessResponse "Activity ended successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden"
// @Failure 404 {object} utils.ErrorResponse "Activity not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to update activity"
// @Router /activity/{id} [put]
func PutActivity(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the activity ID from the URL
	activityID := c.Param("id")
	activityIDInt, err := strconv.Atoi(activityID)
	if err != nil {
		fmt.Println("Error converting activity ID to int:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID: " + activityID})
		return
	}

	// Bind the JSON body to an Activity struct
	var activity utils.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		fmt.Println("Error binding JSON to activity struct:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Get userId from activity
	userID, err := db.GetUserIdByActivity(c.MustGet("db").(*sql.DB), activityIDInt)
	if err != nil {
		fmt.Println("Error getting user ID by activity ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	// Check if the user ID from the token matches the user ID from the activity
	if user.ID != userID || activity.ID != activityIDInt {
		fmt.Println("User ID does not match the activity owner")
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this activity"})
		return
	}

	// Set the user ID for the activity
	activity.UserID = user.ID

	// Update the activity in the database
	err = db.UpdateActivity(c.MustGet("db").(*sql.DB), activity)
	if err != nil {
		fmt.Println("Error updating activity:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity updated successfully"})
}
