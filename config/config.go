package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	DBMaxIdleConns    int
	DBMaxOpenConns    int
	DBConnMaxLifetime time.Duration
	RabbitURL         string
	ProductAPI        string
	ServicePort       string
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func Load() *Config {
	return &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPass:            getEnv("DB_PASS", "postgres"),
		DBName:            getEnv("DB_NAME", "product_db"),
		DBMaxIdleConns:    getEnvInt("POSTGRES_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns:    getEnvInt("POSTGRES_MAX_OPEN_CONNS", 100),
		DBConnMaxLifetime: time.Duration(getEnvInt("POSTGRES_CONN_MAX_LIFETIME_MINUTES", 30)) * time.Minute,
		RabbitURL:         getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ProductAPI:        getEnv("PRODUCT_API_URL", "http://localhost:3001/products"),
		ServicePort:       getEnv("PORT", "3002"),
	}
}
