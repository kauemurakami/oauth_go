package main

import (
	"fmt"
	"log"
	"net/http"
	app_config "oauth2/core/config"
	"oauth2/data/db"
	"oauth2/routes"
)

func main() {
	app_config.SetupEnvironments()
	db.SetupDB()
	fmt.Printf("Run API :%s", app_config.API_port)
	router := routes.SetupAppRoutes()
	// router.Use(middlewares.SetupHeadersMiddleware)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", app_config.API_port), router))
}
