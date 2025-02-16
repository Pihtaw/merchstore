package handlers

import (
	"net/http"
  "log"

	"github.com/gin-gonic/gin"
	"merch-store/models"
	"merch-store/repository"
)

type MerchHandler struct {
	MerchRepo    repository.MerchRepository
	PurchaseRepo repository.PurchaseRepository
	UserRepo     repository.UserRepository
}

// BuyItem выполняет покупку мерча по названию товара.
// Эндпоинт: GET /api/buy/:item
func (h *MerchHandler) BuyItem(c *gin.Context) {
	itemName := c.Param("item")
	merch, err := h.MerchRepo.GetByName(itemName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "товар не найден"})
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неавторизованный доступ"})
		return
	}
	userID := userIDInterface.(int)
	user, err := h.UserRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "пользователь не найден"})
		return
	}
	if user.Balance < merch.Price {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недостаточно монет"})
		return
	}

	// Списываем стоимость товара с баланса
	newBalance := user.Balance - merch.Price
	if err := h.UserRepo.UpdateBalance(userID, newBalance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить баланс"})
		return
	}

	// Регистрируем покупку
	purchase := models.Purchase{
		UserID:  userID,
		MerchID: merch.ID,
	}
	if err := h.PurchaseRepo.Create(&purchase); err != nil {
    log.Println("Ошибка при создании покупки:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "покупка успешна", "new_balance": newBalance})
}
