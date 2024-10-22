package entity

import "time"

type Task struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Deadline time.Time `json:"deadline"`
	Tag      string    `json:"tag"`
}
