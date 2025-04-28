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
	router.POST("/auth/login")                      // ritorna access token

	// GPX API
	router.POST("/gpx/upload") // carica un file GPX
	router.GET("/gpx/:id")     // scarica un GPX specifico
	router.GET("/gpx/list")    // lista dei GPX caricati dall'utente
	router.DELETE("/gpx/:id")  // elimina un GPX

	// STATS API
	router.GET("/stats/global")  // statistiche globali (tot km, tot salite, etc)
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
