package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetFiles(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get files from the database
	files, err := db.GetFiles(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		fmt.Println("Error fetching files:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching files"})
		return
	}

	// Return the files as JSON
	c.JSON(http.StatusOK, gin.H{"gpx_files": files})
}
