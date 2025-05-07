package main

import (
	"OnTrek/api"
	"OnTrek/db"
	"github.com/gin-gonic/gin"
)

func main() {

	// Call the setup function from dbsetup package
	db.SetupDatabase()

	router := gin.Default()

	router.Use(db.DatabaseMiddleware())

	// AUTH API
	router.POST("/auth/register", api.PostRegister) // crea un nuovo utente
	router.POST("/auth/login", api.PostLogin)       // ritorna access token

	// GPX API
	router.POST("/gpx/upload", api.PostUpload) // carica un file GPX
	router.GET("/gpx/", api.GetFiles)          // lista dei GPX caricati dall'utente
	router.GET("/gpx/:id", api.GetFile)        // scarica un GPX specifico
	router.DELETE("/gpx/:id", api.DeleteFile)  // elimina un GPX

	// Activity API
	router.POST("/activity/", api.PostActivity)        // crea una nuova attività
	router.PUT("/activity/:id", api.PutActivity)       // aggiorna un'attività
	router.GET("/activity/", api.GetActivities)        // lista delle attività
	router.GET("/activity/:id", api.GetActivity)       // scarica un'attività specifica
	router.DELETE("/activity/:id", api.DeleteActivity) // elimina un'attività

	// STATS API
	router.GET("/stats", api.GetStats) // statistiche globali (tot km, tot salite, etc)

	// SESSION API
	router.POST("/session")     // crea una nuova sessione
	router.POST("/session/:id") // aggiorna la posizione della sessione
	router.GET("/session/:id")  // ritorna la sessione attiva

	// FRIENDS API
	router.PUT("/friends/:id", api.PutFriend)       // aggiungi un amico
	router.GET("/friends/", api.GetFriends)         // lista degli amici
	router.DELETE("/friends/:id", api.DeleteFriend) // elimina un amico

	// USER API
	router.GET("/user", api.GetProfile)       // dati personali
	router.DELETE("/user", api.DeleteProfile) // cancella l'account

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
