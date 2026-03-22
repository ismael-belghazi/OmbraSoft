package routes

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/api/handlers"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/middleware"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/services"
)

// SeriesRoutes initialise les routes pour la gestion des séries
func SeriesRoutes(router *gin.RouterGroup, bookmarkService *services.BookmarkSeriesService) {
	seriesHandler := handlers.NewSeriesHandler(bookmarkService)

	series := router.Group("/series")
	series.Use(middleware.AuthMiddleware())
	{
		// Créer un bookmark avec scraping automatique des chapitres
		series.POST("", func(c *gin.Context) {
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
			ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				// Fetch the series info from URL
				seriesModel, err := seriesHandler.BookmarkService.FetchSeries(req.SourceURL)
				if err != nil {
					done <- err
					return
				}
				done <- seriesHandler.BookmarkService.AddBookmarkWithChapters(userID, seriesModel)
			}()

			select {
			case err := <-done:
				if err != nil {
					log.Printf("[CreateBookmark] error: %v", err)
					status := http.StatusInternalServerError
					msg := err.Error()

					if containsIgnoreCase(msg, "cloudflare") {
						status = http.StatusBadGateway
						msg = "Blocked by target site"
					} else if containsIgnoreCase(msg, "no episodes") {
						status = http.StatusNotFound
						msg = "No chapters found"
					}

					c.JSON(status, gin.H{"success": false, "data": nil, "error": msg})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"success": true, "data": nil, "error": ""})
			case <-ctx.Done():
				c.JSON(http.StatusGatewayTimeout, gin.H{"success": false, "data": nil, "error": "Scraping timeout"})
			}
		})

		// Obtenir tous les chapitres d'une série
		series.GET("/:id/chapters", func(c *gin.Context) {
			seriesID, ok := getSeriesID(c)
			if !ok {
				return
			}

			chapters, err := bookmarkService.GetChaptersForSeries(seriesID)
			if err != nil {
				log.Printf("[GetChapters] error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "error": "Failed to get chapters"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": chapters, "error": ""})
		})

		// Mettre à jour les chapitres d'une série via scraping
		series.POST("/:id/update-chapters", func(c *gin.Context) {
			seriesID, ok := getSeriesID(c)
			if !ok {
				return
			}

			ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				done <- bookmarkService.UpdateChaptersForSeries(seriesID)
			}()

			select {
			case err := <-done:
				if err != nil {
					log.Printf("[UpdateChapters] error: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "error": "Failed to update chapters: " + err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"success": true, "data": nil, "error": ""})
			case <-ctx.Done():
				c.JSON(http.StatusGatewayTimeout, gin.H{"success": false, "data": nil, "error": "Scraping timeout"})
			}
		})

		// Tester le scraping sans créer de bookmark
		series.POST("/test-scrape", func(c *gin.Context) {
			var req struct {
				URL string `json:"url" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "error": "URL is required"})
				return
			}

			ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
			defer cancel()

			done := make(chan []services.Episode, 1)
			errChan := make(chan error, 1)

			go func() {
				chapters, err := bookmarkService.FetchEpisodes(req.URL)
				if err != nil {
					errChan <- err
				} else {
					done <- chapters
				}
			}()

			select {
			case chapters := <-done:
				c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"chapters": chapters, "count": len(chapters)}, "error": ""})
			case err := <-errChan:
				log.Printf("[TestScrape] error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "error": err.Error()})
			case <-ctx.Done():
				c.JSON(http.StatusGatewayTimeout, gin.H{"success": false, "data": nil, "error": "Scraping timeout"})
			}
		})

		// Recherche d'une série par URL ou query
		series.GET("/search", seriesHandler.SearchSeries)
	}

	// Redirection vers un chapitre spécifique
	router.GET("/chapter/:id", seriesHandler.RedirectToChapter)

	// Redirection vers l'URL originale d'une série
	router.GET("/series/redirect", seriesHandler.RedirectToSeries)
}

// ==========================
// HELPERS
// ==========================

func getSeriesID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "error": "Series ID is required"})
		return "", false
	}
	return id, true
}

func containsIgnoreCase(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}
