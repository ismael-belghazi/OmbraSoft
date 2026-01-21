package db

import (
	"fmt"
	"log"

	"github.com/bebeb/ombrasoft-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(host, port, user, password, dbname string) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erreur de connexion à PostgreSQL:", err)
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (d *Database) Migrate() error {
	return d.DB.AutoMigrate(
		&models.User{},
		&models.Series{},
		&models.Bookmark{},
	)
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
