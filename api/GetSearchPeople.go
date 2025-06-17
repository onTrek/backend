package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSearchPeople godoc
// @Summary Search for users by username
// @Description Searches for users by username
// @Tags search
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Param query query string true "Search query"
// @Success 200 {array} utils.UserEssentials "Returns a list of users matching the search query"
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "No users found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /search [get]
func GetSearchPeople(c *gin.Context) {

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

	// Get the search query from the request
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	// Fetch users matching the search query from the database
	users, err := db.SearchUsers(c.MustGet("db").(*sql.DB), query, user.ID)
	if err != nil {
		fmt.Println("Error searching users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	// Check if any users were found
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No users found"})
		return
	}
	// Return the list of users matching the search query
	c.JSON(http.StatusOK, users)

}
