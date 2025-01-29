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

	topic := cfg.KafkaUserRegisterTopic
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Error retrieving partitions: %v", err)
	}

	var wg sync.WaitGroup

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for _, partition := range partitions {
		wg.Add(1)

		go func(partition int32) {
			defer wg.Done()

			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)

			if err != nil {
				log.Fatalf("Error creating partition consumer: %v", err)
			}

			defer func() {
				if err := partitionConsumer.Close(); err != nil {
					log.Fatalf("Error closing partition consumer: %v", err)
				}
			}()

			for {
				select {
				case msg := <-partitionConsumer.Messages():
					log.Printf("Received message from partition %d at offset %d: %s\n", msg.Partition, msg.Offset, string(msg.Value))
					message := email.UserRegisterEvent{}
					err := json.Unmarshal(msg.Value, &message)

					if err != nil {
						log.Printf("Error unmarshaling the json", err)
					}

					email.SendWelcomeEmail(*cfg, message)
				case err := <-partitionConsumer.Errors():
					log.Printf("Error consuming message: %v", err)
				case <-signals:
					fmt.Println("Interrupt signal received. Closing consumer...")
					return
				}
			}

		}(partition)
	}

	wg.Wait()
	fmt.Println("Consumer closed.")
}
