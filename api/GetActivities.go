package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetActivities godoc
// @Summary Get a user's activities
// @Description Fetches a list of activities for a user based on their authorization token
// @Tags activity
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Success 200 {object} []utils.Activity "List of activities"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to get activities"
// @Router /activity/ [get]
func GetActivities(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get activities from the database
	activities, err := db.GetActivities(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		fmt.Println("Error getting activities:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activities: " + err.Error()})
		return
	}

	// Return the activities as JSON
	c.JSON(http.StatusOK, gin.H{"activities": activities})
}
