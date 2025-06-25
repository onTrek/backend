package api

import (
	"OnTrek/db/models"
	"OnTrek/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// GetFriendRequests godoc
// @Summary Get friend requests
// @Description Retrieve all friend requests for the authenticated user
// @Tags friends
// @Accept json
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Success 200 {array} utils.UserEssentials "List of friend requests"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to retrieve friend requests"
// @Router /friends/requests/ [get]
func GetFriendRequests(c *gin.Context) {

	user := c.MustGet("user").(utils.UserInfo)

	friendRequests, err := models.GetFriendRequestsByUserId(c.MustGet("db").(*gorm.DB), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve friend requests"})
		return
	}

	if len(friendRequests) == 0 {
		c.JSON(http.StatusOK, []utils.UserEssentials{})
		return
	}

	c.JSON(http.StatusOK, friendRequests)
}
