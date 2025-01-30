package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/IBM/sarama"
	"github.com/azizkhan030/go-kafka/notification/config"
	"github.com/azizkhan030/go-kafka/notification/email"
)

type ConsumerHandler struct {
	cfg *config.Config
}

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

		log.Printf("Received message from topic %s: %+v", msg.Topic, event)

		switch msg.Topic {
		case h.cfg.KafkaUserRegisterTopic:
			email.SendWelcomeEmail(*h.cfg, event)
		case h.cfg.KafkaUserLoginTopic:
			email.SendLoginNotification(*h.cfg, event)
		}

		session.MarkMessage(msg, "")
	}

	return nil
}
