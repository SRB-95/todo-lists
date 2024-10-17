package repositories

import (
	"log"
	"time"
	"todo-lists/models"

	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

// CreateTask saves a new task in the database
func (r *TaskRepository) CreateTask(task *models.Task) error {
	return r.DB.Create(task).Error
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

// GetTaskById method retrieves a task by ID from the database
func (r *TaskRepository) GetTaskById(id int) (models.Task, error) {
	var task models.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return task, err
		}
		log.Println("Error fetching task:", err)
		return task, err
	}
	return task, nil
}

// GetTaskByTag method retrieves a task by tag name from the database
func (r *TaskRepository) GetTasksByTag(tag string) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB.Where("tag = ?", tag).Find(&tasks).Error; err != nil {
		log.Println("Error fetching tasks by tag:", err)
		return nil, err
	}
	return tasks, nil
}

// UpdateTask method updates a task in the database
func (r *TaskRepository) UpdateTask(task *models.Task) error {
	return r.DB.Save(task).Error
}

// SearchTasksByName method searches for tasks by keyword in their name
func (r *TaskRepository) SearchTasksByName(keyword string) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB.Where("name LIKE ?", "%"+keyword+"%").Find(&tasks).Error; err != nil {
		log.Println("Error searching tasks by keyword:", err)
		return nil, err
	}
	return tasks, nil
}

// FilterTasksByDeadline method retrieves tasks with deadlines within the specified range
func (r *TaskRepository) FilterTasksByDeadline(start, end time.Time) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB.Where("deadline BETWEEN ? AND ?", start, end).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// DeleteTask deletes a task by its ID
func (r *TaskRepository) DeleteTask(id int) error {
	result := r.DB.Delete(&models.Task{}, id)
	return result.Error
}
