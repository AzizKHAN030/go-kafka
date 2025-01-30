package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/azizkhan030/go-kafka/registration/config"
	"github.com/azizkhan030/go-kafka/registration/db"
	"github.com/azizkhan030/go-kafka/registration/handlers"
	"github.com/azizkhan030/go-kafka/registration/kafka"
	"github.com/azizkhan030/go-kafka/registration/migrations"
	"github.com/gofiber/fiber/v2"
)

// Init Kafka producer
func main() {
	cfg := config.LoadConfig()

	db.ConnectDB(cfg.DatabaseURL)
	defer db.CloseDB()

	migrations.RunMigrations()

	kafka.ConnectProducer(cfg.KafkaBrokers)
	defer kafka.CloseProducer()

	topics := []string{}
	topics = append(topics, cfg.KafkaUserLoginTopic, cfg.KafkaUserRegisterTopic, cfg.KafkaDQLTopic)

	kafka.CreateTopics(cfg.KafkaBrokers, topics, 3, 1)

	app := fiber.New()

	handlers.SetupRoutes(app)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		kafka.CloseProducer()
		db.CloseDB()
		os.Exit(0)
	}()

	log.Printf("Starting server on port %s...", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
