package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/services"
)

type SeriesHandler struct {
	BookmarkService *services.BookmarkSeriesService
}

func NewSeriesHandler(bookmarkService *services.BookmarkSeriesService) *SeriesHandler {
	return &SeriesHandler{BookmarkService: bookmarkService}
}

// SearchSeries recherche une série par son nom (query)
func (h *SeriesHandler) SearchSeries(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query"})
		return
	}

	// Utilisation du service BookmarkSeriesService pour rechercher la série
	series, err := h.BookmarkService.FetchSeries(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search series"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": series})
}

// RedirectToSeries redirige l'utilisateur vers une URL spécifique
func (h *SeriesHandler) RedirectToSeries(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing url"})
		return
	}

	c.Redirect(http.StatusFound, url)
}
