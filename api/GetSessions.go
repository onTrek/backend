package api

import (
	"OnTrek/db/functions"
	"OnTrek/utils"
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
// @Param Bearer header string true "Bearer token for user authentication"
// @Success 200 {object} []utils.SessionDoc "List of sessions"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Error fetching files"
// @Router /sessions/ [get]
func GetSessions(c *gin.Context) {

	// Get the user from the context
	user := c.MustGet("user").(utils.User)

	sessions, err := functions.GetSessionsByUserId(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		fmt.Println("Error getting sessions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if len(sessions) == 0 {
		c.JSON(http.StatusOK, []utils.SessionDoc{})
		return
	}

	c.JSON(http.StatusOK, sessions)
}
