package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drive-deep/task-microservice/handlers"
	"github.com/drive-deep/task-microservice/models"
	"github.com/drive-deep/task-microservice/services"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	service := services.NewMockTaskService()
	handler := handlers.NewTaskHandler(service)

	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
		Priority:    1,
		Assignee:    "john.doe",
	}

	body, _ := json.Marshal(task)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.CreateTask(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var createdTask models.Task
	json.NewDecoder(rr.Body).Decode(&createdTask)
	assert.Equal(t, task.Title, createdTask.Title)
}

func TestGetTask(t *testing.T) {
	service := services.NewMockTaskService()
	handler := handlers.NewTaskHandler(service)

	task := models.Task{
		ID:          "1",
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
		Priority:    1,
		Assignee:    "john.doe",
	}
	service.CreateTask(&task)

	req, err := http.NewRequest("GET", "/tasks/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", handler.GetTask).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var fetchedTask models.Task
	json.NewDecoder(rr.Body).Decode(&fetchedTask)
	assert.Equal(t, task.Title, fetchedTask.Title)
}

func TestGetAllTasks(t *testing.T) {
	service := services.NewMockTaskService()
	handler := handlers.NewTaskHandler(service)

	task1 := models.Task{
		Title:       "Test Task 1",
		Description: "Test Description 1",
		Status:      "Pending",
		Priority:    1,
		Assignee:    "john.doe",
	}
	task2 := models.Task{
		Title:       "Test Task 2",
		Description: "Test Description 2",
		Status:      "Completed",
		Priority:    2,
		Assignee:    "jane.doe",
	}
	service.CreateTask(&task1)
	service.CreateTask(&task2)

	req, err := http.NewRequest("GET", "/tasks", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetAllTasks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var tasks []models.Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	assert.Len(t, tasks, 2)
}

func TestUpdateTask(t *testing.T) {
	service := services.NewMockTaskService()
	handler := handlers.NewTaskHandler(service)

	task := models.Task{
		ID:          "1",
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
		Priority:    1,
		Assignee:    "john.doe",
	}
	service.CreateTask(&task)

	updatedTask := models.Task{
		ID:          "1",
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      "In Progress",
		Priority:    2,
		Assignee:    "jane.doe",
	}

	body, _ := json.Marshal(updatedTask)
	req, err := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.UpdateTask(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var fetchedTask models.Task
	json.NewDecoder(rr.Body).Decode(&fetchedTask)
	assert.Equal(t, updatedTask.Title, fetchedTask.Title)
}

func TestDeleteTask(t *testing.T) {
	service := services.NewMockTaskService()
	handler := handlers.NewTaskHandler(service)

	task := models.Task{
		ID:          "1",
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
		Priority:    1,
		Assignee:    "john.doe",
	}
	service.CreateTask(&task)

	req, err := http.NewRequest("DELETE", "/tasks/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", handler.DeleteTask).Methods("DELETE")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}
