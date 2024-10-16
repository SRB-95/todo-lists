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
