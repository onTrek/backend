package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PostSession godoc
// @Summary Create a new session for the user
// @Description Creates a new session for the authenticated user with location data (latitude, longitude, altitude, accuracy)
// @Tags sessions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Param session body utils.SessionInfoUpdate true "Session information"
// @Success 201 {object} integer "session_id"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /sessions/ [post]
func PostSession(c *gin.Context) {

	var sessionInfo utils.SessionInfo
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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

	// Validate the request body
	if input.Latitude == 0 || input.Longitude == 0 || input.Altitude == 0 || input.Accuracy == 0 {
		fmt.Println("Missing required fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	sessionInfo.Latitude = input.Latitude
	sessionInfo.Longitude = input.Longitude
	sessionInfo.Altitude = input.Altitude
	sessionInfo.Accuracy = input.Accuracy

	// Create a new session
	session, err := db.CreateSession(c.MustGet("db").(*sql.DB), user.ID, sessionInfo)
	if err != nil {
		fmt.Println("Error creating session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return the session ID
	c.JSON(http.StatusCreated, gin.H{"session_id": session.ID})

}
