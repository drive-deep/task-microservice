package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/drive-deep/task-microservice/config"
	"github.com/drive-deep/task-microservice/models"
	"github.com/drive-deep/task-microservice/services"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	Service services.TaskService
}

func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{Service: service}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var err error
	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.CreateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.Service.GetTaskByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	query := r.URL.Query()

	// Pagination parameters
	page := cfg.Server.Page
	pageSize := cfg.Server.PageSize

	if p := query.Get("page"); p != "" {
		page, err = strconv.Atoi(p)
		if err != nil {
			http.Error(w, "Invalid page", http.StatusBadRequest)
			return
		}
	}

	if ps := query.Get("page_size"); ps != "" {
		pageSize, err = strconv.Atoi(ps)
		if err != nil {
			http.Error(w, "Invalid page_size", http.StatusBadRequest)
			return
		}
	}

	// Sorting parameters
	sortBy := query.Get("sort_by")
	if sortBy == "" {
		sortBy = "updated_at asc"
	}
	order := query.Get("order")
	if sortBy != "" && order != "" {
		if order != "asc" && order != "desc" {
			http.Error(w, "Invalid order", http.StatusBadRequest)
			return
		}
		sortBy = sortBy + " " + order
	}

	// Filtering parameters
	filter := make(map[string]interface{})
	if status := query.Get("status"); status != "" {
		filter["status"] = status
	}
	if priority := query.Get("priority"); priority != "" {
		filter["priority"] = priority
	}

	tasks, err := h.Service.GetAllTasks(filter, sortBy, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
