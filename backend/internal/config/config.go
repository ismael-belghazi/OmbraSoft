package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret  string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	GINMode    string
	Port       string
}

var AppConfig *Config

func Load() *Config {
	_ = godotenv.Load()

	AppConfig = &Config{
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "ombrasoft"),
		GINMode:    getEnv("GIN_MODE", "debug"),
		Port:       getEnv("PORT", "8080"),
	}

	log.Printf("✓ Configuration chargée depuis variables d'environnement")
	if AppConfig.GINMode == "debug" {
		log.Println("⚠ Mode DEBUG activé - À désactiver en production")
	}
	return AppConfig
}

func GetJWTSecret() string {
	if AppConfig != nil {
		return AppConfig.JWTSecret
	}
	return "your-secret-key-change-in-production"
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
