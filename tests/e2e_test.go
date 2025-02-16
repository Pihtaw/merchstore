package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"merch-store/config"
	"merch-store/handlers"
	"merch-store/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupRouter() *gin.Engine {
	cfg := config.LoadConfig()
	db, err := sqlx.Connect("postgres", cfg.DBSource)
	if err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	authHandler := &handlers.AuthHandler{
		UserRepo: userRepo,
		Config:   cfg,
	}

	transactionHandler := &handlers.TransactionHandler{
		TransactionRepo: transactionRepo,
	}

	router := gin.Default()
	router.POST("/auth", authHandler.Auth)
	router.POST("/sendCoin", transactionHandler.SendCoin) // Добавляем эндпоинт перевода

	return router
}


func TestRegisterAndLogin(t *testing.T) {
	router := setupRouter()

	registerPayload := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	payloadBytes, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/auth", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.Code)
	}

	req, _ = http.NewRequest("POST", "/auth", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}
func TestTransferCoins(t *testing.T) {
	router := setupRouter()

	// Регистрируем двух пользователей
	user1 := map[string]string{"username": "alice", "password": "password123"}
	user2 := map[string]string{"username": "bob", "password": "password123"}

	register := func(user map[string]string) *httptest.ResponseRecorder {
		payload, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/auth", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		return resp
	}

	register(user1)
	register(user2)

	// Логиним Alice и получаем токен
	login := func(user map[string]string) string {
		payload, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/auth", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("Login failed for %s", user["username"])
		}

		var res map[string]string
		json.Unmarshal(resp.Body.Bytes(), &res)
		return res["token"]
	}

	aliceToken := login(user1)

	// Перевод монет от Alice к Bob
	transferPayload := map[string]interface{}{
		"from_user": "alice",
		"to_user":   "bob",
		"amount":    10,
	}
	payloadBytes, _ := json.Marshal(transferPayload)

	req, _ := http.NewRequest("POST", "/sendCoin", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+aliceToken)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Transfer failed, expected 200, got %d", resp.Code)
	}
}
