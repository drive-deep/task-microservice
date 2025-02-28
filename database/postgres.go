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

	// Add worker nodes to the coordinator
	workerNodes := []string{"postgres-worker1:5432", "postgres-worker2:5432"} // Replace with your actual worker node addresses
	for _, node := range workerNodes {
		if err := p.db.Exec(fmt.Sprintf("SELECT * from master_add_node('%s')", node)).Error; err != nil {
			log.Printf("Failed to add worker node %s: %v", node, err)
		}
	}

	// Auto migrate all models in the models package
	if err := p.db.AutoMigrate(
		&models.Task{},
		// Add other models here
	); err != nil {
		return nil, err
	}

	// Distribute the Task table
	if err := p.db.Exec("SELECT create_distributed_table('tasks', 'id')").Error; err != nil {
		return nil, err
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
