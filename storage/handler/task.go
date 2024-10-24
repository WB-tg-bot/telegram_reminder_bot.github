package handler

import (
	"net/http"
	"telegram_reminder_bot/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createTask(c *gin.Context) {
	var input models.Task
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Task.CreateTask(input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task created"})
}

func (h *Handler) Tasks(c *gin.Context) {
	output, err := h.services.Task.Tasks()
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, output)
}
