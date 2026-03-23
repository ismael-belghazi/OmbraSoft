package routes

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/models"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/services"
)

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

func BookmarksRoutes(router *gin.RouterGroup, bookmarkService *services.BookmarkSeriesService) {
	bookmarks := router.Group("/bookmarks")
	{
		// Récupérer tous les bookmarks d'un utilisateur
		bookmarks.GET("", func(c *gin.Context) {
			userID, ok := getUserID(c)
			if !ok {
				return
			}

			result, err := bookmarkService.GetBookmarksByUser(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": result, "error": ""})
		})

		// Ajouter un nouveau bookmark
		bookmarks.POST("", func(c *gin.Context) {
			userID, ok := getUserID(c)
			if !ok {
				return
			}

			var payload struct {
				Title     string `json:"title"`
				SourceURL string `json:"sourceURL" binding:"required"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "error": "Invalid body: " + err.Error()})
				return
			}

			payload.Title = strings.TrimSpace(payload.Title)
			payload.SourceURL = strings.TrimSpace(payload.SourceURL)
			if _, ok := isValidHTTPURL(payload.SourceURL); !ok {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "error": "Invalid SourceURL"})
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				series := models.Series{Title: payload.Title, SourceURL: payload.SourceURL}
				done <- bookmarkService.AddBookmarkWithChapters(userID, &series)
			}()

			select {
			case err := <-done:
				if err != nil {
					msg := err.Error()
					status := http.StatusInternalServerError
					lower := strings.ToLower(msg)
					if strings.Contains(lower, "cloudflare") {
						status = http.StatusBadGateway
						msg = "Blocked by target site (Cloudflare)"
					} else if strings.Contains(lower, "no episodes") {
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

		// Supprimer un bookmark
		bookmarks.DELETE("/:id", func(c *gin.Context) {
			userID, ok := getUserID(c)
			if !ok {
				return
			}

			bookmarkID := strings.TrimSpace(c.Param("id"))
			if bookmarkID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "error": "Bookmark ID is required"})
				return
			}

			if err := bookmarkService.DeleteBookmarkAndSeriesForUser(bookmarkID, userID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": nil, "error": ""})
		})

		// Marquer un chapitre comme lu
		bookmarks.PATCH("/:bookmarkID/chapters/:chapterNumber/read", func(c *gin.Context) {
			userID, ok := getUserID(c)
			if !ok {
				return
			}

			bookmarkID := c.Param("bookmarkID")
			chapterNumber, err := strconv.Atoi(c.Param("chapterNumber"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid chapter number"})
				return
			}

			lastRead, err := bookmarkService.MarkChapterAsRead(userID, bookmarkID, chapterNumber)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"lastReadChapter": lastRead}})
		})

		// Marquer toute la série comme lue
		bookmarks.PATCH("/:bookmarkID/series/read", func(c *gin.Context) {
			userID, ok := getUserID(c)
			if !ok {
				return
			}

			bookmarkID := c.Param("bookmarkID")

			lastRead, err := bookmarkService.MarkSeriesAsRead(userID, bookmarkID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"lastReadChapter": lastRead}})
		})
	}
}

// ========================
// HELPERS
// ========================

func getUserID(c *gin.Context) (string, bool) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "data": nil, "error": "Unauthorized"})
		return "", false
	}
	return userID, true
}
