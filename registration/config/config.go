package config

import (
	"log"
	"os"
)

// Load Kafka broker and topic settings
// Provide configs for environment-based variables (e.g., dotenv)

type Config struct {
	AppPort       string
	DatabaseURL   string
	KafkaBrokers  string
	KafkaTopic    string
	KafkaClientID string
	JWTSecret     string
}

func LoadConfig() *Config {
	cfg := &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		KafkaBrokers:  getEnv("KAFKA_BROKERS", ""),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "user-registration"),
		KafkaClientID: getEnv("KAFKA_CLIENT_ID", "registration-service"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
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
