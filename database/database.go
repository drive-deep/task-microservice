package database

import (
	"github.com/drive-deep/task-microservice/models"
	"gorm.io/gorm"
)

type Database interface {
	Connect() (*gorm.DB, error)
	Close() error
	CreateTask(task *models.Task) error
	GetTaskByID(id string) (*models.Task, error)
	GetAllTasksPaginated(page, pageSize int) ([]models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id string) error
}
