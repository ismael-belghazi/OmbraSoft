package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port    string
	GINMode string

	// Database - Production (DATABASE_URL) ou Local (HOST/PORT/USER/etc)
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string

	// Redis
	RedisURL string
	RedisDB  int

	// JWT
	JWTSecret string
	JWTExpiry int

	// CORS
	AllowOrigins []string
}

var AppConfig *Config

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		// Server
		Port:    getEnv("PORT", "8080"),
		GINMode: getEnv("GIN_MODE", "debug"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "postgres"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "ombrasoft"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),

		// Redis
		RedisURL: getEnv("REDIS_URL", ""),
		RedisDB:  0,

		// JWT
		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTExpiry: 24,

		// CORS
		AllowOrigins: parseOrigins(getEnv("ALLOW_ORIGINS", "http://localhost:5173,http://localhost:3000")),
	}

	// Mode production = SSL obligatoire si DATABASE_URL absent
	if cfg.GINMode == "release" {
		if cfg.DBSSLMode == "disable" && cfg.DatabaseURL == "" {
			cfg.DBSSLMode = "require"
		}
	}

	validateConfig(cfg)
	logConfig(cfg)

	AppConfig = cfg
	return cfg
}

func validateConfig(c *Config) {
	if c.JWTSecret == "" || c.JWTSecret == "dev_secret" {
		log.Fatal("JWT_SECRET non configuré! Utilisez une clé sécurisée en production")
	}
	if c.GINMode == "release" && len(c.AllowOrigins) == 0 {
		log.Fatal("ALLOW_ORIGINS vide en production!")
	}
}

func logConfig(c *Config) {
	log.Println("✓ Configuration chargée")
	log.Printf("  - Mode: %s", c.GINMode)
	log.Printf("  - Port: %s", c.Port)
	if c.DatabaseURL != "" {
		log.Printf("  - DB: DATABASE_URL configurée")
	} else {
		log.Printf("  - DB: %s@%s:%s (SSL: %s)", c.DBUser, c.DBHost, c.DBPort, c.DBSSLMode)
	}
	if c.RedisURL != "" {
		log.Printf("  - Redis: Activé")
	} else {
		log.Printf("  - Redis: Désactivé")
	}
	if c.GINMode == "debug" {
		log.Println("⚠  Mode DEBUG = Désactiver en production!")
	}
}

func parseOrigins(originsStr string) []string {
	if originsStr == "" {
		return []string{}
	}
	origins := strings.Split(originsStr, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetJWTSecret() string {
	if AppConfig != nil {
		return AppConfig.JWTSecret
	}
	return ""
}
