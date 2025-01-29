package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/azizkhan030/go-kafka/notification/config"
	"github.com/azizkhan030/go-kafka/notification/email"
)

type ConsumerHandler struct {
	cfg config.Config
}

func StartConsumer(cfg *config.Config) {
	brokersURLSlice := strings.Split(cfg.KafkaBrokers, ",")
	consumer, err := sarama.NewConsumer(brokersURLSlice, nil)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing consumer: %v", err)
		}
	}()

	// Define the topics to listen to
	topics := []string{cfg.KafkaUserRegisterTopic, cfg.KafkaUserLoginTopic}

	var wg sync.WaitGroup
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for _, topic := range topics {
		partitions, err := consumer.Partitions(topic)
		if err != nil {
			log.Fatalf("Error retrieving partitions for topic %s: %v", topic, err)
		}

		for _, partition := range partitions {
			wg.Add(1)
			go func(topic string, partition int32) {
				defer wg.Done()

				partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
				if err != nil {
					log.Fatalf("Error creating partition consumer for topic %s: %v", topic, err)
				}

				defer func() {
					if err := partitionConsumer.Close(); err != nil {
						log.Fatalf("Error closing partition consumer for topic %s: %v", topic, err)
					}
				}()

				log.Printf("Listening on topic: %s, partition: %d\n", topic, partition)

				for {
					select {
					case msg := <-partitionConsumer.Messages():
						log.Printf("Received message from topic %s, partition %d, offset %d: %s\n",
							topic, msg.Partition, msg.Offset, string(msg.Value))

						// Process the message based on topic type
						var event email.UserRegisterEvent
						err := json.Unmarshal(msg.Value, &event)
						if err != nil {
							log.Printf("Error unmarshaling JSON for topic %s: %v", topic, err)
							continue
						}

						processEvent(cfg, event, topic)

					case err := <-partitionConsumer.Errors():
						log.Printf("Error consuming message from topic %s: %v", topic, err)
					case <-signals:
						fmt.Println("Interrupt signal received. Closing consumer...")
						return
					}
				}
			}(topic, partition)
		}
	}

	wg.Wait()
	fmt.Println("Consumer closed.")
}

// processEvent routes messages to the correct handler based on the event type
func processEvent(cfg *config.Config, event email.UserRegisterEvent, topic string) {
	switch topic {
	case cfg.KafkaUserRegisterTopic:
		log.Println("Processing user registration event...")
		email.SendWelcomeEmail(*cfg, event)

	case cfg.KafkaUserLoginTopic:
		log.Println("Processing user login event...")
		email.SendLoginNotification(*cfg, event)

	default:
		log.Printf("Unknown topic: %s, skipping...\n", topic)
	}
}
