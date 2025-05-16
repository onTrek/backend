package api

import (
	"OnTrek/db"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PostActivity godoc
// @Summary Create a new activity for the user
// @Description Allows a user to create a new activity by providing the details of the activity, including start time, end time, distance, and elevation data
// @Tags activity
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Param activity body utils.ActivityInput true "Activity input"
// @Success 201 {object} utils.SuccessResponse "Activity created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to save activity"
// @Router /activity/ [post]
func PostActivity(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the request body
	var activity utils.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate the request body
	if activity.StartTime == "" {
		fmt.Println("Missing required field: StartTime")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: StartTime"})
		return
	}
	if activity.EndTime == "" {
		fmt.Println("Missing required field: EndTime")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: EndTime"})
		return
	}
	if activity.Distance == 0 {
		fmt.Println("Missing required field: Distance")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: Distance"})
		return
	}
	if activity.TotalAscent == 0 {
		fmt.Println("Missing required field: TotalAscent")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: TotalAscent"})
		return
	}
	if activity.TotalDescent == 0 {
		fmt.Println("Missing required field: TotalDescent")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: TotalDescent"})
		return
	}
	if activity.StartingElevation == 0 {
		fmt.Println("Missing required field: StartingElevation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: StartingElevation"})
		return
	}
	if activity.MaximumElevation == 0 {
		fmt.Println("Missing required field: MaximumElevation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: MaximumElevation"})
		return
	}
	if activity.AverageSpeed == 0 {
		fmt.Println("Missing required field: AverageSpeed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: AverageSpeed"})
		return
	}

	if activity.AverageHeartRate == 0 {
		fmt.Println("Missing required field: AverageHeartRate")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: AverageHeartRate"})
		return
	}

	// Set the user ID
	activity.UserID = user.ID

	// Save the activity to the database
	err = db.SaveActivity(c.MustGet("db").(*sql.DB), activity)
	if err != nil {
		fmt.Println("Error saving activity:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save activity: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Activity created successfully"})
}
