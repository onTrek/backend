package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteProfile godoc
// @Summary Delete user profile
// @Description Deletes the user profile based on the provided authorization token.
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Success 200 {object} utils.SuccessResponse "User deleted successfully"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to delete user"
// @Router /profile [delete]
func DeleteProfile(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Delete user from the database
	err = db.DeleteUser(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		fmt.Println("Error deleting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
