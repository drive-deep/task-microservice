package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/drive-deep/task-microservice/models"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
    tasks map[string]models.Task
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler() *TaskHandler {
    return &TaskHandler{
        tasks: make(map[string]models.Task),
    }
}

// CreateTask handles the creation of a new task
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var task models.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    h.tasks[task.ID] = task
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

// GetTask handles retrieving a task by ID
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    task, exists := h.tasks[id]
    if !exists {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(task)
}

// GetAllTasks handles retrieving all tasks
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    tasks := []models.Task{}
    for _, task := range h.tasks {
        tasks = append(tasks, task)
    }
    json.NewEncoder(w).Encode(tasks)
}

// UpdateTask handles updating an existing task
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    var task models.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if _, exists := h.tasks[id]; !exists {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    h.tasks[id] = task
    json.NewEncoder(w).Encode(task)
}

// DeleteTask handles deleting a task by ID
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    if _, exists := h.tasks[id]; !exists {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    delete(h.tasks, id)
    w.WriteHeader(http.StatusNoContent)
}