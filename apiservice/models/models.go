package models

import (
	"time"
)

type Task struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Text       string     `json:"text"`
	CreateTime time.Time  `json:"create_time"`
	Complete   bool       `json:"complete"`
	CompleteAt *time.Time `json:"complete_at"`
}

type CreateTaskRequest struct {
	Name string `json:"name"`
	Text string `json:"text"`
}
