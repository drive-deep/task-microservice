package cache

import (
    "github.com/drive-deep/task-microservice/models"
)

type Task = models.Task

type Cache interface {
    Connect() (Cache, error)
    Close() error
    AddTask(task Task) error
    GetTask(id string) (Task, error)
    GetPaginatedTasks(page, pageSize int) ([]Task, error)
    UpdateTask(task Task) error
    DeleteTask(id string) error
}