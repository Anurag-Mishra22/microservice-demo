package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction: 'up' or 'down'")
	}

	direction := os.Args[1]

	// Build connection string from environment variables
	var connString string

	// Check if DATABASE_URL is provided first
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connString = dbURL
	} else {
		// Build connection string from individual environment variables
		host := getEnvOrDefault("DB_HOST", "localhost")
		port := getEnvOrDefault("DB_PORT", "5432")
		dbname := getEnvOrDefault("DB_NAME", "dbname")
		user := getEnvOrDefault("DB_USER", "username")
		password := getEnvOrDefault("DB_PASSWORD", "password")
		sslmode := getEnvOrDefault("DB_SSL_MODE", "disable")

		connString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbname, sslmode)
	}

	log.Printf("Connecting to database: %s", connString)

	// Parse the connection string
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		log.Fatal(err)
	}

	// Open database connection
	db := stdlib.OpenDB(*config)
	defer db.Close()

	// Test the connection
	if err := db.PingContext(context.Background()); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migrations rolled back successfully")
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'.")
	}
}

// getEnvOrDefault returns the environment variable value or a default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
