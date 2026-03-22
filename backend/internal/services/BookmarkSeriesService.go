package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
	"gorm.io/gorm"

	"github.com/ismael-belghazi/ombrasoft-backend/internal/models"
)

type BookmarkSeriesService struct {
	DB *gorm.DB
}

func NewBookmarkSeriesService(db *gorm.DB) *BookmarkSeriesService {
	return &BookmarkSeriesService{DB: db}
}

type Episode struct {
	URL    string
	Number float64
	Title  string
}

// =======================
// Fetch Episodes
// =======================
func (s *BookmarkSeriesService) FetchEpisodes(seriesURL string) ([]Episode, error) {
	fmt.Println("▶ FetchEpisodes:", seriesURL)
	startTotal := time.Now()

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("Playwright start error: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{"--no-sandbox", "--disable-blink-features=AutomationControlled"},
	})
	if err != nil {
		return nil, fmt.Errorf("Browser launch error: %w", err)
	}
	defer browser.Close()

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/122 Safari/537.36"),
	})
	if err != nil {
		return nil, fmt.Errorf("Context creation error: %w", err)
	}

	page, err := context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("Page creation error: %w", err)
	}

	_, err = page.Goto(seriesURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(60000),
	})
	if err != nil {
		return nil, fmt.Errorf("Page navigation error: %w", err)
	}

	page.WaitForTimeout(2000)

	// Sélecteur par site
	var selector string
	switch {
	case strings.Contains(seriesURL, "mangadex.org"):
		selector = "a[href*='/chapter/']"
	case strings.Contains(seriesURL, "namicomi.com"):
		selector = "a.chapter-link"
	case strings.Contains(seriesURL, "raijin-scans.fr"):
		selector = "a.cairo"
	default:
		selector = "a[href]"
	}

	elements, err := page.QuerySelectorAll(selector)
	if err != nil {
		return nil, fmt.Errorf("QuerySelectorAll error: %w", err)
	}

	numRegexp := regexp.MustCompile(`\d+(\.\d+)?`)
	episodes := []Episode{}
	seen := make(map[string]bool)

	for _, el := range elements {
		href, _ := el.GetAttribute("href")
		if href == "" {
			continue
		}

		link := href
		if strings.HasPrefix(link, "/") {
			base := ""
			if strings.Contains(seriesURL, "mangadex.org") {
				base = "https://mangadex.org"
			} else if strings.Contains(seriesURL, "namicomi.com") {
				base = "https://namicomi.com"
			} else if strings.Contains(seriesURL, "raijin-scans.fr") {
				base = "https://raijin-scans.fr"
			}
			link = base + link
		}
		if seen[link] {
			continue
		}
		seen[link] = true

		title := ""
		if t, _ := el.GetAttribute("title"); t != "" {
			title = strings.TrimSpace(t)
		} else {
			txt, _ := el.TextContent()
			title = strings.TrimSpace(txt)
		}

		number := 0.0
		if matches := numRegexp.FindAllString(title, -1); len(matches) > 0 {
			number, _ = strconv.ParseFloat(matches[len(matches)-1], 64)
		}

		episodes = append(episodes, Episode{
			URL:    link,
			Number: number,
			Title:  title,
		})
	}

	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].Number < episodes[j].Number
	})

	fmt.Printf("Total episodes found: %d (%.2fs)\n", len(episodes), time.Since(startTotal).Seconds())
	return episodes, nil
}

// =======================
// Fetch Series Title
// =======================
func (s *BookmarkSeriesService) FetchSeries(seriesURL string) (*models.Series, error) {
	fmt.Println("▶ FetchSeries:", seriesURL)
	start := time.Now()

	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}

	_, err = page.Goto(seriesURL)
	if err != nil {
		return nil, err
	}

	page.WaitForTimeout(2000)

	title, err := page.Title()
	if err != nil || title == "" {
		title = "Unknown Series"
	}
	title = cleanTitle(title)

	fmt.Printf("Series title fetched: %s (%.2fs)\n", title, time.Since(start).Seconds())
	return &models.Series{
		ID:        uuid.New(),
		Title:     title,
		SourceURL: seriesURL,
	}, nil
}

