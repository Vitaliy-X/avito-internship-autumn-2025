package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	AppPort    string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBName:     getEnv("DB_NAME"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		AppPort:    getEnv("APP_PORT"),
	}

	log.Printf("Config loaded: %+v\n", cfg)
	return cfg
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	log.Fatalf("Required environment variable %s not set", key)
	return ""
}
