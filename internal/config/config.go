package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Server   serverConfig
	Database databaseConfig
}

type serverConfig struct {
	Port string
}

type databaseConfig struct {
	Url string
}

func Load() *config {
	godotenv.Load()

	return &config{
		Server: serverConfig{
			Port: getEnvOptional("PORT", "8080"),
		},
		Database: databaseConfig{
			Url: getEnv("DB_URL"),
		},
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s must be set", key)
	}

	return value
}

func getEnvOptional(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
