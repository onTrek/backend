package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSearchUsers godoc
// @Summary Search for users by username
// @Description Searches for users by username
// @Tags search
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param username query string true "Search for username"
// @Param friendsOnly query bool false "Search for friends only (optional, true/false, default is false)"
// @Success 200 {array} utils.UserSearchResponse "Returns a list of users matching the search query ordered by username. State is set to -1 if the user is not a friend, 0 if there is a request sent, and 1 if the user is a friend."
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "No users found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /search/users/ [get]
func GetSearchUsers(c *gin.Context) {

	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	// Get the search query from the request
	query := c.Query("username")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	friends := c.Query("friendsOnly")
	if friends == "" {
		friends = "false"
	}

	friendsValue, err := strconv.ParseBool(friends)
	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'friends' parameter"})
	}

	// Fetch users matching the search query from the database
	users, err := models.SearchUsers(c.MustGet("db").(*gorm.DB), query, friendsValue, user.ID)
	if err != nil {
		fmt.Println("Error searching users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	// Check if any users were found
	if len(users) == 0 {
		c.JSON(http.StatusOK, []utils.UserEssentials{})
		return
	}
	// Return the list of users matching the search query
	c.JSON(http.StatusOK, users)

}
