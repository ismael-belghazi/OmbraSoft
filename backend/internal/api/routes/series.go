package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/middleware"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/services"
)

func SeriesRoutes(router *gin.RouterGroup, bookmarkService *services.BookmarkSeriesService) {
	series := router.Group("/series")
	series.Use(middleware.AuthMiddleware())
	{
		series.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"series": []string{}})
		})

		series.GET("/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"series": nil})
		})

		series.POST("", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"message": "series created"})
		})

		series.GET("/chapters", func(c *gin.Context) {
			seriesURL := c.DefaultQuery("url", "")
			if seriesURL == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
				return
			}

			chapters, err := bookmarkService.FetchChapters(seriesURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chapters"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"chapters": chapters})
		})
	}
}
