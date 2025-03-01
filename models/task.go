package models

import "time"

// Task represents a task with a title, description, status, priority, and timestamps
type Task struct {
    ID          string    `json:"id" gorm:"type:string;primaryKey;index"`
    Title       string    `json:"title" gorm:"type:varchar(100)"`
    Description string    `json:"description" gorm:"type:text"`
    Status      string    `json:"status" gorm:"type:varchar(20)"`
    Priority    int       `json:"priority" gorm:"type:int"`
    CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;default:current_timestamp;autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;default:current_timestamp;autoUpdateTime"`
}