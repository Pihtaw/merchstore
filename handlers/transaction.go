package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"merch-store/models"
	"merch-store/repository"
)

type TransactionHandler struct {
	UserRepo        repository.UserRepository
	TransactionRepo repository.TransactionRepository
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required,gt=0"`
}

// SendCoin выполняет перевод монет от аутентифицированного пользователя другому пользователю.
func (h *TransactionHandler) SendCoin(c *gin.Context) {
	var req SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	senderIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неавторизованный доступ"})
		return
	}
	senderID := senderIDInterface.(int)

	// Получаем данные получателя по имени
	recipient, err := h.UserRepo.GetByUsername(req.ToUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "получатель не найден"})
		return
	}

	// Проверяем баланс отправителя
	sender, err := h.UserRepo.GetByID(senderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения данных пользователя"})
		return
	}
	if sender.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недостаточно средств"})
		return
	}

	// Обновляем балансы
	if err := h.UserRepo.UpdateBalance(senderID, sender.Balance-req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить баланс отправителя"})
		return
	}
	newRecipientBalance := recipient.Balance + req.Amount
	if err := h.UserRepo.UpdateBalance(recipient.ID, newRecipientBalance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить баланс получателя"})
		return
	}

	// Регистрируем транзакцию
	transaction := models.Transaction{
		FromUserID: senderID,
		ToUserID:   recipient.ID,
		Amount:     req.Amount,
	}
  if err := h.TransactionRepo.Create(&transaction); err != nil {
  	c.JSON(http.StatusInternalServerError, gin.H{
  		"error": "не удалось создать транзакцию",
  		"details": err.Error(),  // Печатаем точную ошибку
  	})
  	return
  }

	c.JSON(http.StatusOK, gin.H{"message": "перевод выполнен"})
}
