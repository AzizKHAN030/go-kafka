package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/azizkhan030/go-kafka/notification/config"
	"github.com/azizkhan030/go-kafka/notification/email"
)

type ConsumerHandler struct {
	cfg *config.Config
}

var MaxRetries = 5

func StartConsumer(cfg *config.Config) {
	brokersURLSlice := strings.Split(cfg.KafkaBrokers, ",")
	groupId := "notification-service-group"

	consumerGroup, err := sarama.NewConsumerGroup(brokersURLSlice, groupId, nil)
	if err != nil {
		log.Fatalf("Error starting kafka consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Define the topics to listen to
	topics := []string{cfg.KafkaUserRegisterTopic, cfg.KafkaUserLoginTopic}

	handler := ConsumerHandler{cfg: cfg}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for {
			err := consumerGroup.Consume(context.Background(), topics, &handler)

			if err != nil {
				log.Printf("Error consuming messages: %v", err)
			}
		}
	}()

	<-signals
	log.Println("Shutting down consumer...")
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event email.UserRegisterEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		success := false

		for attempt := 1; attempt <= MaxRetries; attempt++ {
			log.Printf("Processing message from topic %s (Attempt %d/%d)", msg.Topic, attempt, MaxRetries)

			err := processEvent(h.cfg, event, msg.Topic)

			if err == nil {
				success = true
				break
			}

			log.Printf("Error processing message (Attempt %d/%d): %v", attempt, MaxRetries, err)
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}

		if !success {
			log.Printf("Message failed after %d retries. Sending to DQL: %s", MaxRetries, msg.Value)
			sendToDQL(h.cfg, msg)
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func processEvent(cfg *config.Config, event email.UserRegisterEvent, topic string) error {
	switch topic {
	case cfg.KafkaUserRegisterTopic:
		return email.SendWelcomeEmail(*cfg, event)
	case cfg.KafkaUserLoginTopic:
		return email.SendLoginNotification(*cfg, event)
	default:
		log.Printf("Unknown topic: %s, skipping ...\n", topic)
		return nil
	}
}

func sendToDQL(cfg *config.Config, msg *sarama.ConsumerMessage) {
	producer, err := sarama.NewSyncProducer(strings.Split(cfg.KafkaBrokers, ","), nil)

	if err != nil {
		log.Printf("Failed to create DQL producer: %v", err)
		return
	}

	defer producer.Close()

	dqlTopic := cfg.KafkaDQLTopic

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: dqlTopic,
		Value: sarama.StringEncoder(msg.Value),
	})

	if err != nil {
		log.Printf("Failed to send message to DQL: %v", err)
	} else {
		log.Printf("Message sent to DQL: %s", string(msg.Value))
	}
}
