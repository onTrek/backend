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
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition

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
		auth.POST("/register", api.PostRegister) // crea un nuovo utente
		auth.POST("/login", api.PostLogin)       // ritorna access token
	}

	// GPX API
	gpx := router.Group("/gpx")
	gpx.Use(functions.AuthMiddleware())
	{
		gpx.DELETE("/:id", api.DeleteFile)    // elimina un GPX
		gpx.POST("/", api.PostUpload)         // carica un file GPX
		gpx.GET("/", api.GetFiles)            // lista dei GPX caricati dall'utente
		gpx.GET("/:id", api.GetFileInfo)      // ritorna i dati di un file GPX specifico
		gpx.GET("/:id/download", api.GetFile) // scarica un file GPX specifico
		gpx.GET("/:id/map", api.GetFileMap)   // scarica la mappa di un file GPX specifico
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
		search.GET("/", api.GetSearchPeople) // ricerca persone da aggiungere agli amici
	}

	// FRIENDS API
	friends := router.Group("/friends")
	friends.Use(functions.AuthMiddleware())
	{
		friends.GET("/", api.GetFriends)                                  // lista degli amici
		friends.DELETE("/:id", api.DeleteFriend)                          // elimina un amico
		friends.POST("/requests/:id", api.PostAddFriendRequest)           // aggiungi un amico
		friends.GET("/requests/received/", api.GetFriendRequestsReceived) // lista delle richieste di amicizia ricevute
		friends.GET("/requests/sent/", api.GetFriendRequestsSent)         // lista delle richieste di amicizia inviate
		friends.PUT("/requests/:id", api.PutAcceptFriendRequest)          // accetta una richiesta di amicizia
		friends.DELETE("/requests/:id", api.DeleteDeclineFriendRequest)   // rifiuta una richiesta di amicizia
	}

	// PROFILE API
	profile := router.Group("/profile")
	profile.Use(functions.AuthMiddleware())
	{
		profile.GET("", api.GetProfile)            // dati personali
		profile.PUT("/image", api.PutProfileImage) // carica la foto profilo
		profile.DELETE("", api.DeleteProfile)      // cancella l'account
	}

	// USERS API
	users := router.Group("/users")
	users.Use(functions.AuthMiddleware())
	{
		users.GET("/:id/image", api.GetProfileImage) // scarica la foto profilo
	}

	router.Static("/docs", "./docs")

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
