package api

import (
	"OnTrek/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProfile(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := utils.IsLogged(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user profile from the database
	var userInfo utils.UserInfo

	userInfo.ID = user.ID
	userInfo.Name = user.Name
	userInfo.Email = user.Email

	c.JSON(200, userInfo)
}
