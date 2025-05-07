package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteFriend(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the friend ID from the URL parameters
	friendID := c.Param("id")

	// Check if the friend ID is valid
	if friendID == "" {
		fmt.Println("Invalid friend ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	user2, err := db.GetUserById(c.MustGet("db").(*sql.DB), friendID)
	if err != nil {
		fmt.Println("Error getting user by ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Call the function to delete the friend
	err = db.DeleteFriend(c.MustGet("db").(*sql.DB), user.ID, user2.ID)
	if err != nil {
		if err == sql.ErrNoRows {
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
