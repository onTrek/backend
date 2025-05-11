package main

import (
	"OnTrek/api"
	"OnTrek/db"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	// Call the setup function from dbsetup package
	db.SetupDatabase()
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition

	router := gin.Default()

	router.Use(db.DatabaseMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// AUTH API
	auth := router.Group("/auth")
	{
		auth.POST("/auth/register", api.PostRegister) // crea un nuovo utente
		auth.POST("/auth/login", api.PostLogin)       // ritorna access token
	}

	// GPX API
	gpx := router.Group("/gpx")
	{
		gpx.POST("/gpx/", api.PostUpload)      // carica un file GPX
		gpx.GET("/gpx/", api.GetFiles)         // lista dei GPX caricati dall'utente
		gpx.GET("/gpx/:id", api.GetFile)       // scarica un GPX specifico
		gpx.DELETE("/gpx/:id", api.DeleteFile) // elimina un GPX
	}

	// Activity API
	activity := router.Group("/activity")
	{
		activity.POST("/activity/", api.PostActivity)        // crea una nuova attività
		activity.PATCH("/activity/:id", api.PatchActivity)   // termina un'attività
		activity.GET("/activity/", api.GetActivities)        // lista delle attività
		activity.GET("/activity/:id", api.GetActivity)       // scarica un'attività specifica
		activity.DELETE("/activity/:id", api.DeleteActivity) // elimina un'attività
	}

	// STATS API
	stats := router.Group("/stats")
	{
		stats.GET("/stats", api.GetStats) // statistiche globali (tot km, tot salite, etc)
	}

	// SESSION API
	sessions := router.Group("/sessions")
	{
		sessions.POST("/sessions/", api.PostSession)      // crea una nuova sessione
		sessions.PUT("/sessions/:id", api.PutSession)     // aggiorna la posizione della sessione
		sessions.PATCH("/sessions/:id", api.PatchSession) // termina la sessione
		sessions.POST("/sessions/:id", api.PostSessionId) // partecipa a una sessione
		sessions.GET("/sessions/:id")                     // ritorna la sessione
		sessions.GET("/sessions/")                        // lista delle sessioni
	}

	// FRIENDS API
	friends := router.Group("/friends")
	{
		friends.PUT("/friends/:id", api.PutFriend)       // aggiungi un amico
		friends.GET("/friends/", api.GetFriends)         // lista degli amici
		friends.DELETE("/friends/:id", api.DeleteFriend) // elimina un amico
	}

	// USER API
	user := router.Group("/profile")
	{
		user.GET("/profile", api.GetProfile)       // dati personali
		user.DELETE("/profile", api.DeleteProfile) // cancella l'account
	}

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
