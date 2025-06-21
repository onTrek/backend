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
// @Param Bearer header string true "Bearer token for user authentication"
// @Param session body utils.SessionInfoCreation true "Session information"
// @Success 201 {object} utils.SessionId "session_id"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /sessions/ [post]
func PostSession(c *gin.Context) {

	var sessionInfo utils.SessionInfo
	// Get token from the header
	token := c.GetHeader("Bearer")
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

	// Get data from the request body
	var input struct {
		Description string  `json:"description" binding:"required"`
		Latitude    float64 `json:"latitude" binding:"required"`
		Longitude   float64 `json:"longitude" binding:"required"`
		Altitude    float64 `json:"altitude" binding:"required"`
		Accuracy    float64 `json:"accuracy" binding:"required"`
		FileId      int     `json:"file_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if input.Description == "" {
		fmt.Println("Description is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description is required"})
		return
	}

	if input.FileId < 0 {
		fmt.Println("Invalid file ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	sessionInfo.Description = input.Description
	sessionInfo.Latitude = input.Latitude
	sessionInfo.Longitude = input.Longitude
	sessionInfo.Altitude = input.Altitude
	sessionInfo.Accuracy = input.Accuracy
	sessionInfo.FileId = input.FileId

	// Check if the file exists for the user
	_, err = db.GetFileByID(c.MustGet("db").(*sql.DB), sessionInfo.FileId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("File not found for user: " + user.Username)
			c.JSON(http.StatusBadRequest, gin.H{"error": "File not found for user: " + user.Username})
			return
		}
		fmt.Println("Error getting file by ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Create a new session
	session, err := db.CreateSession(c.MustGet("db").(*sql.DB), user, sessionInfo)
	if err != nil {
		fmt.Println("Error creating session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the session ID
	c.JSON(http.StatusCreated, gin.H{"session_id": session.ID})

}
