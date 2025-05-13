package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSessions godoc
// @Summary Get sessions by user ID
// @Description Get all sessions for a user by their ID
// @Tags sessions
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Success 200 {object} []utils.SessionDoc "List of sessions"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 500 {object} utils.ErrorResponse "Error fetching files"
// @Router /sessions/ [get]
func GetSessions(c *gin.Context) {

	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sessions, err := db.GetSessionsByUserId(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No sessions found"})
			return
		}
		fmt.Println("Error getting sessions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}