// =======================
// Add Bookmark & Scraping Async
// =======================
func (s *BookmarkSeriesService) AddBookmarkWithChapters(userID string, seriesInfo *models.Series) error {
	if strings.TrimSpace(seriesInfo.Title) == "" {
		seriesInfo.Title = "Temporary Title"
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	var series models.Series
	err = s.DB.Where("source_url = ?", seriesInfo.SourceURL).First(&series).Error
	if err == gorm.ErrRecordNotFound {
		seriesInfo.ID = uuid.New()
		if err := s.DB.Create(seriesInfo).Error; err != nil {
			return err
		}
		series = *seriesInfo
	} else if err != nil {
		return err
	}

	// Création immédiate du bookmark
	bookmark := models.Bookmark{
		ID:       uuid.New(),
		UserID:   uid,
		SeriesID: series.ID,
	}
	if err := s.DB.Create(&bookmark).Error; err != nil {
		return err
	}

	// Lance une goroutine pour cover + update titre + chapitres
	go func(s *BookmarkSeriesService, seriesID uuid.UUID) {
		var sdb models.Series
		if err := s.DB.First(&sdb, "id = ?", seriesID).Error; err != nil {
			fmt.Println("Erreur rechargement série:", err)
			return
		}

		// Scraper la couverture mais ne bloque pas si échec
		if sdb.Cover == "" {
			if err := s.ScrapeAndSaveCover(&sdb); err != nil {
				fmt.Println("Erreur scraping cover:", err)
			}
		}

		// Mettre à jour le vrai titre si c'était temporaire
		if sdb.Title == "Temporary Title" {
			if updatedSeries, err := s.FetchSeries(sdb.SourceURL); err == nil && strings.TrimSpace(updatedSeries.Title) != "" {
				sdb.Title = updatedSeries.Title
				if err := s.DB.Save(&sdb).Error; err != nil {
					fmt.Println("Erreur update title:", err)
				}
			}
		}

		// Scraper les chapitres
		s.scrapeAndSaveChapters(seriesID, sdb.SourceURL)
	}(s, series.ID)

	return nil
}

// Scraper et sauvegarder la couverture
func (s *BookmarkSeriesService) ScrapeAndSaveCover(series *models.Series) error {
	// Si la cover existe déjà, rien à faire
	if series.Cover != "" {
		fmt.Println("Cover déjà présente:", series.Cover)
		return nil
	}

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("Playwright start error: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("Browser launch error: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("Page creation error: %w", err)
	}

	_, err = page.Goto(series.SourceURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(60000),
	})
	if err != nil {
		return fmt.Errorf("Page navigation error: %w", err)
	}

	// Attendre que l'image de couverture soit présente (max 15s)
	if _, err := page.WaitForSelector("img.cover", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(15000),
	}); err != nil {
		fmt.Println("Aucune image de cover trouvée pour", series.SourceURL)
		return nil // pas d'erreur, juste pas de cover
	}

	imgEl, err := page.QuerySelector("img.cover")
	if err != nil || imgEl == nil {
		fmt.Println("Cover selector non trouvé")
		return nil
	}

	src, _ := imgEl.GetAttribute("src")
	if src == "" {
		fmt.Println("Cover src vide")
		return nil
	}

	// Télécharger l'image
	resp, err := http.Get(src)
	if err != nil {
		return fmt.Errorf("Erreur téléchargement cover: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Erreur lecture cover: %w", err)
	}

	// Créer le dossier s'il n'existe pas
	dir := filepath.Join("public", "covers")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("Erreur création dossier covers: %w", err)
		}
	}

	// Sauvegarder l'image
	filename := series.ID.String() + ".jpg"
	filepath := filepath.Join(dir, filename)
	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("Erreur écriture cover: %w", err)
	}
	fmt.Println("Cover sauvegardée:", filepath)

	// Mettre à jour la série avec le chemin accessible par le front
	series.Cover = "/covers/" + filename
	if err := s.DB.Model(&models.Series{}).Where("id = ?", series.ID).Update("cover", series.Cover).Error; err != nil {
		return fmt.Errorf("Erreur update DB cover: %w", err)
	}

	return nil
}

// =======================
// Scrape & Save Chapters
// =======================
func (s *BookmarkSeriesService) scrapeAndSaveChapters(seriesID uuid.UUID, url string) {
	episodes, err := s.FetchEpisodes(url)
	if err != nil {
		fmt.Println("Erreur scraping:", err)
		return
	}

	for _, ep := range episodes {
		chapter := models.Chapter{
			ID:       uuid.New(),
			SeriesID: seriesID,
			URL:      ep.URL,
			Number:   ep.Number,
			Title:    ep.Title,
			Read:     false,
		}
		if err := s.DB.Create(&chapter).Error; err != nil {
			fmt.Println("Erreur création chapitre:", err)
		}
	}

	s.DB.Model(&models.Series{}).Where("id = ?", seriesID).Updates(map[string]interface{}{
		"last_chapter_number": len(episodes),
		"last_checked_at":     time.Now(),
	})
}

// =======================
// Mark Chapter/Series as Read
// =======================
func (s *BookmarkSeriesService) MarkChapterAsRead(userID, bookmarkID string, chapterNumber int) (int, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	bid, err := uuid.Parse(bookmarkID)
	if err != nil {
		return 0, err
	}

	var bookmark models.Bookmark
	if err := s.DB.Where("id = ? AND user_id = ?", bid, uid).First(&bookmark).Error; err != nil {
		return 0, err
	}

	if chapterNumber > bookmark.LastReadChapter {
		bookmark.LastReadChapter = chapterNumber
		if err := s.DB.Save(&bookmark).Error; err != nil {
			return 0, err
		}
	}

	return bookmark.LastReadChapter, nil
}

