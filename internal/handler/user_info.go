package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getInfo(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "Внутренняя ошибка сервера")
		return
	}

	coins, purchases, sender, receiver, err := h.services.User.GetUsersInfo(c.Request.Context(), userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "Внутренняя ошибка сервера")
		return
	}

	newInfoResponse(c, coins, purchases, sender, receiver)
}
