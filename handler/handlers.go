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
	router.GET("/get-all-tasks", h.getAllTasks)
	router.GET("/get-task-by-id", h.getTaskById)
	router.PATCH("/update-task", h.updateTask)
	router.DELETE("/delete-task", h.deleteTask)

	return router
}
