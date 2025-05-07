package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func DeleteActivity(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the activity ID from the URL parameter
	activityID := c.Param("id")
	activityIDInt, err := strconv.Atoi(activityID)
	if err != nil {
		fmt.Println("Error converting activity ID to int:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	activity, err := db.GetActivityByID(c.MustGet("db").(*sql.DB), activityIDInt)
	if err != nil {
		fmt.Println("Error getting activity by ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	if activity.UserID != user.ID {
		fmt.Println("User is not authorized to delete this activity")
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this activity"})
		return
	}

	err = db.DeleteActivity(c.MustGet("db").(*sql.DB), activityIDInt)
	if err != nil {
		fmt.Println("Error deleting activity:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})

}
