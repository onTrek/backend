package api

import (
	"OnTrek/db"
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
// @Param Authorization header string true "Bearer token for user authentication"
// @Param id path int true "Session ID"
// @Param session body utils.SessionInfoUpdate true "Session information"
// @Success 200 {object} utils.SuccessResponse "Session updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /sessions/{id} [put]
func PutSession(c *gin.Context) {

	var sessionInfo utils.SessionInfo
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
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

	// Get data from the request body
	var input struct {
		Latitude    float64 `json:"latitude" binding:"required"`
		Longitude   float64 `json:"longitude" binding:"required"`
		Altitude    float64 `json:"altitude" binding:"required"`
		HelpRequest *bool   `json:"help_request" binding:"required"`
		Accuracy    float64 `json:"accuracy" binding:"required"`
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

	// Check if the session exists
	_, err = db.CheckSessionExistsByIdAndUserId(c.MustGet("db").(*sql.DB), sessionId, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Session not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		} else {
			fmt.Println("Error checking session:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	err = db.UpdateSession(c.MustGet("db").(*sql.DB), user.ID, sessionInfo)
	if err != nil {
		fmt.Println("Error updating session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session updated successfully"})

}
