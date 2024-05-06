package main

import (
	"log"
	"meigens-api/src/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := app.SetRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
