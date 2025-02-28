package cache

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"

	"github.com/drive-deep/task-microservice/config"
	"github.com/drive-deep/task-microservice/models"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client  *redis.Client
	ctx     context.Context
	maxSize int
	lruList *list.List
	lruMap  map[string]*list.Element
}

type lruEntry struct {
	key   string
	value models.Task
}

func NewRedisCache(maxSize int) *RedisCache {
	return &RedisCache{
		ctx:     context.Background(),
		maxSize: maxSize,
		lruList: list.New(),
		lruMap:  make(map[string]*list.Element),
	}
}

func (r *RedisCache) Connect() (Cache, error) {
	cfg := config.GetConfig()
	r.client = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) AddTask(task models.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	if err := r.client.Set(r.ctx, task.ID, data, 0).Err(); err != nil {
		return err
	}

	// Add to sorted set for efficient pagination and sorting
	if err := r.client.ZAdd(r.ctx, "tasks", &redis.Z{
		Score:  float64(task.CreatedAt.Unix()), // Use CreatedAt as the score for sorting
		Member: task.ID,
	}).Err(); err != nil {
		return err
	}

	// Add to status set for filtering
	if err := r.client.SAdd(r.ctx, fmt.Sprintf("tasks:status:%s", task.Status), task.ID).Err(); err != nil {
		return err
	}

	// Update LRU cache
	if elem, exists := r.lruMap[task.ID]; exists {
		r.lruList.MoveToFront(elem)
		elem.Value.(*lruEntry).value = task
	} else {
		if r.lruList.Len() >= r.maxSize {
			// Evict the least recently used item
			evictElem := r.lruList.Back()
			if evictElem != nil {
				r.lruList.Remove(evictElem)
				evictEntry := evictElem.Value.(*lruEntry)
				delete(r.lruMap, evictEntry.key)
				r.client.Del(r.ctx, evictEntry.key)
				r.client.ZRem(r.ctx, "tasks", evictEntry.key)
				r.client.SRem(r.ctx, fmt.Sprintf("tasks:status:%s", evictEntry.value.Status), evictEntry.key)
			}
		}
		entry := &lruEntry{key: task.ID, value: task}
		elem := r.lruList.PushFront(entry)
		r.lruMap[task.ID] = elem
	}

	return nil
}

func (r *RedisCache) GetTask(id string) (models.Task, error) {
	val, err := r.client.Get(r.ctx, id).Result()
	if err != nil {
		return models.Task{}, err
	}
	var task models.Task
	if err := json.Unmarshal([]byte(val), &task); err != nil {
		return models.Task{}, err
	}

	// Update LRU cache
	if elem, exists := r.lruMap[id]; exists {
		r.lruList.MoveToFront(elem)
	}

	return task, nil
}

func (r *RedisCache) GetPaginatedTasks(page, pageSize int) ([]models.Task, error) {
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	ids, err := r.client.ZRange(r.ctx, "tasks", int64(start), int64(end)).Result()
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	for _, id := range ids {
		val, err := r.client.Get(r.ctx, id).Result()
		if err != nil {
			return nil, err
		}
		var task models.Task
		if err := json.Unmarshal([]byte(val), &task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *RedisCache) UpdateTask(task models.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	if err := r.client.Set(r.ctx, task.ID, data, 0).Err(); err != nil {
		return err
	}

	// Update sorted set for efficient pagination and sorting
	if err := r.client.ZAdd(r.ctx, "tasks", &redis.Z{
		Score:  float64(task.CreatedAt.Unix()), // Use CreatedAt as the score for sorting
		Member: task.ID,
	}).Err(); err != nil {
		return err
	}

	// Update status set for filtering
	if err := r.client.SAdd(r.ctx, fmt.Sprintf("tasks:status:%s", task.Status), task.ID).Err(); err != nil {
		return err
	}

	// Update LRU cache
	if elem, exists := r.lruMap[task.ID]; exists {
		r.lruList.MoveToFront(elem)
		elem.Value.(*lruEntry).value = task
	} else {
		if r.lruList.Len() >= r.maxSize {
			// Evict the least recently used item
			evictElem := r.lruList.Back()
			if evictElem != nil {
				r.lruList.Remove(evictElem)
				evictEntry := evictElem.Value.(*lruEntry)
				delete(r.lruMap, evictEntry.key)
				r.client.Del(r.ctx, evictEntry.key)
				r.client.ZRem(r.ctx, "tasks", evictEntry.key)
				r.client.SRem(r.ctx, fmt.Sprintf("tasks:status:%s", evictEntry.value.Status), evictEntry.key)
			}
		}
		entry := &lruEntry{key: task.ID, value: task}
		elem := r.lruList.PushFront(entry)
		r.lruMap[task.ID] = elem
	}

	return nil
}

func (r *RedisCache) DeleteTask(id string) error {
	task, err := r.GetTask(id)
	if err != nil {
		return err
	}

	if err := r.client.Del(r.ctx, id).Err(); err != nil {
		return err
	}

	// Remove from sorted set
	if err := r.client.ZRem(r.ctx, "tasks", id).Err(); err != nil {
		return err
	}

	// Remove from status set
	if err := r.client.SRem(r.ctx, fmt.Sprintf("tasks:status:%s", task.Status), id).Err(); err != nil {
		return err
	}

	// Update LRU cache
	if elem, exists := r.lruMap[id]; exists {
		r.lruList.Remove(elem)
		delete(r.lruMap, id)
	}

	return nil
}
