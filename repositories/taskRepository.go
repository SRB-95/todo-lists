package repositories

import (
	"log"
	"todo-lists/models"

	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

// GetAllTasks fetches all tasks from the database
func (r *TaskRepository) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB.Find(&tasks).Error; err != nil {
		log.Println("Error fetching tasks:", err)
		return nil, err
	}
	return tasks, nil
}

// GetTaskById retrieves a task by ID from the database
func (r *TaskRepository) GetTaskById(id int) (models.Task, error) {
	var task models.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return task, err // No record found
		}
		log.Println("Error fetching task:", err)
		return task, err // Other errors
	}
	return task, nil // Return the found task
}

// GetTaskByTag retrieves a task by tag name from the database
func (r *TaskRepository) GetTasksByTag(tag string) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB.Where("tag = ?", tag).Find(&tasks).Error; err != nil {
		log.Println("Error fetching tasks by tag:", err)
		return nil, err
	}
	return tasks, nil
}
