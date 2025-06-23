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

// DeleteDeclineFriendRequest godoc
// @Summary Delete a friend request
// @Description Delete a friend request by user ID
// @Tags friends
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for authentication"
// @Param id path string true "User ID of the friend request sender"
// @Success 200 {object} utils.SuccessResponse "Friend request declined successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Friend request not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to decline friend request"
// @Router /friends/requests/{id} [delete]
func DeleteDeclineFriendRequest(c *gin.Context) {
	// Get the user ID from the context
	user := c.MustGet("userId").(utils.User)

	// Get the user ID from the URL parameter
	userID := c.Param("id")

	// Validate the user ID
	if userID == "" {
		fmt.Println("Missing user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	// Check if the user ID is valid
	_, err := functions.GetUserById(c.MustGet("db").(*sql.DB), userID)
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
	err = functions.DeleteFriendRequest(c.MustGet("db").(*sql.DB), user.ID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Friend request not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request declined successfully"})
}
