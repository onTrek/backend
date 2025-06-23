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
	_ "OnTrek/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	db.SetupDatabase()
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition

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
		gpx.GET("/:id/download", api.GetFile) // scarica un file GPX specifico
		gpx.GET("/:id/map", api.GetFileMap)   // scarica la mappa di un file GPX specifico
	}

	// SESSION API
	sessions := router.Group("/groups")
	sessions.Use(functions.AuthMiddleware())
	{
		sessions.POST("/", api.PostGroup)
		sessions.GET("/", api.GetGroups)
		sessions.DELETE("/:id", api.DeleteGroup)
		sessions.GET("/:id", api.GetGroup)
		sessions.PATCH("/:id/gpx", api.PatchSessionGpx)
		sessions.PUT("/:id/members/location", api.PutGroupLocation)
		sessions.PUT("/:id/members/", api.PutGroupId)
		sessions.GET("/:id/members/", api.GetMembersInfo)
		sessions.DELETE(":id/members/", api.DeleteLeaveRemoveMember)
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
		friends.GET("/", api.GetFriends)                                // lista degli amici
		friends.DELETE("/:id", api.DeleteFriend)                        // elimina un amico
		friends.POST("/requests/:id", api.PostAddFriendRequest)         // aggiungi un amico
		friends.GET("/requests/", api.GetFriendRequests)                // lista delle richieste di amicizia
		friends.PUT("/requests/:id", api.PutAcceptFriendRequest)        // accetta una richiesta di amicizia
		friends.DELETE("/requests/:id", api.DeleteDeclineFriendRequest) // rifiuta una richiesta di amicizia
	}

	// USER API
	user := router.Group("/profile")
	user.Use(functions.AuthMiddleware())
	{
		user.GET("", api.GetProfile)       // dati personali
		user.DELETE("", api.DeleteProfile) // cancella l'account
	}

	router.Static("/docs", "./docs")

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
