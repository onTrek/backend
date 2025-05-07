package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetActivity(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the activity ID from the URL parameters
	id := c.Param("id")
	activityId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting activity ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	// Retrieve the activity from the database
	activity, err := db.GetActivityByID(c.MustGet("db").(*sql.DB), activityId)
	if err != nil {
		fmt.Println("Error getting activity by ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	if activity.UserID != user.ID {
		fmt.Println("User is not authorized to access this activity")
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this activity"})
		return
	}

	// Return the activity as JSON
	c.JSON(200, activity)
}
