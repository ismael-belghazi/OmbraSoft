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
		Port:    getEnv("PORT", "8080"),
		GINMode: getEnv("GIN_MODE", "debug"),

		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "postgres"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "ombrasoft"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),

		RedisURL: getEnv("REDIS_URL", ""),
		RedisDB:  0,

		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTExpiry: 24,

		AllowOrigins: parseOrigins(getEnv("ALLOW_ORIGINS", "")),
	}

	if cfg.DatabaseURL != "" {
		log.Println("Production detected - using DATABASE_URL")
		cfg.DBHost = ""
		cfg.DBPort = ""
		cfg.DBUser = ""
		cfg.DBPassword = ""
		cfg.DBName = ""
		cfg.DBSSLMode = ""
	} else {
		log.Println("Local mode detected - using local Docker DB")
	}

	validateConfig(cfg)
	logConfig(cfg)

	AppConfig = cfg
	return cfg
}

func validateConfig(c *Config) {
	if c.JWTSecret == "" {
		log.Fatal("JWT_SECRET requis! Configurez une variable d'environnement")
	}

	if c.GINMode == "release" {
		if c.DatabaseURL == "" {
			log.Fatal("DATABASE_URL obligatoire en production (GIN_MODE=release)")
		}

		if len(c.AllowOrigins) == 0 {
			log.Fatal("ALLOW_ORIGINS obligatoire en production! (ex: https://example.vercel.app)")
		}

		log.Println("Production mode - Configuration stricte validée")
	}
}

func logConfig(c *Config) {
	log.Println("Configuration chargée")
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
		log.Println("Mode DEBUG = Désactiver en production!")
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
