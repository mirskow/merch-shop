//go:build integration
// +build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"avito-internship/internal/handler"
	"avito-internship/internal/repository"
	"avito-internship/internal/repository/postgres"
	"avito-internship/internal/server"
	"avito-internship/internal/services"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type AuthResponse struct {
	Token string `json:"token"`
}

type BuyResponse struct {
	Message string `json:"message"`
}

var (
	db  *sqlx.DB
	url string
)

func setupTestServer(t *testing.T) *handler.Handler {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file")
	}

	viper.Set("db.host", "localhost")
	viper.Set("db.port", "5433")
	viper.Set("db.username", "testuser")
	viper.Set("db.dbname", "testdb")
	viper.Set("db.sslmode", "disable")
	viper.Set("port", "8000")

	var err error
	db, err = postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD_TEST"),
	})
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	repos := repository.NewRepository(db)
	services := services.NewServices(repos)
	handlers := handler.NewHandler(services)

	return handlers
}

func AuthTest(t *testing.T) string {
	username := "user"
	password := "user"

	purchaseRequest := map[string]interface{}{
		"username": username,
		"password": password,
	}

	requestBody, err := json.Marshal(purchaseRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	resp, err := http.Post("http://localhost:8000/api/auth", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to send auth request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp AuthResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&authResp)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.NotEmpty(t, authResp.Token, "Token should not be empty")

	var count int
	err = db.QueryRow("SELECT count(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query the database (username: %s): %v", username, err)
	}

	assert.Equal(t, 1, count, "User should be recorded in the database")

	return authResp.Token
}

func SendCoinTest(t *testing.T) {
	token := AuthTest(t)
	toUser := "test"
	amount := 20

	sendRequest := map[string]interface{}{
		"toUser": toUser,
		"amount": amount,
	}

	requestBody, err := json.Marshal(sendRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/api/sendCoin", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var count int
	err = db.QueryRow("SELECT count(*) FROM transactions JOIN users ON users.id = transactions.to_user WHERE users.username = $1", toUser).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query the database (username: %s): %v", toUser, err)
	}

	assert.Equal(t, 1, count, "Transaction should be recorded in the database")
}

func BuyItemTest(t *testing.T) {
	token := AuthTest(t)
	item := "cup"

	req, err := http.NewRequest("GET", "http://localhost:8000/api/buy/"+item, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var count int
	err = db.QueryRow("SELECT count(*) FROM purchases JOIN merch ON purchases.item_id = merch.id WHERE merch.item = $1", item).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query the database (username: %s): %v", item, err)
	}

	assert.Equal(t, 1, count, "Purchases should be recorded in the database")
}

func TestIntegration(t *testing.T) {
	// запуск тестового сервера
	handlers := setupTestServer(t)

	defer db.Close()

	srv := new(server.Server)

	errChan := make(chan error, 1)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {

			errChan <- fmt.Errorf("Error while running HTTP server: %v", err)
		}
	}()

	select {
	case err := <-errChan:
		t.Fatalf("Test failed: %v", err)
	case <-time.After(2 * time.Second):

	}

	t.Run("Send Coin", SendCoinTest)
	t.Run("Buy Item", BuyItemTest)
}
