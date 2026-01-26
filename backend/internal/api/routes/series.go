package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/middleware"
)

func SeriesRoutes(router *gin.Engine) {
	series := router.Group("/series")
	series.Use(middleware.AuthMiddleware())
	{
		series.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{"series": []string{}})
		})
		series.GET("/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{"series": nil})
		})
		series.POST("", func(c *gin.Context) {
			c.JSON(201, gin.H{"message": "series created"})
		})
	}
}
