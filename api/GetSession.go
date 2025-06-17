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

// GetSession godoc
// @Summary Get session information
// @Description Get session information by session ID
// @Tags sessions
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Param id path int true "Session ID"
// @Success 200 {object} utils.SessionInfoResponseDoc "Session information"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 500 {object} utils.ErrorResponse "Error fetching files"
// @Router /sessions/{id} [get]
func GetSession(c *gin.Context) {

	var sessionInfo utils.SessionInfoResponse
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

	// Get session info
	sessionInfo, err = db.GetSessionInfoMember(c.MustGet("db").(*sql.DB), sessionId, user.ID)
	if err != nil {
		fmt.Println("Error getting session info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, sessionInfo)
}
