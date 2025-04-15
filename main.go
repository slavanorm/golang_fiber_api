package main

import (
	"log"
	"rentincome/config"
	"rentincome/database"
	"rentincome/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	database.InitDB(config.Config("DatabaseDSN"))
	defer database.CloseDB()

	router.SetupRoutes(app)

	log.Fatal(app.Listen(config.Config("Host") + ":" + config.Config("Port")))
}
