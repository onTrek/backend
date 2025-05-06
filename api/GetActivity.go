package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetActivity(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := utils.IsLogged(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the activity ID from the URL parameters
	id := c.Param("id")
	activityId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	// Retrieve the activity from the database
	activity, err := db.GetActivityByID(c.MustGet("db").(*sql.DB), activityId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	if activity.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this activity"})
		return
	}

	// Return the activity as JSON
	c.JSON(200, activity)
}
