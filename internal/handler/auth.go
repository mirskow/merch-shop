package handler

import (
	"avito-internship/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) auth(c *gin.Context) {
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "Неверный запрос.")
		return
	}

	token, err := h.services.Authorization.CreateUser(c.Request.Context(), user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "Внутренняя ошибка сервера.")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}
