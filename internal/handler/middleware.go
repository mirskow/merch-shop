package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Authorization header is missing", "Неавторизован.")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid format in Authorization header", "Неавторизован.")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "Token is missing", "Неавторизован.")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error(), "Неавторизован.")
		return
	}

	c.Set("userID", userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userID")
	if !ok {
		return 0, errors.New("user id not found at context")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is not the correct type")
	}

	return idInt, nil
}
