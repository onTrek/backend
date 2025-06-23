package api

import (
	"OnTrek/db/functions"
	"OnTrek/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PostAddFriendRequest godoc
// @Summary Send a friend request
// @Description Allows a user to add another user as a friend by their user ID
// @Tags friends
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path string true "User ID of the friend to be added"
// @Success 201 {object} utils.SuccessResponse "Friend added successfully"
// @Failure 400 {object} utils.ErrorResponse "Missing or invalid user ID"ss
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 409 {object} utils.ErrorResponse "Users are already friends"
// @Failure 409 {object} utils.ErrorResponse "User cannot add themselves as a friend"
// @Failure 500 {object} utils.ErrorResponse "Failed to add friend"
// @Router /friends/requests/{id} [post]
func PostAddFriendRequest(c *gin.Context) {

	// Get the user from the context
	user := c.MustGet("user").(utils.User)

	// Get the user ID from the URL parameter
	userID := c.Param("id")

	// Validate the user ID
	if userID == "" {
		fmt.Println("Missing user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	// Check if the user ID is valid
	user2, err := functions.GetUserById(c.MustGet("db").(*sql.DB), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		fmt.Println("Error getting user by ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the user is trying to add themselves as a friend
	if user2.ID == user.ID {
		fmt.Println("User cannot add themselves as a friend")
		c.JSON(http.StatusConflict, gin.H{"error": "You cannot add yourself as a friend"})
		return
	}

	// Add the friend to the database
	err = functions.AddFriend(c.MustGet("db").(*sql.DB), user.ID, user2.ID)
	if err != nil {
		if err.Error() == "Users are already friends" {
			fmt.Println("Users are already friends")
			c.JSON(http.StatusConflict, gin.H{"error": "Users are already friends"})
			return
		}
		fmt.Println("Error adding friend:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add friend"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Friend request sent successfully"})
}
