package api

import (
	"OnTrek/db/functions"
	"OnTrek/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DeleteSessionId godoc
// @Summary Leave a session using session ID
// @Description Allows a user to leave a session by providing their session ID
// @Tags sessions
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path string true "Session ID"
// @Success 201 {object} utils.SuccessResponse "Successfully left session"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /sessions/{id}/members [delete]
func DeleteSessionId(c *gin.Context) {

	user := c.MustGet("user").(utils.User)

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

	err = functions.LeaveSessionById(c.MustGet("db").(*sql.DB), user.ID, sessionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found in session")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found in session"})
			return
		}
		fmt.Println("Error joining session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join session"})
		return
	}
	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully left session"})
}
