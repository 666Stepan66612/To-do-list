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

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
}