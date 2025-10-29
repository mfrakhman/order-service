package config

import "os"

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPass       string
	DBName       string
	RabbitURL    string
	ProductAPI   string
	ServicePort  string
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() *Config {
	return &Config{
		DBHost:      getEnv("POSTGRES_HOST", "localhost"),
		DBPort:      getEnv("POSTGRES_PORT", "5432"),
		DBUser:      getEnv("POSTGRES_USER", "postgres"),
		DBPass:      getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:      getEnv("POSTGRES_DB", "product_db"),
		RabbitURL:   getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ProductAPI:  getEnv("PRODUCT_API_URL", "http://localhost:3001/products"),
		ServicePort: getEnv("PORT", "3002"),
	}
}
