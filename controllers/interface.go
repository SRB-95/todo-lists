package controllers

import "github.com/gin-gonic/gin"

// TaskControllerInterface defines the methods that a TaskController should implement.
type IController interface {
	CreateTask(ctx *gin.Context)
	GetTasks(ctx *gin.Context)
	GetTaskById(ctx *gin.Context)
	GetTaskByTag(ctx *gin.Context)
	UpdateTask(ctx *gin.Context)
	SearchTasks(ctx *gin.Context)
	FilterTasksByDeadline(ctx *gin.Context)
	DeleteTask(ctx *gin.Context)
}
