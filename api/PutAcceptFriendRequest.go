package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PutAcceptFriendRequest godoc
// @Summary Accept a friend request
// @Description Accept a friend request from another user
// @Tags friends
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for authentication"
// @Param id path string true "User ID of the friend request sender"
// @Success 201 {object} utils.SuccessResponse "Friend request accepted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Friend request not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to accept friend request"
// @Router /friends/requests/{id} [put]
func PutAcceptFriendRequest(c *gin.Context) {
	// Get the user ID from the context
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
	_, err := db.GetUserById(c.MustGet("db").(*sql.DB), userID)
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

	// Call the service to accept the friend request
	err = db.AcceptFriendRequest(c.MustGet("db").(*sql.DB), user.ID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Friend request not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}
