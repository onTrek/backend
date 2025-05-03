package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStats(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := utils.IsLogged(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Calculate global stats
	var globalStats utils.GlobalStats
	globalStats, err = db.CalculateGlobalStats(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve global stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": globalStats})
}