func (s *BookmarkSeriesService) MarkSeriesAsRead(userID, bookmarkID string) (int, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	bid, err := uuid.Parse(bookmarkID)
	if err != nil {
		return 0, err
	}

	var bookmark models.Bookmark
	if err := s.DB.Preload("Series").
		Where("id = ? AND user_id = ?", bid, uid).
		First(&bookmark).Error; err != nil {
		return 0, err
	}

	lastChapter := bookmark.Series.LastChapterNumber
	if lastChapter > bookmark.LastReadChapter {
		bookmark.LastReadChapter = lastChapter
		if err := s.DB.Save(&bookmark).Error; err != nil {
			return 0, err
		}
	}

	return bookmark.LastReadChapter, nil
}

// =======================
// Utils
// =======================
func (s *BookmarkSeriesService) GetBookmarksByUser(userID string) ([]models.Bookmark, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var bookmarks []models.Bookmark
	err = s.DB.Preload("Series.Chapters", func(db *gorm.DB) *gorm.DB {
		return db.Order("number ASC")
	}).Where("user_id = ?", uid).Find(&bookmarks).Error
	return bookmarks, err
}

func (s *BookmarkSeriesService) GetChaptersForSeries(seriesID string) ([]models.Chapter, error) {
	sid, err := uuid.Parse(seriesID)
	if err != nil {
		return nil, err
	}

	var chapters []models.Chapter
	err = s.DB.Where("series_id = ?", sid).Order("number ASC").Find(&chapters).Error
	return chapters, err
}

// =======================
// Update Chapters for Series
// =======================
func (s *BookmarkSeriesService) UpdateChaptersForSeries(seriesID string) error {
	sid, err := uuid.Parse(seriesID)
	if err != nil {
		return err
	}

	var series models.Series
	if err := s.DB.First(&series, "id = ?", sid).Error; err != nil {
		return err
	}

	// Met à jour cover si manquante
	if series.Cover == "" {
		go s.ScrapeAndSaveCover(&series)
	}

	episodes, err := s.FetchEpisodes(series.SourceURL)
	if err != nil {
		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("series_id = ?", sid).Delete(&models.Chapter{})
		for _, ep := range episodes {
			tx.Create(&models.Chapter{
				ID:       uuid.New(),
				SeriesID: sid,
				URL:      ep.URL,
				Number:   ep.Number,
				Title:    ep.Title,
				Read:     false,
			})
		}
		series.LastChapterNumber = len(episodes)
		series.LastCheckedAt = time.Now()
		return tx.Save(&series).Error
	})
}

func (s *BookmarkSeriesService) DeleteBookmarkAndSeriesForUser(bookmarkID, userID string) error {
	bid, err := uuid.Parse(bookmarkID)
	if err != nil {
		return err
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		var bookmark models.Bookmark
		if err := tx.Preload("Series").Where("id = ? AND user_id = ?", bid, uid).First(&bookmark).Error; err != nil {
			return err
		}

		seriesID := bookmark.SeriesID
		coverPath := bookmark.Series.Cover

		// Supprimer le bookmark
		if err := tx.Delete(&bookmark).Error; err != nil {
			return err
		}

		// Vérifier si d'autres bookmarks utilisent cette série
		var count int64
		if err := tx.Model(&models.Bookmark{}).Where("series_id = ?", seriesID).Count(&count).Error; err != nil {
			return err
		}

		if count == 0 {
			// Supprimer les chapitres
			if err := tx.Where("series_id = ?", seriesID).Delete(&models.Chapter{}).Error; err != nil {
				return err
			}

			// Supprimer la cover si ce n’est pas le default
			deleteCover(coverPath)

			// Supprimer la série
			if err := tx.Delete(&models.Series{}, "id = ?", seriesID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Fonction auxiliaire pour supprimer la cover
func deleteCover(coverPath string) {
	defaultCover := "/covers/default-cover.jpg"
	if coverPath == "" || coverPath == defaultCover {
		return
	}

	// Nettoyer le chemin et construire le chemin absolu
	cleaned := filepath.Clean(coverPath)
	absPath := filepath.Join("public", cleaned)

	// Supprimer le fichier
	if err := os.Remove(absPath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Cover %s déjà supprimée", absPath)
		} else {
			log.Printf("Erreur suppression cover %s: %v", absPath, err)
		}
	} else {
		log.Printf("Cover %s supprimée avec succès", absPath)
	}
}

// =======================
// Clean Title
// =======================
func cleanTitle(title string) string {
	title = strings.TrimSpace(title)
	re := regexp.MustCompile(`(?i)\s*[-|–|—]?\s*( VF |Scan VF|Scan FR|Raijin Scan).*`)
	title = re.ReplaceAllString(title, "")
	return strings.TrimSpace(title)
}
