package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/api/routes"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/config"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/db"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/middleware"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/services"
)

func main() {
	cfg := config.Load()

	if err := db.Init(); err != nil {
		log.Fatalf("Database error: %v", err)
	}

	if cfg.RedisURL != "" {
		if err := services.InitRedis(); err != nil {
			log.Printf("Error initializing Redis: %v", err)
		}
	}

	if cfg.GINMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	routes.AuthRoutes(r)

	bookmarkService := services.NewBookmarkSeriesService(db.DB)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	routes.SeriesRoutes(protected, bookmarkService)
	routes.BookmarksRoutes(protected, bookmarkService)
	routes.NotificationRoutes(protected)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Backend started on port %s (mode: %s)", cfg.Port, cfg.GINMode)

	// Gestion de l'arrêt du serveur (graceful shutdown)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}

		if err := db.Close(); err != nil {
			log.Printf("Error closing DB connection: %v", err)
		}

		if err := services.CloseRedis(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
