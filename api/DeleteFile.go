package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func DeleteFile(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	file := c.Param("id")
	// Validate the file ID
	fileID, err := strconv.Atoi(file)
	if err != nil {
		fmt.Println("Error converting file ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Fetch the file from the database
	gpx, err := db.GetFileByID(c.MustGet("db").(*sql.DB), fileID, user.ID)
	if err != nil {
		fmt.Println("Error fetching file from database:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Delete file from the disk
	err = utils.DeleteFile(gpx.StoragePath)
	if err != nil {
		fmt.Println("Error deleting file from disk:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Delete file from the database
	err = db.DeleteFileById(c.MustGet("db").(*sql.DB), fileID, user.ID)
	if err != nil {
		fmt.Println("Error deleting file from database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
	// Return success response

}
