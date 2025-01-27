package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB(dbURL string) {
	var err error

	DB, err = pgxpool.New(context.Background(), dbURL)

	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}

	log.Println("Connected to the database successfully!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed.")
	}
}
