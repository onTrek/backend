package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetFriends godoc
// @Summary Get list of friends for a user
// @Description Retrieves the list of friends for the authenticated user based on the token
// @Tags friends
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Success 200 {array} utils.UserEssentials "List of friends"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "No friends found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /friends/ [get]
func GetFriends(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch friends from the database
	friends, err := db.GetFriends(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		fmt.Println("Error fetching friends:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch friends"})
		return
	}

	// Check if the user has any friends
	if len(friends) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No friends found"})
		return
	}

	// Return the list of friends
	c.JSON(http.StatusOK, friends)
}
