package config

import (
	"log"
	"os"
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
	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "pr_service"),
		DBUser:     getEnv("DB_USER", "pr_user"),
		DBPassword: getEnv("DB_PASSWORD", "pr_password"),
		AppPort:    getEnv("APP_PORT", "8080"),
	}

	log.Printf("Config loaded: %+v\n", cfg)

	return cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
