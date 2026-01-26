package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bebeb/ombrasoft-backend/internal/api/routes"
	"github.com/bebeb/ombrasoft-backend/internal/config"
	"github.com/bebeb/ombrasoft-backend/internal/db"
	"github.com/bebeb/ombrasoft-backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := db.Init(); err != nil {
		log.Fatalf("Erreur DB: %v", err)
	}

	if cfg.RedisURL != "" {
		if err := services.InitRedis(); err != nil {
			log.Fatalf("Erreur Redis: %v", err)
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
	routes.SeriesRoutes(r)
	routes.BookmarksRoutes(r)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Backend démarré sur port %s (Mode: %s)", cfg.Port, cfg.GINMode)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("⏸  Arrêt gracieux...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Erreur shutdown: %v", err)
		}

		_ = db.Close()
		_ = services.CloseRedis()
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Erreur serveur: %v", err)
	}
}
