package repositories

import (
    "fmt"
    "github.com/drive-deep/task-microservice/models"

    "gorm.io/gorm"
)

type Task = models.Task

type TaskRepository struct {
    db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
    return &TaskRepository{db}
}

func (r *TaskRepository) Create(entity *Task) error {
    return r.db.Create(entity).Error
}

func (r *TaskRepository) GetByID(id string) (*Task, error) {
    var task Task
    err := r.db.First(&task, "id = ?", id).Error
    return &task, err
}

func (r *TaskRepository) GetAll(filter map[string]interface{}, sort string, page, pageSize int) ([]Task, error) {
    var tasks []Task
    query := r.db.Model(&Task{})

    // Apply filters
    for key, value := range filter {
        query = query.Where(fmt.Sprintf("%s = ?", key), value)
    }

    // Apply sorting
    if sort != "" {
        query = query.Order(sort)
    }

    // Apply pagination
    offset := (page - 1) * pageSize
    err := query.Limit(pageSize).Offset(offset).Find(&tasks).Error
    return tasks, err
}

func (r *TaskRepository) Update(entity *Task) error {
    return r.db.Save(entity).Error
}

func (r *TaskRepository) Delete(id string) error {
    return r.db.Delete(&Task{}, "id = ?", id).Error
}