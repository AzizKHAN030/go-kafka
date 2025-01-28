package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func CreateTopics(brokers string, topics []string, partitionsNum int32, replicationFactor int16) {
	admin, err := sarama.NewClusterAdmin([]string{brokers}, sarama.NewConfig())

	if err != nil {
		log.Fatalf("Failed to create kafka adminL %v", err)
	}
	defer admin.Close()

	for _, topic := range topics {
		err = admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     partitionsNum,
			ReplicationFactor: replicationFactor,
		}, false)

		if err != nil && err != sarama.ErrTopicAlreadyExists {
			log.Fatalf("Error creating topic %s: %v", topic, err)
		}

		log.Printf("Topic %s created or already exists.", topic)
	}
}
