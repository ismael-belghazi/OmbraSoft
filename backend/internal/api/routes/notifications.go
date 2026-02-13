package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/db"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/middleware"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/models"
)

type NotificationSettingsRequest struct {
	Push      bool   `json:"push"`
	DiscordID string `json:"discord_id"`
}

func GetNotificationSettings(c *gin.Context) {
	userID := c.GetString("userID")
	var settings models.UserNotifications

	database := db.GetDB()
	if err := database.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"push": true, "discord_id": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"push":       settings.Push,
		"discord_id": settings.DiscordID,
	})
}

func UpdateNotificationSettings(c *gin.Context) {
	userID := c.GetString("userID")
	var req NotificationSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	database := db.GetDB()
	var settings models.UserNotifications

	if err := database.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		settings = models.UserNotifications{
			ID:        uuid.NewString(),
			UserID:    userID,
			Push:      req.Push,
			DiscordID: req.DiscordID,
		}
		database.Create(&settings)
	} else {
		settings.Push = req.Push
		settings.DiscordID = req.DiscordID
		database.Save(&settings)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Préférences mises à jour"})
}

func NotificationRoutes(router *gin.RouterGroup) {
	notify := router.Group("/user")
	notify.Use(middleware.AuthMiddleware())
	{
		notify.GET("/notifications", GetNotificationSettings)
		notify.POST("/notifications", UpdateNotificationSettings)
	}
}
