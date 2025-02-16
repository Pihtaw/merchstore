package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"merch-store/config"
	"merch-store/handlers"
	"merch-store/repository"
	_ "github.com/lib/pq"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	cfg := config.LoadConfig()

	// Подключаемся к базе данных
	db, err := sqlx.Connect("postgres", cfg.DBSource)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}
	defer db.Close()

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(db)
	merchRepo := repository.NewMerchRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Инициализируем обработчики
	authHandler := &handlers.AuthHandler{
		UserRepo: userRepo,
		Config:   cfg,
	}
	transactionHandler := &handlers.TransactionHandler{
		UserRepo:        userRepo,
		TransactionRepo: transactionRepo,
	}
	merchHandler := &handlers.MerchHandler{
		MerchRepo:    merchRepo,
		PurchaseRepo: purchaseRepo,
		UserRepo:     userRepo,
	}
	infoHandler := &handlers.InfoHandler{
		UserRepo:        userRepo,
		PurchaseRepo:    purchaseRepo,
		TransactionRepo: transactionRepo,
		MerchRepo:       merchRepo,
	}

	// Инициализируем Gin и маршруты
	router := gin.Default()

	// Группа публичных маршрутов (без авторизации)
	api := router.Group("/api")
	{
		// Аутентификация: POST /api/auth
		api.POST("/auth", authHandler.Auth)
	}

	// Группа маршрутов, требующих авторизации (JWT)
	authorized := api.Group("/")
	authorized.Use(AuthMiddleware(cfg.JWTSecret))
	{
		// Перевод монет: POST /api/sendCoin
		authorized.POST("/sendCoin", transactionHandler.SendCoin)
		// Покупка мерча по названию: GET /api/buy/:item
		authorized.GET("/buy/:item", merchHandler.BuyItem)
		// Получение информации о монетах, инвентаре и истории транзакций: GET /api/info
		authorized.GET("/info", infoHandler.Info)
	}

	log.Println("Server is running on port " + cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
}

// AuthMiddleware проверяет JWT токен и устанавливает userID в контекст
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "отсутствует токен"})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена"})
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "неверный токен"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "неверный токен"})
			return
		}
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "неверный токен"})
			return
		}
		c.Set("userID", int(userIDFloat))
		c.Next()
	}
}
