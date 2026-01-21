package routes

import (
	"github.com/bebeb/ombrasoft-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func BookmarksRoutes(router *gin.Engine) {
	bookmarks := router.Group("/bookmarks")
	bookmarks.Use(middleware.AuthMiddleware())
	{
		bookmarks.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{"bookmarks": []string{}})
		})
		bookmarks.POST("", func(c *gin.Context) {
			c.JSON(201, gin.H{"message": "bookmark created"})
		})
		bookmarks.PATCH("/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "bookmark updated"})
		})
		bookmarks.DELETE("/:id", func(c *gin.Context) {
			c.JSON(204, nil)
		})
	}
}
