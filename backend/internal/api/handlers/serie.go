package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/models"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/services"
	"gorm.io/gorm"
)

type SeriesHandler struct {
	BookmarkService *services.BookmarkSeriesService
	ScraperService  *services.BookmarkSeriesService
}

func NewSeriesHandler(bookmarkService *services.BookmarkSeriesService) *SeriesHandler {
	return &SeriesHandler{
		BookmarkService: bookmarkService,
		ScraperService:  bookmarkService,
	}
}

////////////////////////////////////////////////////////////
// UTILS
////////////////////////////////////////////////////////////

func isValidHTTPURL(raw string) (*url.URL, bool) {
	u, err := url.ParseRequestURI(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, false
	}
	return u, true
}

////////////////////////////////////////////////////////////
// CREATE BOOKMARK
////////////////////////////////////////////////////////////

func (h *SeriesHandler) CreateBookmark(c *gin.Context) {
	var req struct {
		SourceURL string `json:"sourceURL" binding:"required"`
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "data": nil, "error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "error": "Invalid body: " + err.Error()})
		return
	}

	req.SourceURL = strings.TrimSpace(req.SourceURL)

	if _, ok := isValidHTTPURL(req.SourceURL); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "error": "Invalid URL"})
		return
	}

	seriesInfo, err := h.BookmarkService.FetchSeries(req.SourceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "error": "Failed to fetch series: " + err.Error()})
		return
	}

	if strings.TrimSpace(seriesInfo.Title) == "" {
		seriesInfo.Title = "Unknown Series"
	}

	if err := h.BookmarkService.AddBookmarkWithChapters(userID, seriesInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "error": "Failed to create bookmark: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": gin.H{"series": seriesInfo}, "error": ""})
}

////////////////////////////////////////////////////////////
// REDIRECT CHAPTER
////////////////////////////////////////////////////////////

func (h *SeriesHandler) RedirectToChapter(c *gin.Context) {
	chapterID := strings.TrimSpace(c.Param("id"))
	if chapterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing chapter ID"})
		return
	}

	var chapter models.Chapter
	err := h.BookmarkService.DB.First(&chapter, "id = ?", chapterID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	u, ok := isValidHTTPURL(chapter.URL)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter URL"})
		return
	}

	c.Redirect(http.StatusFound, u.String())
}

////////////////////////////////////////////////////////////
// REDIRECT SERIES
////////////////////////////////////////////////////////////

func (h *SeriesHandler) RedirectToSeries(c *gin.Context) {
	raw := strings.TrimSpace(c.Query("url"))
	if raw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing url"})
		return
	}

	u, ok := isValidHTTPURL(raw)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	c.Redirect(http.StatusFound, u.String())
}

////////////////////////////////////////////////////////////
// SEARCH SERIES
////////////////////////////////////////////////////////////

func (h *SeriesHandler) SearchSeries(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	urlParam := strings.TrimSpace(c.Query("url"))

	if query != "" {
		c.JSON(http.StatusOK, gin.H{"success": true, "results": []models.Series{}})
		return
	}

	if urlParam != "" {
		if _, ok := isValidHTTPURL(urlParam); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
			return
		}

		series, err := h.ScraperService.FetchSeries(urlParam)
		if err != nil || series == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch series", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "series": series})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'q' or 'url' parameter"})
}
