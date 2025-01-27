package migrations

import (
	"context"
	"log"

	"github.com/azizkhan030/go-kafka/registration/db"
)

// RunMigrations runs database migrations
func RunMigrations() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.DB.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v\n", err)
	}

	log.Println("Database migrations applied successfully!")
}
