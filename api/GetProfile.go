package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetProfile godoc
// @Summary Get user profile by token
// @Description Fetches the profile information of the user based on the provided token in the Authorization header
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Success 200 {object} utils.UserInfo "User profile information"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /profile [get]
func GetProfile(c *gin.Context) {
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

	// Get user profile from the database
	var userInfo utils.UserInfo

	userInfo.ID = user.ID
	userInfo.Username = user.Username
	userInfo.Email = user.Email

	c.JSON(200, userInfo)
}
