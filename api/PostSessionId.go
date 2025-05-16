package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// PostSessionId godoc
// @Summary Join a session using session ID and location data
// @Description Allows a user to join a session by providing their session ID and location data (latitude, longitude, altitude, accuracy)
// @Tags sessions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Param id path string true "Session ID"
// @Param session body utils.SessionInfoJoin true "Session information"
// @Success 201 {object} utils.SuccessResponse "Successfully joined session"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 400 {object} utils.ErrorResponse "Session ID is required"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 409 {object} utils.ErrorResponse "User already joined the session"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /sessions/{id} [post]
func PostSessionId(c *gin.Context) {
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
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
		Altitude  float64 `json:"altitude" binding:"required"`
		Accuracy  float64 `json:"accuracy" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sessionInfo.SessionID = sessionId
	sessionInfo.Latitude = input.Latitude
	sessionInfo.Longitude = input.Longitude
	sessionInfo.Altitude = input.Altitude
	sessionInfo.Accuracy = input.Accuracy

	// Chekc if the session exists
	s, err := db.CheckSessionExistsById(c.MustGet("db").(*sql.DB), sessionId)
	if err != nil {
		fmt.Println("Error checking session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// Check if the session is valid
	if s.ID == -1 {
		fmt.Println("Session not found or expired")
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found or expired"})
		return
	}

	err = db.JoinSession(c.MustGet("db").(*sql.DB), user.ID, sessionInfo)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			fmt.Println("User already joined the session")
			c.JSON(http.StatusConflict, gin.H{"error": "User already joined the session"})
			return
		}
		fmt.Println("Error joining session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join session"})
		return
	}
	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully joined session"})
}
