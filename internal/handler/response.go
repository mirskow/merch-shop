package handler

import (
	"avito-internship/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, errorMessage string, clientMessage string) {
	logrus.Error(errorMessage)

	c.AbortWithStatusJSON(statusCode, ErrorResponse{clientMessage})
}

type InfoResponse struct {
	Coins       int         `json:"coins" binding:"required"`
	Inventory   []Inventory `json:"inventory" binding:"required" `
	CoinHistory CoinHistory `json:"coinHistory" binding:"required" `
}

type Inventory struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []CoinTransaction `json:"received"`
	Sent     []CoinTransaction `json:"sent"`
}

type CoinTransaction struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

func newInfoResponse(c *gin.Context, coins int, purchases []model.Purchases, sender []model.Transaction, receiver []model.Transaction) {
	var inventory []Inventory

	for _, v := range purchases {
		inventory = append(inventory, Inventory{
			Item:     v.ItemName,
			Quantity: v.Quantity,
		})
	}

	var history CoinHistory

	for _, v := range sender {
		history.Sent = append(history.Sent, CoinTransaction{
			Username: v.ReceiverName,
			Amount:   v.Amount,
		})
	}

	for _, v := range receiver {
		history.Received = append(history.Received, CoinTransaction{
			Username: v.SenderName,
			Amount:   v.Amount,
		})
	}

	c.JSON(http.StatusOK, InfoResponse{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: history,
	})
}
