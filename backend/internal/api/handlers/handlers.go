package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/models"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

var users = make(map[uuid.UUID]*models.User)

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for _, user := range users {
		if user.Email == req.Email {
			c.JSON(400, gin.H{"error": "Email déjà utilisé"})
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erreur lors du hashage"})
		return
	}

	user := &models.User{
		ID:           uuid.New(), // ✅ FIX
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	users[user.ID] = user

	token, err := utils.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(201, AuthResponse{
		Token: token,
		User: &models.User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user *models.User
	for _, u := range users {
		if u.Email == req.Email {
			user = u
			break
		}
	}

	if user == nil {
		c.JSON(401, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(200, AuthResponse{
		Token: token,
		User: &models.User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
