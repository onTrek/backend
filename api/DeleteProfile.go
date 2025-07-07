package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// DeleteProfile godoc
// @Summary Delete user profile
// @Description Deletes the user profile based on the provided authorization token.
// @Tags profile
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Success 204 "No Content"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to delete user"
// @Router /profile [delete]
func DeleteProfile(c *gin.Context) {

	user := c.MustGet("user").(utils.UserInfo)

	err := models.DeleteUser(c.MustGet("db").(*gorm.DB), user.ID)
	if err != nil {
		fmt.Println("Error deleting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
