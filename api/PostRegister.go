package api

import (
	"OnTrek/db/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// PostRegister godoc
// @Summary Register a new user
// @Description Registers a new user with an email, password, and name.
// @Tags auth
// @Accept json
// @Produce json
// @Param register body utils.RegisterInput true "User registration credentials. Email must be unique. Password must be at least 8 characters long."
// @Success 201 {object} utils.SuccessResponse "User registered successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 409 {object} utils.ErrorResponse "User with this email already exists"
// @Failure 500 {object} utils.ErrorResponse "Failed to register user"
// @Router /auth/register [post]
func PostRegister(c *gin.Context) {
	// Get the request body
	var user models.User
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate the request body
	if input.Email == "" || input.Password == "" || input.Username == "" {
		fmt.Println("Missing required fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	user.Email = input.Email
	user.Username = input.Username
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedPassword)

	err := models.RegisterUser(c.MustGet("db").(*gorm.DB), user)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			fmt.Println("Email already exists")
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		fmt.Println("Error registering user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
