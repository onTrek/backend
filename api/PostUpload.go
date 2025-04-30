package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostUpload(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := utils.IsLogged(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the gpx file from the form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	// Check if the file is a GPX file
	if file.Header.Get("Content-Type") != "application/gpx+xml" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}
	// Save the file to the server
	var gpx utils.Gpx
	gpx.UserID = user.ID
	gpx.Filename = file.Filename
	gpx.Stats = ""
	err = db.SaveFile(c.MustGet("db").(*sql.DB), gpx, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file" + err.Error()})
		return
	}

	return
}
