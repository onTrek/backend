package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostActivity(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := utils.IsLogged(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the request body
	var activity utils.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate the request body
	if activity.StartTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: StartTime"})
		return
	}
	if activity.EndTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: EndTime"})
		return
	}
	if activity.Distance == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: Distance"})
		return
	}
	if activity.TotalAscent == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: TotalAscent"})
		return
	}
	if activity.TotalDescent == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: TotalDescent"})
		return
	}
	if activity.StartingElevation == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: StartingElevation"})
		return
	}
	if activity.MaximumElevation == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: MaximumElevation"})
		return
	}
	if activity.AverageSpeed == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: AverageSpeed"})
		return
	}

	// Set the user ID
	activity.UserID = user.ID

	// Save the activity to the database
	err = db.SaveActivity(c.MustGet("db").(*sql.DB), activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save activity: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity created successfully"})
}
