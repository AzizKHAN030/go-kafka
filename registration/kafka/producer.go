package kafka

import (
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/azizkhan030/go-kafka/registration/config"
)

// Initialize Sarama producer with broker settings
// Provide a function to send messages to Kafka
// Handle errors during message production

var producer sarama.SyncProducer

func ConnectProducer(brokersURL string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	brokersURLSlice := strings.Split(brokersURL, ",")
	var err error

	producer, err = sarama.NewSyncProducer(brokersURLSlice, config)

	if err != nil {
		log.Fatalf("%s", err)

		return nil, err
	}

	log.Println("Kafka producer connected successfully!")
	return producer, nil
}

func CloseProducer() {
	if producer != nil {
		err := producer.Close()
		if err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		} else {
			log.Println("Kafka producer closed successfully.")
		}
	}
}

func Publish(message []byte, topic string) error {
	cfg := config.LoadConfig()

	producer, err := ConnectProducer(cfg.KafkaBrokers)

	if err != nil {
		log.Fatalf("%s", err)

		return err
	}

	defer CloseProducer()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)

	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
