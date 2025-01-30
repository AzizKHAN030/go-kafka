package config

import (
	"log"
	"os"
)

type Config struct {
	AppPort                string
	KafkaBrokers           string
	KafkaUserRegisterTopic string
	KafkaUserLoginTopic    string
	KafkaClientID          string
	KafkaDQLTopic          string
	SmtpHost               string
	SmtpPort               string
	SmtpUsername           string
	SmtpPassword           string
	EmailFrom              string
	SmtpType               string
}

func LoadConfig() *Config {
	cfg := &Config{
		AppPort:                getEnv("APP_PORT", "8081"),
		KafkaBrokers:           getEnv("KAFKA_BROKERS", ""),
		KafkaUserRegisterTopic: getEnv("KAFKA_USER_REGISTER_TOPIC", "user_registered"),
		KafkaUserLoginTopic:    getEnv("KAFKA_USER_LOGIN_TOPIC", "user_login"),
		KafkaClientID:          getEnv("KAFKA_CLIENT_ID", "registration-service"),
		KafkaDQLTopic:          getEnv("KAFKA_DQL_TOPIC", "kafka_dql"),
		SmtpHost:               getEnv("SMTP_HOST", "SMTP"),
		SmtpPort:               getEnv("SMTP_PORT", ""),
		SmtpUsername:           getEnv("SMTP_USERNAME", ""),
		SmtpPassword:           getEnv("SMTP_PASSWORD", ""),
		SmtpType:               getEnv("SMTP_TYPE", "test"),
		EmailFrom:              getEnv("EMAIL_FROM", ""),
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
	if cfg.SmtpHost == "" {
		log.Fatal("SMTP_HOST is required!")
	}

	if cfg.SmtpPort == "" {
		log.Fatal("SMTP_PORT is required!")
	}

	if cfg.SmtpUsername == "" {
		log.Fatal("SMTP_USERNAME is required!")
	}
}
