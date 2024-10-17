package routing

import (
	"log"
	"todo-lists/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer(taskController *controllers.TaskController) {
	router := gin.Default()

	// Task API
	router.POST("/tasks", taskController.CreateTask)
	router.GET("/tasks", taskController.GetTasks)
	router.GET("/tasks/:id", taskController.GetTaskById)
	router.GET("/tasks/tag/:tag", taskController.GetTaskByTag)
	router.PUT("/tasks/:id", taskController.UpdateTask)
	router.GET("/tasks/search", taskController.SearchTasks)
	router.GET("/tasks/filter", taskController.FilterTasksByDeadline)
	router.DELETE("/tasks/:id", taskController.DeleteTask)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
