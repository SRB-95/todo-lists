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
