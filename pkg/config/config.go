package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	MQUsername string
	MQPassword string
	MQPort     string
	MQHost     string

	SMTPHost string
	SMTPPort string
}

var Envs = initConfig()

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("POST", "8080"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "mypassword"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBName:     getEnv("DB_NAME", "library"),
		DBPort:     getEnv("DB_PORT", "5432"),
		MQUsername: getEnv("MQ_USERNAME", "guest"),
		MQPassword: getEnv("MQ_PASSWORD", "guest"),
		MQHost:     getEnv("MQ_HOST", "localhost"),
		MQPort:     getEnv("MQ_PORT", "5672"),
		SMTPHost:   getEnv("SMTP_HOST", "mailhog"),
		SMTPPort:   getEnv("SMTP_PORT", "1025"),
	}
}
