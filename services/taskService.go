package services

import (
	"todo-lists/models"
	"todo-lists/repositories"
)

type TaskService struct {
	Repo *repositories.TaskRepository
}

// GetAllTasks retrieves all tasks from the repository
func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.Repo.GetAllTasks()
}

// GetTaskById retrieves a task by ID from the repository
func (s *TaskService) GetTaskById(id int) (models.Task, error) {
	return s.Repo.GetTaskById(id)
}

// GetTaskByTag retrieves a task by tag name from the repository
func (s *TaskService) GetTasksByTag(tag string) ([]models.Task, error) {
	return s.Repo.GetTasksByTag(tag)
}
