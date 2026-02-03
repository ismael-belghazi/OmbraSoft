package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/db"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/models"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	SecretPhrase string `json:"secret_phrase" binding:"required,min=10"`
}

type ForgotPasswordRequest struct {
	Email        string `json:"email" binding:"required,email"`
	SecretPhrase string `json:"secret_phrase" binding:"required,min=10"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var existing models.User
	err := db.GetDB().Where("email = ?", req.Email).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email déjà utilisé"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	hashedSecret, _ := bcrypt.GenerateFromPassword([]byte(req.SecretPhrase), bcrypt.DefaultCost)

	user := &models.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		SecretHash:   string(hashedSecret),
	}

	if err := db.GetDB().Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email)
	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := db.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email)
	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", Register)
		auth.POST("/login", Login)
	}
}
