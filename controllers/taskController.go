package controllers

import (
	"net/http"
	"strconv"
	"todo-lists/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// GetTaskById retrieves a task by its ID and responds with JSON
func (c *TaskController) GetTaskById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := c.Service.GetTaskById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching task"})
		}
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// GetTaskByTag retrieves a task by its Tag name and responds with JSON
func (c *TaskController) GetTaskByTag(ctx *gin.Context) {
	tag := ctx.Param("tag")

	tasks, err := c.Service.GetTasksByTag(tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tasks"})
		return
	}

	if len(tasks) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No tasks found with the specified tag"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
