package api

import (
	"OnTrek/db/functions"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DeleteSession godoc
// @Summary Delete a session
// @Description Deletes a session based on the provided session ID from the URL and the user's token from the header
// @Tags sessions
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path string true "Session ID to be deleted" example:"12345"
// @Success 200 {object} utils.SuccessResponse "Session deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid session ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - User is not the leader of the session"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to delete session"
// @Router /sessions/{id} [delete]
func DeleteSession(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.User)

	// Get session ID from the URL
	session := c.Param("id")
	sessionId, err := strconv.Atoi(session)
	if err != nil {
		fmt.Println("Error converting session ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}
	if sessionId < 0 {
		fmt.Println("Session ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Chekc if the session exists
	s, err := functions.CheckSessionExistsById(c.MustGet("db").(*sql.DB), sessionId)
	if err != nil {
		fmt.Println("Error checking session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// Check if the session is valid
	if s.ID == -1 {
		fmt.Println("Session not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	leader, err := functions.GetLeaderBySession(c.MustGet("db").(*sql.DB), sessionId)
	if err != nil {
		fmt.Println("Error getting leader by session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session leader"})
		return
	}

	// Check if the user is the leader of the session
	if user.ID != leader {
		fmt.Println("User is not the leader of the session")
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this session"})
		return
	}

	// Call the database function to delete the session
	err = functions.DeleteSessionById(c.MustGet("db").(*sql.DB), user.ID, sessionId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete session"})
		return
	}

	// Return a success response
	c.JSON(200, gin.H{"message": "Session deleted successfully"})
}
