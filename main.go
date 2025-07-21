// @title Ontrek API
// @version 1.0
// @description API documentation for Ontrek backend
// @host localhost:3000
// @BasePath /
package main

import (
	"OnTrek/api"
	"OnTrek/db"
	"OnTrek/db/functions"
	"OnTrek/db/models"
	_ "OnTrek/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db.SetupDatabase()
	url := ginSwagger.URL("/swagger/doc.json")

	err := models.CleanUnusedFiles(db.DB)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(db.DatabaseMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// AUTH API
	auth := router.Group("/auth")
	{
		auth.POST("/register", api.PostRegister)
		auth.POST("/login", api.PostLogin)
	}

	// GPX API
	gpx := router.Group("/gpx")
	gpx.Use(functions.AuthMiddleware())
	{
		gpx.DELETE("/:id", api.DeleteFile)
		gpx.POST("/", api.PostUpload)
		gpx.GET("/", api.GetFiles)
		gpx.GET("/:id", api.GetFileInfo)
		gpx.GET("/:id/download", api.GetFile)
		gpx.GET("/:id/map", api.GetFileMap)
	}

	// SESSION API
	groups := router.Group("/groups")
	groups.Use(functions.AuthMiddleware())
	{
		groups.POST("/", api.PostGroup)
		groups.GET("/", api.GetGroups)
		groups.DELETE("/:id", api.DeleteGroup)
		groups.GET("/:id", api.GetGroup)
		groups.PATCH("/:id/gpx", api.PatchSessionGpx)
		groups.PUT("/:id/members/location", api.PutGroupLocation)
		groups.PUT("/:id/members/", api.PutGroupId)
		groups.GET("/:id/members/", api.GetMembersInfo)
		groups.DELETE(":id/members/", api.DeleteLeaveRemoveMember)
	}

	search := router.Group("/search")
	search.Use(functions.AuthMiddleware())
	{
		search.GET("/", api.GetSearchPeople)
	}

	// FRIENDS API
	friends := router.Group("/friends")
	friends.Use(functions.AuthMiddleware())
	{
		friends.GET("/", api.GetFriends)
		friends.DELETE("/:id", api.DeleteFriend)
		friends.POST("/requests/:id", api.PostAddFriendRequest)
		friends.GET("/requests/received/", api.GetFriendRequestsReceived)
		friends.GET("/requests/sent/", api.GetFriendRequestsSent)
		friends.PUT("/requests/:id", api.PutAcceptFriendRequest)
		friends.DELETE("/requests/:id", api.DeleteDeclineFriendRequest)
	}

	// PROFILE API
	profile := router.Group("/profile")
	profile.Use(functions.AuthMiddleware())
	{
		profile.GET("", api.GetProfile)
		profile.PUT("/image", api.PutProfileImage)
		profile.DELETE("", api.DeleteProfile)
	}

	// USERS API
	users := router.Group("/users")
	users.Use(functions.AuthMiddleware())
	{
		users.GET("/:id/image", api.GetProfileImage)
	}

	router.Static("/docs", "./docs")

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
