package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetActivities(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := utils.IsLogged(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get activities from the database
	activities, err := db.GetActivities(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activities: " + err.Error()})
		return
	}

	// Return the activities as JSON
	c.JSON(http.StatusOK, gin.H{"activities": activities})
}
