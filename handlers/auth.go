package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"merch-store/config"
	"merch-store/models"
	"merch-store/repository"
)

type AuthHandler struct {
	UserRepo repository.UserRepository
	Config   *config.Config
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

// Auth выполняет аутентификацию и возвращает JWT-токен.
func (h *AuthHandler) Auth(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserRepo.GetByUsername(req.Username)
	if err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user = &models.User{
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
			Balance:      1000,
		}
		err = h.UserRepo.Create(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать пользователя"})
			return
		}
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверные учетные данные"})
			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать токен"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: tokenString})
}
