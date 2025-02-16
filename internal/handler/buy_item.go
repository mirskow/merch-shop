package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MerchType string

const (
	MerchTypeTShirt    MerchType = "t-shirt"
	MerchTypeCup       MerchType = "cup"
	MerchTypeBook      MerchType = "book"
	MerchTypePen       MerchType = "pen"
	MerchTypePowerBank MerchType = "powerbank"
	MerchTypeHoody     MerchType = "hoody"
	MerchTypeUmbrella  MerchType = "umbrella"
	MerchTypeSocks     MerchType = "socks"
	MerchTypeWallet    MerchType = "wallet"
	MerchTypePinkHoody MerchType = "pink-hoody"
)

func validateReqMerchType(item MerchType) bool {
	switch item {
	case MerchTypeTShirt, MerchTypeCup, MerchTypeBook, MerchTypePen, MerchTypePowerBank,
		MerchTypeHoody, MerchTypeUmbrella, MerchTypeSocks, MerchTypeWallet, MerchTypePinkHoody:
		return true
	default:
		return false
	}
}

func (h *Handler) buyItem(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "Внутренняя ошибка сервера")
		return
	}

	item := c.Param("item")
	if !validateReqMerchType(MerchType(item)) {
		newErrorResponse(c, http.StatusBadRequest, "Не получилось считать параметр item", "Неверный запрос")
		return
	}

	err = h.services.Buyer.Buy(c.Request.Context(), userID, item)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "Внутренняя ошибка сервера")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}
