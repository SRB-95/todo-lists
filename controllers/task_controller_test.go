package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo-lists/entity"
	"todo-lists/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	t.Run("Successful creation", func(t *testing.T) {
		tc := TaskController{
			Service: mockService,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Request = &http.Request{
			Method: http.MethodPost,
			Header: map[string][]string{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(bytes.NewBuffer([]byte(`{"name": "test", "deadline": "2024-10-22T17:00:00+05:30", "tag":"high"}`))),
		}
		mockService.EXPECT().CreateTask(gomock.Any()).Return(nil).Times(1)
		tc.CreateTask(ginContext)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("Invalid payload", func(t *testing.T) {
		tc := TaskController{
			Service: mockService,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Request = &http.Request{
			Method: http.MethodPost,
			Header: map[string][]string{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(bytes.NewBuffer([]byte(``))),
		}
		tc.CreateTask(ginContext)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid input"}`, w.Body.String())
	})

	t.Run("Error creating task", func(t *testing.T) {
		tc := TaskController{
			Service: mockService,
		}

		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Request = &http.Request{
			Method: http.MethodPost,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: io.NopCloser(bytes.NewBuffer([]byte(`{"name": "test", "deadline": "2024-10-22T17:00:00+05:30", "tag":"high"}`))),
		}

		mockService.EXPECT().CreateTask(gomock.Any()).Return(errors.New("creation error")).Times(1)

		tc.CreateTask(ginContext)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Error creating task"}`, w.Body.String())
	})
}

func TestGetTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)
	tc := TaskController{
		Service: mockService,
	}

	gin.SetMode(gin.TestMode)

	t.Run("Successful retrieval of tasks", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		expectedTasks := []entity.Task{
			{Name: "Task 1", Deadline: time.Now(), Tag: "high"},
			{Name: "Task 2", Deadline: time.Now().Add(48 * time.Hour), Tag: "medium"},
		}

		mockService.EXPECT().GetAllTasks().Return(expectedTasks, nil).Times(1)

		tc.GetTasks(ginContext)

		assert.Equal(t, http.StatusOK, w.Code)

		// Unmarshal response into a variable for comparison
		var actualTasks []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualTasks)
		assert.NoError(t, err)

		// Compare relevant fields
		for i, task := range actualTasks {
			assert.Equal(t, expectedTasks[i].Name, task["name"])
			assert.Equal(t, expectedTasks[i].Tag, task["tag"])
		}
	})

	t.Run("Error fetching tasks", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		mockService.EXPECT().GetAllTasks().Return(nil, errors.New("fetching error")).Times(1)

		tc.GetTasks(ginContext)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Error fetching tasks"}`, w.Body.String())
	})

	t.Run("Service returns nil slice without error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		mockService.EXPECT().GetAllTasks().Return([]entity.Task{}, nil).Times(1)

		tc.GetTasks(ginContext)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[]`, w.Body.String())
	})

}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	tc := TaskController{
		Service: mockService,
	}

	gin.SetMode(gin.TestMode)

	t.Run("Successful update", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		taskToUpdate := entity.Task{
			ID:       1,
			Name:     "Updated Task",
			Deadline: time.Now(),
			Tag:      "medium",
		}

		// Set up the request body
		body, _ := json.Marshal(taskToUpdate)
		ginContext.Request = &http.Request{
			Method: http.MethodPut,
			Header: map[string][]string{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(bytes.NewBuffer(body)),
		}

		mockService.EXPECT().UpdateTask(gomock.Any()).Return(nil).Times(1)

		tc.UpdateTask(ginContext)

		assert.Equal(t, http.StatusOK, w.Code)

		var updatedTask entity.Task
		err := json.Unmarshal(w.Body.Bytes(), &updatedTask)
		assert.NoError(t, err)
		assert.Equal(t, taskToUpdate.Name, updatedTask.Name)
		assert.Equal(t, taskToUpdate.Tag, updatedTask.Tag)
	})

	t.Run("Invalid ID format", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

		tc.UpdateTask(ginContext)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid task ID"}`, w.Body.String())
	})

	t.Run("Invalid input payload", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		ginContext.Request = &http.Request{
			Method: http.MethodPut,
			Header: map[string][]string{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(bytes.NewBuffer([]byte(`invalid json`))),
		}

		tc.UpdateTask(ginContext)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid input"}`, w.Body.String())
	})

	t.Run("Task not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		taskToUpdate := entity.Task{
			ID:   1,
			Name: "Non-existent Task",
		}

		body, _ := json.Marshal(taskToUpdate)
		ginContext.Request = &http.Request{
			Method: http.MethodPut,
			Header: map[string][]string{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(bytes.NewBuffer(body)),
		}

		mockService.EXPECT().UpdateTask(gomock.Any()).Return(gorm.ErrRecordNotFound).Times(1)

		tc.UpdateTask(ginContext)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error": "Task not found"}`, w.Body.String())
	})

	t.Run("Error updating task", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)
		ginContext.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		taskToUpdate := entity.Task{
			ID:   1,
			Name: "Task with update error",
		}

		body, _ := json.Marshal(taskToUpdate)
		ginContext.Request = &http.Request{
			Method: http.MethodPut,
			Header: map[string][]string{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(bytes.NewBuffer(body)),
		}

		mockService.EXPECT().UpdateTask(gomock.Any()).Return(errors.New("update error")).Times(1)

		tc.UpdateTask(ginContext)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Error updating task"}`, w.Body.String())
	})
}

func TestSearchTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	tc := TaskController{
		Service: mockService,
	}

	gin.SetMode(gin.TestMode)

	t.Run("Successful search", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the query parameter for the search
		ginContext.Request, _ = http.NewRequest(http.MethodGet, "/tasks/search?keyword=test", nil)

		expectedTasks := []entity.Task{
			{Name: "Test Task 1", Deadline: time.Now(), Tag: "high"},
			{Name: "Test Task 2", Deadline: time.Now().Add(48 * time.Hour), Tag: "medium"},
		}

		mockService.EXPECT().SearchTasksByName("test").Return(expectedTasks, nil).Times(1)

		tc.SearchTasks(ginContext)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualTasks []entity.Task
		err := json.Unmarshal(w.Body.Bytes(), &actualTasks)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedTasks), len(actualTasks))
	})

	t.Run("No tasks found", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the query parameter for the search
		ginContext.Request, _ = http.NewRequest(http.MethodGet, "/tasks/search?keyword=test", nil)

		mockService.EXPECT().SearchTasksByName("test").Return([]entity.Task{}, nil).Times(1)

		tc.SearchTasks(ginContext)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"message": "No tasks found matching with name"}`, w.Body.String())
	})

	t.Run("Error searching tasks", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the query parameter for the search
		ginContext.Request, _ = http.NewRequest(http.MethodGet, "/tasks/search?keyword=test", nil)

		mockService.EXPECT().SearchTasksByName("test").Return(nil, errors.New("search error")).Times(1)

		tc.SearchTasks(ginContext)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Error searching tasks"}`, w.Body.String())
	})
}

