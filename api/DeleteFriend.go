package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// DeleteFriend godoc
// @Summary Delete a friend from the user's friend list
// @Description Deletes a friend based on the provided friend ID from the URL and the user's token from the header
// @Tags friends
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Param id path string true "Friend ID to be deleted" example:"12345"
// @Success 200 {object} utils.SuccessResponse "Friend deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid friend ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 404 {object} utils.ErrorResponse "Friend not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to delete friend"
// @Router /friends/{id} [delete]
func DeleteFriend(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.UserInfo)

	// Get the friend ID from the URL parameters
	friendID := c.Param("id")

	// Check if the friend ID is valid
	if friendID == "" {
		fmt.Println("Invalid friend ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	user2, err := models.GetUserById(c.MustGet("db").(*gorm.DB), friendID)
	if err != nil {
		fmt.Println("Error getting user by ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = models.DeleteFriend(c.MustGet("db").(*gorm.DB), user.ID, user2.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Friend not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
			return
		}
		fmt.Println("Error deleting friend:", err)
		c.JSON(500, gin.H{"error": "Failed to delete friend"})
		return
	}

	c.JSON(200, gin.H{"message": "Friend deleted successfully"})
}
