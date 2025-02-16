package handler

import (
	"avito-internship/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.POST("/auth", h.auth)

		protected := api.Group("/", h.userIdentity)
		{
			protected.GET("/info", h.getInfo)
			protected.POST("/sendCoin", h.sendCoin)
			protected.GET("/buy/:item", h.buyItem)
		}
	}

	return router
}
