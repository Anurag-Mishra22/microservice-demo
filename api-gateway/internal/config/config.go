package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnvOrDefault("SERVER_PORT", "8080"),
			Host:         getEnvOrDefault("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnvOrDefault("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "username"),
			Password: getEnvOrDefault("DB_PASSWORD", "password"),
			DBName:   getEnvOrDefault("DB_NAME", "dbname"),
			SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
			MaxConns: getIntEnvOrDefault("DB_MAX_CONNS", 25),
		},
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.DBName, c.Database.SSLMode)
}
