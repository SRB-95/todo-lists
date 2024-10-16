package controllers

import (
	"net/http"
	"todo-lists/services"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	Service *services.TaskService
}

// GetTasks retrieves all tasks and responds with JSON
func (c *TaskController) GetTasks(ctx *gin.Context) {
	tasks, err := c.Service.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tasks"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
