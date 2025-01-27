package main

import (
	"log"

	"github.com/azizkhan030/go-kafka/registration/config"
	"github.com/azizkhan030/go-kafka/registration/db"
	"github.com/azizkhan030/go-kafka/registration/handlers"
	"github.com/azizkhan030/go-kafka/registration/migrations"
	"github.com/gofiber/fiber/v2"
)

// Init Kafka producer
// TODO Login handler is not working
func main() {
	cfg := config.LoadConfig()

	db.ConnectDB(cfg.DatabaseURL)
	defer db.CloseDB()

	migrations.RunMigrations()

	app := fiber.New()

	handlers.SetupRoutes(app)

	log.Printf("Starting server on port %s...", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
