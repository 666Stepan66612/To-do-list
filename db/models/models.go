package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		text TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := r.DB.Exec(query)
	return err
}

func (r *TaskRepository) CreateTask(task *Task) error {
	query := `INSERT INTO tasks (name, text, completed, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	task.CreatedAt = time.Now()
	return r.DB.QueryRow(query, task.Name, task.Text, task.Completed, task.CreatedAt).Scan(&task.ID)
}

func (r *TaskRepository) GetAllTasks() ([]Task, error) {
	query := `SELECT id, name, text, completed, created_at FROM tasks ORDER BY created_at DESC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Text, &task.Completed, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetCompletedTasks() ([]Task, error) {
	query := `SELECT id, name, text, completed, created_at FROM tasks WHERE completed = TRUE ORDER BY created_at DESC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Text, &task.Completed, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetUncompletedTasks() ([]Task, error) {
	query := `SELECT id, name, text, completed, created_at FROM tasks WHERE completed = FALSE ORDER BY created_at DESC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Text, &task.Completed, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetTaskByID(id int) (*Task, error) {
	query := `SELECT id, name, text, completed, created_at FROM tasks WHERE id = $1`
	var task Task
	err := r.DB.QueryRow(query, id).Scan(&task.ID, &task.Name, &task.Text, &task.Completed, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetIDByName(name string) (*Task, error) {
	query := `SELECT id, name, text, completed, created_at FROM tasks WHERE name = $1`
	var task Task
	err := r.DB.QueryRow(query, name).Scan(&task.ID, &task.Name, &task.Text, &task.Completed, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *TaskRepository) CompleteTask(id int) error {
	query := `UPDATE tasks SET completed = TRUE WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
