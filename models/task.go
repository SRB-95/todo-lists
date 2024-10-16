package models

import (
	"time"
)

// Task represents the task model
type Task struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"not null" json:"name"`
	Deadline time.Time `gorm:"not null" json:"deadline"`
	Tag      string    `gorm:"type:enum('less', 'medium', 'high');not null" json:"tag"`
}
