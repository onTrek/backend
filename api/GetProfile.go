package api

import (
	"OnTrek/utils"
	"github.com/gin-gonic/gin"
)

// GetProfile godoc
// @Summary Get user profile by token
// @Description Fetches the profile information of the user based on the provided token in the Authorization header
// @Tags profile
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Success 200 {object} utils.UserInfo "User profile information"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /profile [get]
func GetProfile(c *gin.Context) {
	// Get the user from the context
	user := c.MustGet("user").(utils.User)

	// Get user profile from the database
	var userInfo utils.UserInfo

	userInfo.ID = user.ID
	userInfo.Username = user.Username
	userInfo.Email = user.Email

	c.JSON(200, userInfo)
}
