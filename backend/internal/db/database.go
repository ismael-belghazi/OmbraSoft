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
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Printf("Erreur connexion BD: %v", err)
		return err
	}

	log.Println("Connecté à PostgreSQL avec succès")

	if err := Migrate(); err != nil {
		log.Printf("Erreur migrations: %v", err)
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

	log.Println("✓ Migrations complétées")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
