package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Task created"})
}

func (h *Handler) getAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "All tasks"})
}

func (h *Handler) getTaskById(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Task by ID"})
}

func (h *Handler) updateTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (h *Handler) deleteTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

/*
// Связываем входящие данные с структурой json
	if err := c.BindJSON(&json); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Создаем задачу через сервис
	if err := h.services.Create(json.ChatID, json.Task, json.Interval, json.Unit, json.Username); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Возвращаем ответ
	c.JSON(http.StatusOK, gin.H{"status": "task scheduled"})
*/
