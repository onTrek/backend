package api

import (
	"OnTrek/db"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DeleteActivity godoc
// @Summary Delete an activity by ID
// @Description Deletes an activity after verifying the user's token and authorization to perform the action
// @Tags activity
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for user authentication"
// @Param id path int true "Activity ID"
// @Success 200 {object} utils.SuccessResponse "Activity deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid activity ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden"
// @Failure 404 {object} utils.ErrorResponse "Activity not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to delete activity"
// @Router /activity/{id} [delete]
func DeleteActivity(c *gin.Context) {
	// Get token from the header
	token := c.GetHeader("Authorization")
	user, err := db.GetUserById(c.MustGet("db").(*sql.DB), token)
	if err != nil {
		fmt.Println("Error getting user by token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the activity ID from the URL parameter
	activityID := c.Param("id")
	activityIDInt, err := strconv.Atoi(activityID)
	if err != nil {
		fmt.Println("Error converting activity ID to int:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	activity, err := db.GetActivityByID(c.MustGet("db").(*sql.DB), activityIDInt)
	if err != nil {
		fmt.Println("Error getting activity by ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	if activity.UserID != user.ID {
		fmt.Println("User is not authorized to delete this activity")
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this activity"})
		return
	}

	err = db.DeleteActivity(c.MustGet("db").(*sql.DB), activityIDInt)
	if err != nil {
		fmt.Println("Error deleting activity:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})

}
