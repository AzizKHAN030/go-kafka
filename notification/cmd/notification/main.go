package main

import (
	"github.com/azizkhan030/go-kafka/notification/config"
	"github.com/azizkhan030/go-kafka/notification/kafka"
)

func main() {
	cfg := config.LoadConfig()

	kafka.StartConsumer(cfg)
}
