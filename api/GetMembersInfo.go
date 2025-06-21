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

// GetMembersInfo godoc
// @Summary Get members information of a session
// @Description Get members information by session ID
// @Tags sessions
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path int true "Session ID"
// @Success 200 {array} utils.MemberInfo "List of members in the session"
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Session not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /sessions/{id}/members/ [get]
func GetMembersInfo(c *gin.Context) {

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
	s, err := db.CheckSessionExistsById(c.MustGet("db").(*sql.DB), sessionId)
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

	members, err := db.GetMembersInfoBySessionId(c.MustGet("db").(*sql.DB), sessionId)
	if err != nil {
		fmt.Println("Error getting sessions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if len(members) == 0 {
		c.JSON(http.StatusOK, []utils.MemberInfo{})
		return
	}

	c.JSON(http.StatusOK, members)
}
