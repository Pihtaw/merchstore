package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"merch-store/models"
	"merch-store/repository"
)

type InfoHandler struct {
	UserRepo        repository.UserRepository
	PurchaseRepo    repository.PurchaseRepository
	TransactionRepo repository.TransactionRepository
	MerchRepo       repository.MerchRepository
}

// GET /api/info
func (h *InfoHandler) Info(c *gin.Context) {
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

	purchases, err := h.PurchaseRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить покупки"})
		return
	}

	transactions, err := h.TransactionRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить транзакции"})
		return
	}

	received := []map[string]interface{}{}
	sent := []map[string]interface{}{}
	for _, t := range transactions {
		if t.ToUserID == userID {
			received = append(received, map[string]interface{}{
				"fromUser": t.FromUserID,
				"amount":   t.Amount,
			})
		}
		if t.FromUserID == userID {
			sent = append(sent, map[string]interface{}{
				"toUser": t.ToUserID,
				"amount": t.Amount,
			})
		}
	}

	inventoryMap := make(map[string]int)
	for _, p := range purchases {
		merch, err := h.MerchRepo.GetByID(p.MerchID)
		if err != nil {
			continue
		}
		inventoryMap[merch.Name]++
	}
	inventory := []models.InventoryItem{}
	for name, qty := range inventoryMap {
		inventory = append(inventory, models.InventoryItem{
			Type:     name,
			Quantity: qty,
		})
	}

	infoResponse := models.InfoResponse{
		Coins:       user.Balance,
		Inventory:   inventory,
		CoinHistory: models.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}
	c.JSON(http.StatusOK, infoResponse)
}
