package config

import (
	"log"
	"os"
)

// Load Kafka broker and topic settings
// Provide configs for environment-based variables (e.g., dotenv)

type Config struct {
	AppPort                string
	DatabaseURL            string
	KafkaBrokers           string
	KafkaUserRegisterTopic string
	KafkaUserLoginTopic    string
	KafkaClientID          string
	KafkaDQLTopic          string
	JWTSecret              string
}

func LoadConfig() *Config {
	cfg := &Config{
		AppPort:                getEnv("APP_PORT", "8080"),
		DatabaseURL:            getEnv("DATABASE_URL", ""),
		KafkaBrokers:           getEnv("KAFKA_BROKERS", ""),
		KafkaUserRegisterTopic: getEnv("KAFKA_USER_REGISTER_TOPIC", "user_registered"),
		KafkaUserLoginTopic:    getEnv("KAFKA_USER_LOGIN_TOPIC", "user_login"),
		KafkaClientID:          getEnv("KAFKA_CLIENT_ID", "registration-service"),
		KafkaDQLTopic:          getEnv("KAFKA_DQL_TOPIC", "kafka_dql"),
		JWTSecret:              getEnv("JWT_SECRET", ""),
	}

	validateConfig(cfg)

	return cfg
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func validateConfig(cfg *Config) {
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required!")
	}

	if cfg.KafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS is required!")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required!")
	}
}
