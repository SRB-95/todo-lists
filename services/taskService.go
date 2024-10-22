package services

import (
	"time"
	"todo-lists/entity"
	"todo-lists/repositories"
)

type TaskService struct {
	Repo *repositories.TaskRepository
}

// CreateTask method creates a new task
func (s *TaskService) CreateTask(task *entity.Task) error {
	return s.Repo.CreateTask(task)
}

// GetAllTasks method retrieves all tasks
func (s *TaskService) GetAllTasks() ([]entity.Task, error) {
	return s.Repo.GetAllTasks()
}

// GetTaskById method retrieves a task by ID
func (s *TaskService) GetTaskById(id int) (entity.Task, error) {
	return s.Repo.GetTaskById(id)
}

// GetTaskByTag method retrieves a task by tag name
func (s *TaskService) GetTasksByTag(tag string) ([]entity.Task, error) {
	return s.Repo.GetTasksByTag(tag)
}

// UpdateTask method updates an existing task
func (s *TaskService) UpdateTask(task *entity.Task) error {
	return s.Repo.UpdateTask(task)
}

// SearchTasksByName method searches for tasks by keyword in their name
func (s *TaskService) SearchTasksByName(keyword string) ([]entity.Task, error) {
	return s.Repo.SearchTasksByName(keyword)
}

// FilterTasksByDeadline method retrieves tasks within a specified date range
func (s *TaskService) FilterTasksByDeadline(start, end time.Time) ([]entity.Task, error) {
	return s.Repo.FilterTasksByDeadline(start, end)
}

// DeleteTask deletes a task by its ID
func (s *TaskService) DeleteTask(id int) error {
	return s.Repo.DeleteTask(id)
}
