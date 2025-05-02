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

	// STATS API
	router.GET("/stats/global")  // statistiche globali (tot km, tot salite, etc)
	router.GET("/stats/gpx/")    // statistiche di tutte le tracce
	router.GET("/stats/gpx/:id") // statistiche di una singola traccia

	// USER API
	router.GET("/user/profile")    // dati personali
	router.PUT("/user/profile")    // aggiorna i dati personali
	router.DELETE("/user/account") // cancella l'account

	// ACTIVITY API
	router.GET("/activity/find") // manda la posizione e riceve l'id della traccia GPX più vicina

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
