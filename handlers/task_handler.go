package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	page := 1
	query := r.URL.Query()
	var err error
	pageSize := 10
	if ps, ok := query["pageSize"]; ok {
		pageSize, err = strconv.Atoi(ps[0])
		if err != nil {
			http.Error(w, "Invalid pageSize", http.StatusBadRequest)
			return
		}
	}

	if p, ok := query["page"]; ok {
		page, err = strconv.Atoi(p[0])
		if err != nil {
			http.Error(w, "Invalid page", http.StatusBadRequest)
			return
		}
	}

	sort := query.Get("sort")
	filter := make(map[string]interface{})
	if f := query.Get("filter"); f != "" {
		if err := json.Unmarshal([]byte(f), &filter); err != nil {
			http.Error(w, "Invalid filter format", http.StatusBadRequest)
			return
		}
	}

	tasks, err := h.Service.GetAllTasks(filter, sort, page, pageSize)
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
