package models

import "time"

type User struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	SecretHash   string    `json:"-"`
}

type Series struct {
	ID                string    `json:"id" gorm:"primaryKey"`
	Title             string    `json:"title" gorm:"not null"`
	SourceSite        string    `json:"source_site"`
	SourceURL         string    `json:"source_url" gorm:"not null"`
	LastChapterNumber int       `json:"last_chapter_number"`
	LastCheckedAt     time.Time `json:"last_checked_at"`
}

type Bookmark struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	UserID          string    `json:"user_id" gorm:"not null"`
	SeriesID        string    `json:"series_id" gorm:"not null"`
	LastReadChapter int       `json:"last_read_chapter"`
	UpdatedAt       time.Time `json:"updated_at"`
	User            User      `json:"user" gorm:"foreignKey:UserID"`
	Series          Series    `json:"series" gorm:"foreignKey:SeriesID"`
}

type UserNotifications struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"uniqueIndex"`
	Email     bool
	Push      bool
	DiscordID string `gorm:"column:discord_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Chapter struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	SeriesID  string    `json:"series_id" gorm:"not null"`
	URL       string    `json:"url" gorm:"not null"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
