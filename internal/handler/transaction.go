package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func validateReqAmount(amount int) bool {
	return amount > 0
}

func (h *Handler) sendCoin(c *gin.Context) {
	fromUserID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "Внутренняя ошибка сервера")
		return
	}

	var req struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}

	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "Неверный запрос")
		return
	}

	if !validateReqAmount(req.Amount) {
		newErrorResponse(c, http.StatusBadRequest, "Negative coin value for transfer", "Неверный запрос.")
		return
	}

	err = h.services.Transaction.SendCoin(c.Request.Context(), fromUserID, req.ToUser, req.Amount)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "Внутренняя ошибка сервера.")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}
