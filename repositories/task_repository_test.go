package repositories

import (
	"errors"
	"regexp"
	"testing"
	"time"
	"todo-lists/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	// Create a mock SQL connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// Mock GORM's initial query (e.g., SELECT VERSION())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT VERSION()")).
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.0.23"))

	// Open GORM DB with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Return a function to clean up resources
	cleanup := func() {
		// db.Close()
	}
	return gormDB, mock, cleanup
}

func TestCreateTask(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	// Define expected behavior for inserting a task
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `tasks`").WillReturnResult(sqlmock.NewResult(1, 1)) // Mock task creation
	mock.ExpectCommit()

	repo := &TaskRepository{DB: gormDB}
	task := &entity.Task{Name: "Test Task", Deadline: time.Now(), Tag: "high"}

	err := repo.CreateTask(task)
	assert.NoError(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAllTasks(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	// Define the tasks to return on a successful query
	tasks := []entity.Task{
		{ID: 1, Name: "Task 1", Deadline: time.Now(), Tag: "high"},
		{ID: 2, Name: "Task 2", Deadline: time.Now(), Tag: "medium"},
	}

	// Mock the successful retrieval of tasks
	mock.ExpectQuery("SELECT \\* FROM `tasks`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "deadline", "tag"}).
			AddRow(tasks[0].ID, tasks[0].Name, tasks[0].Deadline, tasks[0].Tag).
			AddRow(tasks[1].ID, tasks[1].Name, tasks[1].Deadline, tasks[1].Tag))

	// Create the repository instance
	repo := &TaskRepository{DB: gormDB}

	// Call the GetAllTasks method for successful case
	fetchedTasks, err := repo.GetAllTasks()
	assert.NoError(t, err)
	assert.Equal(t, len(tasks), len(fetchedTasks))
	assert.Equal(t, tasks[0].Name, fetchedTasks[0].Name)
	assert.Equal(t, tasks[1].Name, fetchedTasks[1].Name)

	// Verify expectations were met after successful case
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// Now test the error handling
	mock.ExpectQuery("SELECT \\* FROM `tasks`").WillReturnError(errors.New("db error"))

	// Call the GetAllTasks method again for error case
	fetchedTasks, err = repo.GetAllTasks()
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error()) // Check the specific error message
	assert.Nil(t, fetchedTasks)              // Should return nil tasks on error

	// Verify expectations were met again after error case
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetTaskById(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup() // Ensure cleanup is called

	// Define the expected task details for a successful retrieval
	taskID := uint(2) // Use uint to match the ID type in the Task struct
	taskName := "Task 2"
	taskDeadline := time.Now()
	taskTag := "medium"

	// Mock the retrieval of the task by ID (successful case)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks` WHERE `tasks`.`id` = ? ORDER BY `tasks`.`id` LIMIT ?")).
		WithArgs(taskID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "deadline", "tag"}).
			AddRow(taskID, taskName, taskDeadline, taskTag))

	// Create the repository instance
	repo := &TaskRepository{DB: gormDB}

	// Fetch the task by ID
	fetchedTask, err := repo.GetTaskById(int(taskID)) // Convert to int for the method call
	assert.NoError(t, err)
	assert.Equal(t, taskName, fetchedTask.Name)
	assert.Equal(t, taskID, fetchedTask.ID) // Ensure fetchedTask.ID is compared as uint

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// Test for task not found (error scenario)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks` WHERE `tasks`.`id` = ? ORDER BY `tasks`.`id` LIMIT ?")).
		WithArgs(taskID, 1).
		WillReturnRows(sqlmock.NewRows([]string{})) // No rows returned

	fetchedTask, err = repo.GetTaskById(int(taskID))
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err) // Check that the error is the record not found error
	assert.Equal(t, entity.Task{}, fetchedTask)  // Should return an empty entity.Task

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// Test for other errors (error scenario)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks` WHERE `tasks`.`id` = ? ORDER BY `tasks`.`id` LIMIT ?")).
		WithArgs(taskID, 1).
		WillReturnError(errors.New("db error"))

	fetchedTask, err = repo.GetTaskById(int(taskID))
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())    // Check the specific error message
	assert.Equal(t, entity.Task{}, fetchedTask) // Should return an empty entity.Task

	// Ensure all expectations were met again
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetTasksByTag(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	// Mock the database to return tasks when querying by tag
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks` WHERE tag = ?")).
		WithArgs("high").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "deadline", "tag"}).
			AddRow(1, "Task 1", time.Now(), "high").
			AddRow(2, "Task 2", time.Now(), "high")) // Add more tasks with the same tag as needed

	repo := &TaskRepository{DB: gormDB}

	// Fetch tasks by tag
	tasks, err := repo.GetTasksByTag("high")
	assert.NoError(t, err)

	// Check the number of tasks returned
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, "Task 1", tasks[0].Name)
	assert.Equal(t, "Task 2", tasks[1].Name)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// Now test the error handling
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks` WHERE tag = ?")).
		WithArgs("high").
		WillReturnError(errors.New("db error"))

	// Call the GetTasksByTag method again
	tasks, err = repo.GetTasksByTag("high")
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error()) // Check the specific error message
	assert.Nil(t, tasks)                     // Should return nil tasks on error

	// Verify all expectations were met again
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSearchTasksByName(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	// Define the expected tasks
	expectedTasks := []entity.Task{
		{ID: 1, Name: "Task One", Deadline: time.Now(), Tag: "high"},
		{ID: 2, Name: "Task Two", Deadline: time.Now(), Tag: "high"},
	}

	// Mock the search operation
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks` WHERE name LIKE ?")).
		WithArgs("%One%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "deadline", "tag"}).
			AddRow(expectedTasks[0].ID, expectedTasks[0].Name, expectedTasks[0].Deadline, expectedTasks[0].Tag))

	// Create the repository instance
	repo := &TaskRepository{DB: gormDB}

	// Perform the search
	tasks, err := repo.SearchTasksByName("One")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "Task One", tasks[0].Name)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFilterTasksByDeadline(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	// Define expected tasks based on the deadlines
	start := time.Now().Add(-72 * time.Hour)
	end := time.Now().Add(-12 * time.Hour)

	expectedTasks := []entity.Task{
		{ID: 1, Name: "Task A", Deadline: time.Now().Add(-48 * time.Hour), Tag: "high"},
		{ID: 2, Name: "Task B", Deadline: time.Now().Add(-24 * time.Hour), Tag: "high"},
	}

	// Mock the database query for filtering tasks by deadline
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks` WHERE deadline BETWEEN ? AND ?")).
		WithArgs(start, end).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "deadline", "tag"}).
			AddRow(expectedTasks[0].ID, expectedTasks[0].Name, expectedTasks[0].Deadline, expectedTasks[0].Tag).
			AddRow(expectedTasks[1].ID, expectedTasks[1].Name, expectedTasks[1].Deadline, expectedTasks[1].Tag)) // Return both tasks

	// Create the repository instance
	repo := &TaskRepository{DB: gormDB}

	// Perform the filter operation
	tasks, err := repo.FilterTasksByDeadline(start, end)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, expectedTasks[0].Name, tasks[0].Name)
	assert.Equal(t, expectedTasks[1].Name, tasks[1].Name)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteTask(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &TaskRepository{DB: gormDB}
	taskID := 1

	// Test successful deletion
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `tasks` WHERE `tasks`.`id` = ?")).
		WithArgs(taskID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulate successful delete
	mock.ExpectCommit()

	err := repo.DeleteTask(taskID)
	assert.NoError(t, err)

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test deletion when task is not found
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `tasks` WHERE `tasks`.`id` = ?")).
		WithArgs(taskID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // Simulate delete not found
	mock.ExpectCommit()

	err = repo.DeleteTask(taskID)
	assert.NoError(t, err)

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test error during deletion
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `tasks` WHERE `tasks`.`id` = ?")).
		WithArgs(taskID).
		WillReturnError(errors.New("some database error")) // Simulate an error during delete
	mock.ExpectRollback()

	err = repo.DeleteTask(taskID)
	assert.Error(t, err)
	assert.Equal(t, "some database error", err.Error())

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}
