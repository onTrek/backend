package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetStats godoc
// @Summary Get user stats
// @Description Retrieves global statistics for the user based on their token
// @Tags stats
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Success 200 {object} utils.GlobalStats "Global statistics for the user"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /stats [get]
func GetStats(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Calculate global stats
	var globalStats utils.GlobalStats
	globalStats, err = db.CalculateGlobalStats(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		fmt.Println("Error calculating global stats:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve global stats: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": globalStats})
}
