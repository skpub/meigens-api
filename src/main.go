package main

import (
	"log"
	"meigens-api/src/app"
	"meigens-api/src/dbconn"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := dbconn.Conn()
	if err != nil {
		panic("failed to connect db.")
	}

	router := app.SetRouter(db)
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}
