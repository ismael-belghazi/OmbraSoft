package db

import (
	"fmt"
	"log"

	"github.com/bebeb/ombrasoft-backend/internal/config"
	"github.com/bebeb/ombrasoft-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	cfg := config.AppConfig

	var dsn string

	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBSSLMode,
		)
	}

	logLevel := logger.Silent
	if cfg.GINMode == "debug" {
		logLevel = logger.Info
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Printf("❌ Erreur connexion BD: %v", err)
		return err
	}

	log.Println("✓ Connecté à PostgreSQL")

	if err := Migrate(); err != nil {
		log.Printf("❌ Erreur migrations: %v", err)
		return err
	}

	return nil
}

func Migrate() error {
	log.Println("🔄 Exécution des migrations...")
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Series{},
		&models.Bookmark{},
	); err != nil {
		return err
	}
	log.Println("✓ Migrations OK")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
