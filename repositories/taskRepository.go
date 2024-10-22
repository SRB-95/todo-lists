package repositories

import (
	"log"
	"time"
	"todo-lists/entity"
	"todo-lists/models"

	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

// CreateTask saves a new task in the database
func (r *TaskRepository) CreateTask(task *entity.Task) error {
	newTask := &models.Task{
		Name:     task.Name,
		Deadline: task.Deadline,
		Tag:      task.Tag,
	}

	return r.DB.Create(newTask).Error
}

// GetAllTasks fetches all tasks from the database
func (r *TaskRepository) GetAllTasks() ([]entity.Task, error) {
	var allTasks []models.Task
	if err := r.DB.Find(&allTasks).Error; err != nil {
		log.Println("Error fetching tasks:", err)
		return nil, err
	}

	// Convert allTasks to entity.Tasks
	var entityTasks []entity.Task
	for _, mTask := range allTasks {
		entityTasks = append(entityTasks, entity.Task{
			ID:       mTask.ID,
			Name:     mTask.Name,
			Deadline: mTask.Deadline,
			Tag:      mTask.Tag,
		})
	}

	return entityTasks, nil
}

// GetTaskById method retrieves a task by ID from the database
func (r *TaskRepository) GetTaskById(id int) (entity.Task, error) {
	var task models.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Task{}, err
		}
		log.Println("Error fetching task:", err)
		return entity.Task{}, err
	}
	e := entity.Task{
		ID:       task.ID,
		Name:     task.Name,
		Deadline: task.Deadline,
		Tag:      task.Tag,
	}
	return e, nil
}

// GetTaskByTag method retrieves a task by tag name from the database
func (r *TaskRepository) GetTasksByTag(tag string) ([]entity.Task, error) {
	var tasks []models.Task
	if err := r.DB.Where("tag = ?", tag).Find(&tasks).Error; err != nil {
		log.Println("Error fetching tasks by tag:", err)
		return nil, err
	}
	// Convert allTasks to entity.Tasks
	var entityTasks []entity.Task
	for _, mTask := range tasks {
		entityTasks = append(entityTasks, entity.Task{
			ID:       mTask.ID,
			Name:     mTask.Name,
			Deadline: mTask.Deadline,
			Tag:      mTask.Tag,
		})
	}
	return entityTasks, nil
}

// UpdateTask method updates a task in the database
func (r *TaskRepository) UpdateTask(task *entity.Task) error {
	return r.DB.Save(task).Error
}

// SearchTasksByName method searches for tasks by keyword in their name
func (r *TaskRepository) SearchTasksByName(keyword string) ([]entity.Task, error) {
	var tasks []entity.Task
	if err := r.DB.Where("name LIKE ?", "%"+keyword+"%").Find(&tasks).Error; err != nil {
		log.Println("Error searching tasks by keyword:", err)
		return nil, err
	}
	return tasks, nil
}

// FilterTasksByDeadline method retrieves tasks with deadlines within the specified range
func (r *TaskRepository) FilterTasksByDeadline(start, end time.Time) ([]entity.Task, error) {
	var tasks []entity.Task
	if err := r.DB.Where("deadline BETWEEN ? AND ?", start, end).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// DeleteTask method deletes a task by its ID
func (r *TaskRepository) DeleteTask(id int) error {
	result := r.DB.Delete(&entity.Task{}, id)
	return result.Error
}
