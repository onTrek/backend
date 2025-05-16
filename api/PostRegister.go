package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// PostRegister godoc
// @Summary Register a new user
// @Description Registers a new user with an email, password, and name
// @Tags auth
// @Accept json
// @Produce json
// @Param register body utils.RegisterInput true "User registration credentials"
// @Success 201 {object} utils.SuccessResponse "User registered successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 409 {object} utils.ErrorResponse "Email already exists"
// @Failure 500 {object} utils.ErrorResponse "Failed to register user"
// @Router /auth/register [post]
func PostRegister(c *gin.Context) {
	// Get the request body
	var user utils.User
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
	user.Password = string(hashedPassword)
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now().Format(time.RFC3339)

	err := db.RegisterUser(c.MustGet("db").(*sql.DB), user)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			fmt.Println("Email already exists")
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		fmt.Println("Error registering user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
