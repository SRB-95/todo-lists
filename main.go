package main

import (
	"log"
	"todo-lists/config"
	"todo-lists/controllers"
	"todo-lists/repositories"
	"todo-lists/routing"
	"todo-lists/services"
)

func main() {
	log.Println("Starting Todo Lists Service..!")

	log.Println("Initializig configuration")
	db := config.ConnectDB()

	// Initialize the repository, service, and controller
	taskRepo := &repositories.TaskRepository{DB: db}
	taskService := &services.TaskService{Repo: taskRepo}
	taskController := &controllers.TaskController{Service: taskService}

	// Start the server with the task controller
	routing.StartServer(taskController)
}
