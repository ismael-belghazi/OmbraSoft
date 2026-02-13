package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ismael-belghazi/ombrasoft-backend/internal/models"
)

type BookmarkSeriesService struct {
	DB *gorm.DB
}

func NewBookmarkSeriesService(db *gorm.DB) *BookmarkSeriesService {
	return &BookmarkSeriesService{DB: db}
}

func (s *BookmarkSeriesService) FetchSeries(query string) ([]models.Series, error) {
	searchURL := "https://example.com/search?q=" + query

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch series: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch series page: status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse series page: %w", err)
	}

	var seriesList []models.Series
	doc.Find(".series-item").Each(func(i int, sel *goquery.Selection) {
		title := sel.Find(".series-title").Text()
		link, _ := sel.Find("a").Attr("href")

		seriesList = append(seriesList, models.Series{
			ID:        uuid.NewString(),
			Title:     title,
			SourceURL: link,
		})
	})

	return seriesList, nil
}

func (s *BookmarkSeriesService) GetBookmarksByUser(userID string) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	err := s.DB.Preload("Series").Where("user_id = ?", userID).Find(&bookmarks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarks: %w", err)
	}
	return bookmarks, nil
}

func (s *BookmarkSeriesService) DeleteBookmarkForUser(bookmarkID, userID string) error {
	var bookmark models.Bookmark
	err := s.DB.Where("id = ? AND user_id = ?", bookmarkID, userID).First(&bookmark).Error
	if err != nil {
		return fmt.Errorf("bookmark not found or not owned by user: %w", err)
	}

	err = s.DB.Delete(&bookmark).Error
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}
	return nil
}

func (s *BookmarkSeriesService) AddBookmarkWithChapters(userID string, seriesInfo models.Series) error {
	var existingSeries models.Series
	err := s.DB.Where("source_url = ?", seriesInfo.SourceURL).First(&existingSeries).Error

	if err == gorm.ErrRecordNotFound {
		seriesInfo.ID = uuid.NewString()
		if err := s.DB.Create(&seriesInfo).Error; err != nil {
			return fmt.Errorf("failed to create series: %w", err)
		}

		chapters, err := s.FetchChapters(seriesInfo.SourceURL)
		if err != nil {
			return fmt.Errorf("failed to fetch chapters: %w", err)
		}

		for _, link := range chapters {
			chapter := models.Chapter{
				ID:       uuid.NewString(),
				SeriesID: seriesInfo.ID,
				URL:      link,
			}
			if err := s.DB.Create(&chapter).Error; err != nil {
				return fmt.Errorf("failed to save chapter: %w", err)
			}
		}
	}

	bookmark := models.Bookmark{
		ID:       uuid.NewString(),
		UserID:   userID,
		SeriesID: existingSeries.ID,
	}

	var existingBookmark models.Bookmark
	err = s.DB.Where("user_id = ? AND series_id = ?", userID, existingSeries.ID).First(&existingBookmark).Error
	if err == nil {
		return fmt.Errorf("bookmark already exists")
	}

	err = s.DB.Create(&bookmark).Error
	if err != nil {
		return fmt.Errorf("failed to create bookmark: %w", err)
	}
	return nil
}

func (s *BookmarkSeriesService) FetchChapters(seriesURL string) ([]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(seriesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch series page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch series page: status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse series page: %w", err)
	}

	var chapters []string
	doc.Find("li.wp-manga-chapter a").Each(func(i int, sel *goquery.Selection) {
		if link, exists := sel.Attr("href"); exists {
			chapters = append(chapters, link)
		}
	})

	return chapters, nil
}
