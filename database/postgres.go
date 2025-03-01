package database

import (
	"fmt"
	"log"

	"github.com/drive-deep/task-microservice/config"
	"github.com/drive-deep/task-microservice/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

func (p *PostgresDB) Connect() (*gorm.DB, error) {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
	var err error
	p.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Enable Citus extension
	if err := p.db.Exec("CREATE EXTENSION IF NOT EXISTS citus").Error; err != nil {
		return nil, err
	}

	// Add worker nodes
	workerNodes := []string{"postgres-worker1:5432", "postgres-worker2:5432"}
	for _, worker := range workerNodes {
		if err := p.db.Exec("SELECT * from master_add_node(?)", worker).Error; err != nil {
			log.Printf("Failed to add worker node %s: %v", worker, err)
		}
	}
	// Run migrations on all models
    if err := p.db.AutoMigrate(&models.Task{}); err != nil {
        return nil, fmt.Errorf("failed to run migrations: %w", err)
    }

	// Check if the tasks table is already distributed
	var count int64
	if err := p.db.Raw("SELECT count(*) FROM pg_dist_partition WHERE logicalrelid = 'tasks'::regclass").Scan(&count).Error; err != nil {
		log.Printf("failed to check if tasks table is distributed: %v", err)
	}

	if count == 0 {
		// Distribute the tasks table
		if err := p.db.Exec("SELECT create_distributed_table('tasks', 'id')").Error; err != nil {
			return nil, fmt.Errorf("failed to distribute tasks table: %w", err)
		}
		log.Println("tasks table is now distributed")
	} else {
		log.Println("tasks table is already distributed")
	}

	return p.db, nil
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgresDB) CreateTask(task *models.Task) error {
	return p.db.Create(task).Error
}

func (p *PostgresDB) GetTaskByID(id string) (*models.Task, error) {
	var task models.Task
	err := p.db.First(&task, "id = ?", id).Error
	return &task, err
}

func (p *PostgresDB) GetAllTasksPaginated(page, pageSize int) ([]models.Task, error) {
	var tasks []models.Task
	offset := (page - 1) * pageSize
	err := p.db.Limit(pageSize).Offset(offset).Find(&tasks).Error
	return tasks, err
}

func (p *PostgresDB) UpdateTask(task *models.Task) error {
	return p.db.Save(task).Error
}

func (p *PostgresDB) DeleteTask(id string) error {
	return p.db.Delete(&models.Task{}, "id = ?", id).Error
}
