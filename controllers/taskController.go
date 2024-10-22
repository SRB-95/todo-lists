package controllers

import (
	"net/http"
	"strconv"
	"time"
	"todo-lists/entity"
	"todo-lists/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	Service *services.TaskService
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	var task entity.Task

	// Bind the incoming JSON to the task struct
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := c.Service.CreateTask(&task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating task"})
		return
	}

	ctx.JSON(http.StatusCreated, task)
}

// GetTasks method retrieves all tasks and responds with JSON
func (c *TaskController) GetTasks(ctx *gin.Context) {
	tasks, err := c.Service.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tasks"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// GetTaskById method retrieves a task by ID and responds with JSON
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

// UpdateTask method handles the updating of an existing task
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task entity.Task
	// Bind the incoming JSON to the task struct
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set the ID on the task to ensure we update the correct one
	task.ID = uint(id)

	if err := c.Service.UpdateTask(&task); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
		}
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// SearchTasks handles searching for tasks by keyword in the name
func (c *TaskController) SearchTasks(ctx *gin.Context) {
	keyword := ctx.Query("keyword")

	tasks, err := c.Service.SearchTasksByName(keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching tasks"})
		return
	}

	if len(tasks) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No tasks found matching with name"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// FilterTasksByDeadline method filters the tasks by a deadline within a specified date range
func (c *TaskController) FilterTasksByDeadline(ctx *gin.Context) {
	startDateStr := ctx.Query("start")
	endDateStr := ctx.Query("end")

	// format specifier, Go understands that you are expecting the input to be in the format YYYY-MM-DD
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	tasks, err := c.Service.FilterTasksByDeadline(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error filtering tasks"})
		return
	}

	if len(tasks) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No tasks found in the specified date range"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// DeleteTask method deletes a task by ID
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = c.Service.DeleteTask(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting task"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"}) // Respond with success message
}
