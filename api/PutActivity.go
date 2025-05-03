package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PutActivity(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := utils.IsLogged(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the activity ID from the URL
	activityID := c.Param("id")
	activityIDInt, err := strconv.Atoi(activityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID: " + activityID})
		return
	}

	// Bind the JSON body to an Activity struct
	var activity utils.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Get userId from activity
	userID, err := db.GetUserIdByActivity(c.MustGet("db").(*sql.DB), activityIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	// Check if the user ID from the token matches the user ID from the activity
	if user.ID != userID || activity.ID != activityIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this activity"})
		return
	}

	// Set the user ID for the activity
	activity.UserID = user.ID

	// Update the activity in the database
	err = db.UpdateActivity(c.MustGet("db").(*sql.DB), activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity updated successfully"})
}
