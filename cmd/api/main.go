package main

import (
	"log"
	"uptime-checker/internal/api"
	"uptime-checker/internal/config"
	"uptime-checker/internal/database"
)

func main() {
	db, err := database.InitDB(config.GetDSN())
	if err != nil {
		log.Fatal(err)
	}

	r := api.SetupRouter(db, config.GetJWTSecret())

	log.Println("Starting server on port 8080")
	r.Run(":8080")
}
