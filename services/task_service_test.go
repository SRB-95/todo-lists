package services

import (
	"errors"
	"testing"
	"time"
	"todo-lists/entity"
	"todo-lists/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTaskService_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIRepo(ctrl)
	taskService := TaskService{Repo: mockRepo}
	task := &entity.Task{ID: 1, Name: "Test Task"}

	// Test successful creation
	mockRepo.EXPECT().CreateTask(task).Return(nil)
	err := taskService.CreateTask(task)
	assert.NoError(t, err)

	// Test creation error
	mockRepo.EXPECT().CreateTask(task).Return(errors.New("creation error"))
	err = taskService.CreateTask(task)
	assert.Error(t, err)
	assert.Equal(t, "creation error", err.Error())
}

func TestTaskService_GetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIRepo(ctrl)
	taskService := TaskService{Repo: mockRepo}
	tasks := []entity.Task{{ID: 1, Name: "Task 1"}, {ID: 2, Name: "Task 2"}}

	// all tasks successful retrieval
	mockRepo.EXPECT().GetAllTasks().Return(tasks, nil)
	result, err := taskService.GetAllTasks()
	assert.NoError(t, err)
	assert.Equal(t, tasks, result)

	// task retrieval error
	mockRepo.EXPECT().GetAllTasks().Return(nil, errors.New("fetch error"))
	result, err = taskService.GetAllTasks()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "fetch error", err.Error())
}

func TestTaskService_GetTaskById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIRepo(ctrl)
	taskService := TaskService{Repo: mockRepo}
	task := entity.Task{ID: 1, Name: "Test Task"}

	// Task by id successful retrieval
	mockRepo.EXPECT().GetTaskById(1).Return(task, nil)
	result, err := taskService.GetTaskById(1)
	assert.NoError(t, err)
	assert.Equal(t, task, result)

	// Taskt not found error
	mockRepo.EXPECT().GetTaskById(2).Return(entity.Task{}, gorm.ErrRecordNotFound)
	result, err = taskService.GetTaskById(2)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// Service layer error
	mockRepo.EXPECT().GetTaskById(3).Return(entity.Task{}, errors.New("fetch error"))
	result, err = taskService.GetTaskById(3)
	assert.Error(t, err)
	assert.Equal(t, "fetch error", err.Error())
}

func TestTaskService_GetTasksByTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIRepo(ctrl)
	taskService := TaskService{Repo: mockRepo}
	tasks := []entity.Task{{ID: 1, Name: "Task 1", Tag: "urgent"}}

	// Task successful retrieval
	mockRepo.EXPECT().GetTasksByTag("urgent").Return(tasks, nil)
	result, err := taskService.GetTasksByTag("urgent")
	assert.NoError(t, err)
	assert.Equal(t, tasks, result)

	// Test retrieval error
	mockRepo.EXPECT().GetTasksByTag("nonexistent").Return(nil, errors.New("fetch error"))
	result, err = taskService.GetTasksByTag("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "fetch error", err.Error())
}

func TestTaskService_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIRepo(ctrl)
	taskService := TaskService{Repo: mockRepo}
	task := &entity.Task{ID: 1, Name: "Updated Task"}

	// Task successful update
	mockRepo.EXPECT().UpdateTask(task).Return(nil)
	err := taskService.UpdateTask(task)
	assert.NoError(t, err)

	// Task update error
	mockRepo.EXPECT().UpdateTask(task).Return(errors.New("update error"))
	err = taskService.UpdateTask(task)
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
}

func TestTaskService_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIRepo(ctrl)
	taskService := TaskService{Repo: mockRepo}

	// Task successful deletion
	mockRepo.EXPECT().DeleteTask(1).Return(nil)
	err := taskService.DeleteTask(1)
	assert.NoError(t, err)

	// Task deletion error
	mockRepo.EXPECT().DeleteTask(2).Return(errors.New("deletion error"))
	err = taskService.DeleteTask(2)
	assert.Error(t, err)
	assert.Equal(t, "deletion error", err.Error())
}

func TestTaskService_SearchTasksByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIRepo(ctrl)
	taskService := TaskService{Repo: mockRepo}
	tasks := []entity.Task{{ID: 1, Name: "Test Task"}}

	// Task successful search
	mockRepo.EXPECT().SearchTasksByName("Test").Return(tasks, nil)
	result, err := taskService.SearchTasksByName("Test")
	assert.NoError(t, err)
	assert.Equal(t, tasks, result)

	// Task search error
	mockRepo.EXPECT().SearchTasksByName("Error").Return(nil, errors.New("search error"))
	result, err = taskService.SearchTasksByName("Error")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "search error", err.Error())
}

func TestTaskService_FilterTasksByDeadline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIRepo(ctrl)
	taskService := TaskService{Repo: mockRepo}
	start := time.Now()
	end := start.Add(24 * time.Hour)
	tasks := []entity.Task{{ID: 1, Name: "Task within deadline"}}

	// Task successful filtering
	mockRepo.EXPECT().FilterTasksByDeadline(start, end).Return(tasks, nil)
	result, err := taskService.FilterTasksByDeadline(start, end)
	assert.NoError(t, err)
	assert.Equal(t, tasks, result)

	// Task filtering error
	mockRepo.EXPECT().FilterTasksByDeadline(start, end).Return(nil, errors.New("filtering error"))
	result, err = taskService.FilterTasksByDeadline(start, end)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "filtering error", err.Error())
}
