package main

import (
	"log"
	"time"

	"github.com/bebeb/ombrasoft-backend/internal/api/routes"
	"github.com/bebeb/ombrasoft-backend/internal/config"
	"github.com/bebeb/ombrasoft-backend/internal/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := db.Init(); err != nil {
		log.Fatalf("Impossible de démarrer: %v", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173", "https://*.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	routes.AuthRoutes(r)
	routes.SeriesRoutes(r)
	routes.BookmarksRoutes(r)

	log.Println("Backend OmbraSoft démarré sur :" + cfg.Port)
	r.Run(":" + cfg.Port)
}
