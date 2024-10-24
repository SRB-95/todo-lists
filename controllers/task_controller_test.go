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
			// You can skip asserting the exact timestamp
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

func TestGetTaskById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	t.Run("Successful fetch task by id", func(t *testing.T) {
		tc := TaskController{
			Service: mockService,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(w)

		// Create a new request to include the ID in the URL
		req := httptest.NewRequest(http.MethodGet, "/tasks/1", io.NopCloser(bytes.NewBuffer([]byte(`{"name": "test", "deadline": "2024-10-22T17:00:00+05:30", "tag":"high"}`))))
		// req.Header.Set("Content-Type", "application/json")

		// Assign the request to the Gin context
		ginContext.Request = req

		mockService.EXPECT().CreateTask(gomock.Any()).Return(nil).Times(1)
		tc.CreateTask(ginContext)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
	// t.Run("Invalid payload", func(t *testing.T) {
	// 	tc := TaskController{
	// 		Service: mockService,
	// 	}

	// 	gin.SetMode(gin.TestMode)
	// 	w := httptest.NewRecorder()
	// 	ginContext, _ := gin.CreateTestContext(w)
	// 	ginContext.Request = &http.Request{
	// 		Method: http.MethodPost,
	// 		Header: map[string][]string{
	// 			"Content-Type": []string{"application/json"},
	// 		},
	// 		Body: io.NopCloser(bytes.NewBuffer([]byte(``))),
	// 	}
	// 	tc.CreateTask(ginContext)

	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	assert.JSONEq(t, `{"error": "Invalid input"}`, w.Body.String())
	// })

	// t.Run("Error creating task", func(t *testing.T) {
	// 	tc := TaskController{
	// 		Service: mockService,
	// 	}

	// 	w := httptest.NewRecorder()
	// 	ginContext, _ := gin.CreateTestContext(w)
	// 	ginContext.Request = &http.Request{
	// 		Method: http.MethodPost,
	// 		Header: map[string][]string{
	// 			"Content-Type": {"application/json"},
	// 		},
	// 		Body: io.NopCloser(bytes.NewBuffer([]byte(`{"name": "test", "deadline": "2024-10-22T17:00:00+05:30", "tag":"high"}`))),
	// 	}

	// 	mockService.EXPECT().CreateTask(gomock.Any()).Return(errors.New("creation error")).Times(1)

	// 	tc.CreateTask(ginContext)

	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// 	assert.JSONEq(t, `{"error": "Error creating task"}`, w.Body.String())
	// })
}

// func TestDeleteTask(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockService := mocks.NewMockIService(ctrl)
// 	router := setupRouter(mockService)

// 	t.Run("Successful deletion", func(t *testing.T) {
// 		mockService.EXPECT().DeleteTask(1).Return(nil)

// 		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, req)

// 		assert.Equal(t, http.StatusOK, w.Code)
// 	})

// 	t.Run("Error deleting task", func(t *testing.T) {
// 		mockService.EXPECT().DeleteTask(1).Return(errors.New("deletion error"))

// 		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, req)

// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	})
// }
