package functions

import (
	"OnTrek/db/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//token := strings.Split(c.Request.Header.Get("Authorization"), " ")
		token := c.GetHeader("Bearer")
		
		//if token[0] != "Bearer" {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token(Bearer token is required)"})
		//	c.Abort()
		//	return
		//}

		//if token[1] == "" {
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		//if strings.Contains(token[1], " ") {
		if strings.Contains(token, " ") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		//user, err := models.GetUserByToken(c.MustGet("db").(*gorm.DB), token[1])
		user, err := models.GetUserByToken(c.MustGet("db").(*gorm.DB), token)
		if err != nil {
			if err.Error() == "token expired" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			} else {
				fmt.Println("Error getting user by token:", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			}
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
