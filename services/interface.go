package services

import (
	"time"
	"todo-lists/entity"
)

type IService interface {
	CreateTask(task *entity.Task) error
	GetAllTasks() ([]entity.Task, error)
	GetTaskById(id int) (entity.Task, error)
	GetTasksByTag(tag string) ([]entity.Task, error)
	UpdateTask(task *entity.Task) error
	SearchTasksByName(keyword string) ([]entity.Task, error)
	FilterTasksByDeadline(start, end time.Time) ([]entity.Task, error)
	DeleteTask(id int) error
}
