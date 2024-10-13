package main

import (
	"log"
	"ssf/db"
	"ssf/routes"
	views "ssf/views/helpers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	InitApp()
	routes.InitRoutes(App)
	views.PrepareTemplates()
	db.ConnectDatabase()
	StartApp()
}
