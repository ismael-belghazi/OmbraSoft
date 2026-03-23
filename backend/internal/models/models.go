package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-"`
	SecretHash   string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Series struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Title             string    `json:"title"`
	SourceSite        string    `json:"source_site"`
	SourceURL         string    `json:"source_url" gorm:"not null"`
	Cover             string    `json:"cover_image_url"`
	LastChapterNumber int       `json:"last_chapter_number"`
	LastCheckedAt     time.Time `json:"last_checked_at"`

	Chapters []Chapter `json:"chapters" gorm:"foreignKey:SeriesID"`
}

type Bookmark struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	SeriesID        uuid.UUID `json:"series_id" gorm:"type:uuid;not null"`
	LastReadChapter int       `json:"last_read_chapter"`
	UpdatedAt       time.Time `json:"updated_at"`

	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Series Series `json:"series" gorm:"foreignKey:SeriesID"`
}

type Chapter struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	SeriesID  uuid.UUID `json:"series_id" gorm:"type:uuid;not null"`
	URL       string    `json:"url" gorm:"not null"`
	Number    float64   `json:"number"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Read      bool      `json:"read"`
}

type UserNotifications struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Email     bool
	Push      bool
	DiscordID string `gorm:"column:discord_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
