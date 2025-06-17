package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PostLogin godoc
// @Summary Login a user
// @Description Authenticates a user using email and password, returns a token (user ID)
// @Tags auth
// @Accept json
// @Produce json
// @Param user body utils.Login true "User login credentials"
// @Success 200 {object} utils.UserToken "User ID token"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Invalid email or password"
// @Failure 500 {object} utils.ErrorResponse "Failed to login"
// @Router /auth/login [post]
func PostLogin(c *gin.Context) {
	// Get the request body
	var user utils.UserToken
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate the request body
	if input.Email == "" || input.Password == "" {
		fmt.Println("Missing required fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	database := c.MustGet("db").(*sql.DB)
	user, err := db.Login(database, input.Email, input.Password)
	if err != nil {
		if err.Error() == "user not found" {
			fmt.Println("User not found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		fmt.Println("Error logging in:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": user.Token})
}
