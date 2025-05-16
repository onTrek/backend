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
		auth.POST("/register", api.PostRegister) // crea un nuovo utente
		auth.POST("/login", api.PostLogin)       // ritorna access token
	}

	// GPX API
	gpx := router.Group("/gpx")
	{
		gpx.DELETE("/:id", api.DeleteFile) // elimina un GPX
		gpx.POST("/", api.PostUpload)      // carica un file GPX
		gpx.GET("/", api.GetFiles)         // lista dei GPX caricati dall'utente
		gpx.GET("/:id", api.GetFile)       // scarica un GPX specifico
	}

	// Activity API
	activity := router.Group("/activity")
	{
		activity.POST("/", api.PostActivity)        // crea una nuova attività
		activity.PUT("/:id", api.PutActivity)       // aggiorna un'attività
		activity.GET("/", api.GetActivities)        // lista delle attività
		activity.GET("/:id", api.GetActivity)       // scarica un'attività specifica
		activity.DELETE("/:id", api.DeleteActivity) // elimina un'attività
	}

	// STATS API
	stats := router.Group("/stats")
	{
		stats.GET("", api.GetStats) // statistiche globali (tot km, tot salite, etc)
	}

	// SESSION API
	sessions := router.Group("/sessions")
	{
		sessions.POST("/", api.PostSession)      // crea una nuova sessione
		sessions.PUT("/:id", api.PutSession)     // aggiorna la posizione della sessione
		sessions.PATCH("/:id", api.PatchSession) // termina la sessione
		sessions.POST("/:id", api.PostSessionId) // partecipa a una sessione
		sessions.GET("/:id", api.GetSession)     // ritorna la sessione
		sessions.GET("/", api.GetSessions)       // lista delle sessioni
	}

	// FRIENDS API
	friends := router.Group("/friends")
	{
		friends.PUT("/:id", api.PutFriend)       // aggiungi un amico
		friends.GET("/", api.GetFriends)         // lista degli amici
		friends.DELETE("/:id", api.DeleteFriend) // elimina un amico
	}

	// USER API
	user := router.Group("/profile")
	{
		user.GET("", api.GetProfile)       // dati personali
		user.DELETE("", api.DeleteProfile) // cancella l'account
	}

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
