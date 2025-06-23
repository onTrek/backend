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

// PutSession godoc
// @Summary Update session information
// @Description Updates the session with new location data (latitude, longitude, altitude, accuracy)
// @Tags sessions
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "Session ID"
// @Param session body utils.SessionInfoUpdate true "Session information"
// @Success 200 {object} utils.SuccessResponse "Session updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /sessions/{id}/members/ [put]
func PutSession(c *gin.Context) {

	var sessionInfo utils.SessionInfo

	// Get the user from the context
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

	// Get data from the request body
	var input struct {
		Latitude    float64 `json:"latitude" binding:"required"`
		Longitude   float64 `json:"longitude" binding:"required"`
		Altitude    float64 `json:"altitude" binding:"required"`
		Accuracy    float64 `json:"accuracy" binding:"required"`
		HelpRequest *bool   `json:"help_request" binding:"required"`
		GoingTo     string  `json:"going_to" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Create a new session object
	sessionInfo.SessionID = sessionId
	sessionInfo.Latitude = input.Latitude
	sessionInfo.Longitude = input.Longitude
	sessionInfo.Altitude = input.Altitude
	sessionInfo.HelpRequest = *input.HelpRequest
	sessionInfo.Accuracy = input.Accuracy
	sessionInfo.GoingTo = input.GoingTo

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

	err = functions.UpdateSession(c.MustGet("db").(*sql.DB), user.ID, sessionInfo)
	if err != nil {
		fmt.Println("Error updating session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session updated successfully"})

}
