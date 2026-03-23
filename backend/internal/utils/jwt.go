package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ismael-belghazi/ombrasoft-backend/internal/config"
)

var resetTokens = struct {
	sync.RWMutex
	tokens map[string]string
}{tokens: make(map[string]string)}

func SaveResetToken(userID, token string) {
	resetTokens.Lock()
	defer resetTokens.Unlock()
	resetTokens.tokens[userID] = token
}

func GetResetToken(userID string) (string, bool) {
	resetTokens.RLock()
	defer resetTokens.RUnlock()
	t, ok := resetTokens.tokens[userID]
	return t, ok
}

func DeleteResetToken(userID string) {
	resetTokens.Lock()
	defer resetTokens.Unlock()
	delete(resetTokens.tokens, userID)
}

type AppriseRequest struct {
	Title   string `json:"title"`
	Message string `json:"body"`
	Target  string `json:"target"`
}

func SendAppriseEmail(to, subject, message string) error {
	appriseURL := os.Getenv("APPRISE_URL")
	if appriseURL == "" {
		return errors.New("APPRISE_URL non défini")
	}

	payload := AppriseRequest{
		Title:   subject,
		Message: message,
		Target:  to,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(appriseURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("erreur lors de l'envoi via Apprise")
	}

	return nil
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := config.GetJWTSecret()
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	jwtSecret := config.GetJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
