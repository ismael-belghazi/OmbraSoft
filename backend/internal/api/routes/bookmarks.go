package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/models"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/services"
)

func BookmarksRoutes(router *gin.RouterGroup, bookmarkService *services.BookmarkSeriesService) {
	bookmarks := router.Group("/bookmarks")
	{
		bookmarks.GET("", func(c *gin.Context) {
			userID := c.GetString("userID")

			bookmarks, err := bookmarkService.GetBookmarksByUser(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookmarks"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"bookmarks": bookmarks})
		})

		bookmarks.POST("", func(c *gin.Context) {
			var newBookmark struct {
				SeriesID string `json:"series_id"`
			}

			userID := c.GetString("user_id")

			if err := c.ShouldBindJSON(&newBookmark); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			err := bookmarkService.AddBookmarkWithChapters(userID, models.Series{ID: newBookmark.SeriesID})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add bookmark"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Bookmark added successfully"})
		})

		bookmarks.DELETE("/:id", func(c *gin.Context) {
			bookmarkID := c.Param("id")
			userID := c.GetString("user_id")

			if err := bookmarkService.DeleteBookmarkForUser(bookmarkID, userID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bookmark"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Bookmark deleted successfully"})
		})
	}
}
