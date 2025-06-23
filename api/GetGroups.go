package api

import (
	"OnTrek/db/functions"
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetGroups godoc
// @Summary Get groups
// @Description Get all groups for the authenticated user
// @Tags groups
// @Produce json
// @Param Bearer header string true "Bearer token for user authentication"
// @Success 200 {object} []utils.GroupDoc "List of groups"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Error fetching files"
// @Router /groups/ [get]
func GetGroups(c *gin.Context) {

	// Get the user from the context
	user := c.MustGet("user").(utils.User)

	groups, err := functions.GetGroupsByUserId(c.MustGet("db").(*sql.DB), user.ID)
	if err != nil {
		fmt.Println("Error getting sessions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if len(groups) == 0 {
		c.JSON(http.StatusOK, []utils.GroupDoc{})
		return
	}

	c.JSON(http.StatusOK, groups)
}
