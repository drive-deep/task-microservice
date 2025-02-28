package routes

import (
	"github.com/drive-deep/task-microservice/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	taskHandler := handlers.NewTaskHandler()

	// Define the routes
	router.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	router.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")
}