func TestFilterTasksByDeadline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	tc := TaskController{
		Service: mockService,
	}

	gin.SetMode(gin.TestMode)

	t.Run("Successful filter by deadline", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the query parameters for the date range
		ginContext.Request, _ = http.NewRequest(http.MethodGet, "/tasks/filter?start=2024-10-01&end=2024-10-31", nil)

		expectedTasks := []entity.Task{
			{Name: "Task 1", Deadline: time.Now(), Tag: "high"},
			{Name: "Task 2", Deadline: time.Now().Add(48 * time.Hour), Tag: "medium"},
		}

		mockService.EXPECT().FilterTasksByDeadline(gomock.Any(), gomock.Any()).Return(expectedTasks, nil).Times(1)

		tc.FilterTasksByDeadline(ginContext)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualTasks []entity.Task
		err := json.Unmarshal(w.Body.Bytes(), &actualTasks)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedTasks), len(actualTasks))
	})

	t.Run("Invalid start date format", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the query parameter with an invalid start date
		ginContext.Request, _ = http.NewRequest(http.MethodGet, "/tasks/filter?start=invalid-date&end=2024-10-31", nil)

		tc.FilterTasksByDeadline(ginContext)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid start date format"}`, w.Body.String())
	})

	t.Run("Invalid end date format", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the query parameter with an invalid end date
		ginContext.Request, _ = http.NewRequest(http.MethodGet, "/tasks/filter?start=2024-10-01&end=invalid-date", nil)

		tc.FilterTasksByDeadline(ginContext)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid end date format"}`, w.Body.String())
	})

	t.Run("No tasks found in date range", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the query parameters for the date range
		ginContext.Request, _ = http.NewRequest(http.MethodGet, "/tasks/filter?start=2024-10-01&end=2024-10-31", nil)

		mockService.EXPECT().FilterTasksByDeadline(gomock.Any(), gomock.Any()).Return([]entity.Task{}, nil).Times(1)

		tc.FilterTasksByDeadline(ginContext)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"message": "No tasks found in the specified date range"}`, w.Body.String())
	})

	t.Run("Error filtering tasks", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the query parameters for the date range
		ginContext.Request, _ = http.NewRequest(http.MethodGet, "/tasks/filter?start=2024-10-01&end=2024-10-31", nil)

		mockService.EXPECT().FilterTasksByDeadline(gomock.Any(), gomock.Any()).Return(nil, errors.New("filter error")).Times(1)

		tc.FilterTasksByDeadline(ginContext)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Error filtering tasks"}`, w.Body.String())
	})
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	tc := TaskController{
		Service: mockService,
	}

	gin.SetMode(gin.TestMode)

	t.Run("Successful deletion", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the ID parameter for deletion
		ginContext.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		mockService.EXPECT().DeleteTask(1).Return(nil).Times(1)

		tc.DeleteTask(ginContext)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Task deleted successfully"}`, w.Body.String())
	})

	t.Run("Invalid task ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set an invalid ID parameter for deletion
		ginContext.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		tc.DeleteTask(ginContext)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid task ID"}`, w.Body.String())
	})

	t.Run("Error deleting task", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the ID parameter for deletion
		ginContext.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		mockService.EXPECT().DeleteTask(1).Return(errors.New("deletion error")).Times(1)

		tc.DeleteTask(ginContext)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Error deleting task"}`, w.Body.String())
	})
}

func TestGetTaskByTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	tc := TaskController{
		Service: mockService,
	}

	gin.SetMode(gin.TestMode)

	t.Run("Successful retrieval of tasks by tag", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the tag parameter for retrieval
		ginContext.Params = gin.Params{
			{Key: "tag", Value: "important"},
		}

		tasks := []entity.Task{
			{ID: 1, Name: "Task 1", Tag: "important"},
			{ID: 2, Name: "Task 2", Tag: "important"},
		}

		mockService.EXPECT().GetTasksByTag("important").Return(tasks, nil).Times(1)

		tc.GetTaskByTag(ginContext)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entity.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, uint(1), response[0].ID)
		assert.Equal(t, "Task 1", response[0].Name)
	})

	t.Run("No tasks found for the specified tag", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the tag parameter for retrieval
		ginContext.Params = gin.Params{
			{Key: "tag", Value: "nonexistent"},
		}

		mockService.EXPECT().GetTasksByTag("nonexistent").Return([]entity.Task{}, nil).Times(1)

		tc.GetTaskByTag(ginContext)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error": "No tasks found with the specified tag"}`, w.Body.String())
	})

	t.Run("Error fetching tasks from the service", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Set the tag parameter for retrieval
		ginContext.Params = gin.Params{
			{Key: "tag", Value: "error"},
		}

		mockService.EXPECT().GetTasksByTag("error").Return(nil, errors.New("service error")).Times(1)

		tc.GetTaskByTag(ginContext)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Error fetching tasks"}`, w.Body.String())
	})
}
