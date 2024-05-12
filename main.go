package main

import (
	"hermes/database"
	"hermes/servers"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	database.InitDatabase()
	err := godotenv.Load()
	if err != nil {
		log.Printf("[ENV-VAR] - %v - Error loading .env file\n", time.Now())
	}
}


func main() {
	servers.InitHTTPServer()
}