package routing

import (
	"log"
	"todo-lists/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer(taskController *controllers.TaskController) {
	router := gin.Default()

	// Booking API
	router.GET("/tasks", taskController.GetTasks)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
