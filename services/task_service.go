package services

import (
	"github.com/drive-deep/task-microservice/cache"
	"github.com/drive-deep/task-microservice/models"
	"github.com/drive-deep/task-microservice/repositories"
)

type Task = models.Task

type TaskService struct {
	repo  repositories.Repository[Task]
	cache cache.Cache
}

func NewTaskService(repo repositories.Repository[Task], cache cache.Cache) *TaskService {
	return &TaskService{repo, cache}
}

func (s *TaskService) CreateTask(entity *Task) error {
	if err := s.repo.Create(entity); err != nil {
		return err
	}
	if err := s.cache.AddTask(*entity); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) GetTaskByID(id string) (*Task, error) {
	task, err := s.cache.GetTask(id)
	if err == nil {
		return &task, nil
	}
	return s.repo.GetByID(id)
}

func (s *TaskService) GetAllTasks(filter map[string]interface{}, sort string, page, pageSize int) ([]Task, error) {
	// Try to get cached tasks
	if len(filter) == 0 && sort == "" {
		if tasks, err := s.cache.GetPaginatedTasks(page, pageSize); err == nil && len(tasks) == pageSize {
			return tasks, nil
		}
	}

	// If not cached, get from repository
	tasks, err := s.repo.GetAll(filter, sort, page, pageSize)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(entity *Task) error {
	if err := s.repo.Update(entity); err != nil {
		return err
	}
	if err := s.cache.UpdateTask(*entity); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) DeleteTask(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	if err := s.cache.DeleteTask(id); err != nil {
		return err
	}
	return nil
}
