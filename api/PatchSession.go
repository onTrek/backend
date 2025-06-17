package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// PatchSession godoc
// @Summary Close a session
// @Description Closes an active session for the user who is the leader of the session
// @Tags sessions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Param id path int true "Session ID"
// @Success 200 {object} utils.SuccessResponse "Session closed successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid session ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /sessions/{id} [patch]
func PatchSession(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserByToken(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		if err.Error() == "token expired" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get session ID from the URL
	session := c.Param("id")
	if session == "" {
		fmt.Println("Session ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	sessionId, err := strconv.Atoi(session)
	if err != nil {
		fmt.Println("Error converting session ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	leader, err := db.GetLeaderBySession(c.MustGet("db").(*sql.DB), sessionId)
	if err != nil {
		fmt.Println("Error getting leader by session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if leader != user.ID {
		fmt.Println("User is not the leader of the session")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	err = db.CloseSession(c.MustGet("db").(*sql.DB), sessionId)
	if err != nil {
		fmt.Println("Error closing session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Session closed successfully"})
}
