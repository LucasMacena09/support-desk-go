package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	RedisURL  string
	JWTSecret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("warn: .env not found, using system environment variables")
	}

	return &Config{
		Port:      getEnv("PORT", "8080"),
		RedisURL:  getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret: getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}