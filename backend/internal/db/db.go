package db

import (
	"fmt"
	"time"

	"github.com/ismael-belghazi/ombrasoft-backend/internal/config"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init initialise la connexion à la base et exécute les migrations
func Init() error {
	cfg := config.AppConfig

	// Construction de la DSN
	var dsn string
	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBSSLMode,
		)
	}

	// Niveau de log GORM
	logLevel := logger.Silent
	if cfg.GINMode == "debug" {
		logLevel = logger.Info
	}

	// Connexion à la base
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configuration du pool de connexions
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Exécution des migrations
	if err := Migrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// Migrate crée les tables et indexes si nécessaire
func Migrate() error {
	migrationSQL := `
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    source_site VARCHAR(255),
    source_url VARCHAR(500) NOT NULL UNIQUE,
    last_chapter_number INT DEFAULT -1,
    last_checked_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    last_read_chapter INT DEFAULT -1,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, series_id)
);

CREATE TABLE IF NOT EXISTS chapters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    url VARCHAR(1000) NOT NULL,
    number INT NOT NULL,
    title VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(series_id, url)
);

CREATE TABLE IF NOT EXISTS user_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    email BOOLEAN DEFAULT FALSE,
    push BOOLEAN DEFAULT FALSE,
    discord_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_series_id ON bookmarks(series_id);
CREATE INDEX IF NOT EXISTS idx_chapters_series_id ON chapters(series_id);
CREATE INDEX IF NOT EXISTS idx_series_source_url ON series(source_url);
`

	if err := DB.Exec(migrationSQL).Error; err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// AutoMigrate pour les modèles GORM
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Series{},
		&models.Bookmark{},
		&models.UserNotifications{},
		&models.Chapter{},
	); err != nil {
		return fmt.Errorf("failed to auto migrate models: %w", err)
	}

	return nil
}

// GetDB retourne l'instance globale de la DB
func GetDB() *gorm.DB {
	return DB
}

// Close ferme la connexion à la base
func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql DB: %w", err)
	}
	return sqlDB.Close()
}
