package handler

import (
	"telegram_reminder_bot/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/create-task", h.createTask)
	router.GET("/tasks", h.Tasks)

	return router
}
