package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID           int        `json:"id"`
	UserID       int        `json:"user_id"`
	CollectionID *int       `json:"collection_id"`
	Name         string     `json:"name"`
	Text         string     `json:"text"`
	CreateTime   time.Time  `json:"create_time"`
	Complete     bool       `json:"complete"`
	CompleteAt   *time.Time `json:"complete_at"`
}

type Collection struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) CreateTable() error {
	// Создать таблицу users
	_, err := r.DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password_hash VARCHAR(60) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Создать таблицу collections
	_, err = r.DB.Exec(`
	CREATE TABLE IF NOT EXISTS collections (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(100) NOT NULL,
		color VARCHAR(7) DEFAULT '#2564cf',
		icon VARCHAR(50) DEFAULT '📁',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Создать таблицу tasks
	_, err = r.DB.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		collection_id INTEGER REFERENCES collections(id) ON DELETE SET NULL,
		name VARCHAR(255) NOT NULL,
		text TEXT,
		complete BOOLEAN DEFAULT FALSE,
		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		complete_at TIMESTAMP
	)`)
	return err
}

func (r *TaskRepository) CreateTask(task *Task) error {
	return r.DB.QueryRow(`
	INSERT INTO tasks (user_id, collection_id, name, text, complete, create_time) 
	VALUES ($1, $2, $3, $4, FALSE, Now()) 
	RETURNING id, user_id, collection_id, name, text, complete, create_time, complete_at`,
		task.UserID, task.CollectionID, task.Name, task.Text).Scan(
		&task.ID,
		&task.UserID,
		&task.CollectionID,
		&task.Name,
		&task.Text,
		&task.Complete,
		&task.CreateTime,
		&task.CompleteAt)
}

func (r *TaskRepository) GetAllTasksByUser(userID int) ([]Task, error) {
	rows, err := r.DB.Query(`
	SELECT id, user_id, collection_id, name, text, complete, create_time, complete_at FROM tasks
	WHERE user_id = $1
	ORDER BY create_time DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.CollectionID,
			&task.Name,
			&task.Text,
			&task.Complete,
			&task.CreateTime,
			&task.CompleteAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetAllTasks() ([]Task, error) {
	rows, err := r.DB.Query(`
	SELECT * FROM tasks
	ORDER BY createtime DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Text,
			&task.Complete,
			&task.CreateTime,
			&task.CompleteAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetCompletedTasksByUser(userID int) ([]Task, error) {
	rows, err := r.DB.Query(`
	SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks 
	WHERE complete = TRUE AND user_id = $1
	ORDER BY create_time DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Name,
			&task.Text,
			&task.Complete,
			&task.CreateTime,
			&task.CompleteAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetCompletedTasks() ([]Task, error) {
	rows, err := r.DB.Query(`
	SELECT * FROM tasks 
	WHERE complete = TRUE 
	ORDER BY createtime DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Text,
			&task.Complete,
			&task.CreateTime,
			&task.CompleteAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetUncompletedTasksByUser(userID int) ([]Task, error) {
	rows, err := r.DB.Query(`
	SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks
	WHERE complete = FALSE AND user_id = $1
	ORDER BY create_time DESC`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Name,
			&task.Text,
			&task.Complete,
			&task.CreateTime,
			&task.CompleteAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetUncompletedTasks() ([]Task, error) {
	rows, err := r.DB.Query(`
	SELECT * FROM tasks
	WHERE complete = FALSE
	ORDER BY createtime DESC`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Text,
			&task.Complete,
			&task.CreateTime,
			&task.CompleteAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetTaskByID(id int) (*Task, error) {
	var task Task

	err := r.DB.QueryRow(`
	SELECT * FROM tasks 
	WHERE id = $1`,
		id).Scan(&task.ID,
		&task.Name,
		&task.Text,
		&task.Complete,
		&task.CreateTime,
		&task.CompleteAt)

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetIDByName(name string) (*Task, error) {
	var task Task

	err := r.DB.QueryRow(`SELECT * FROM tasks 
	WHERE name = $1`,
		name).Scan(
		&task.ID,
		&task.Name,
		&task.Text,
		&task.Complete,
		&task.CreateTime,
		&task.CompleteAt)

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) DeleteTaskByUser(id, userID int) error {
	result, err := r.DB.Exec(`DELETE FROM tasks WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task not found or access denied")
	}
	return nil
}

func (r *TaskRepository) DeleteTask(id int) error {
	_, err := r.DB.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	return err
}

// Collection methods

func (r *TaskRepository) CreateCollection(collection *Collection) error {
	return r.DB.QueryRow(`
	INSERT INTO collections (user_id, name, color, icon, created_at) 
	VALUES ($1, $2, $3, $4, Now()) 
	RETURNING id, user_id, name, color, icon, created_at`,
		collection.UserID, collection.Name, collection.Color, collection.Icon).Scan(
		&collection.ID,
		&collection.UserID,
		&collection.Name,
		&collection.Color,
		&collection.Icon,
		&collection.CreatedAt)
}

func (r *TaskRepository) GetCollectionsByUser(userID int) ([]Collection, error) {
	rows, err := r.DB.Query(`
	SELECT id, user_id, name, color, icon, created_at FROM collections
	WHERE user_id = $1
	ORDER BY created_at ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var collection Collection
		if err := rows.Scan(
			&collection.ID,
			&collection.UserID,
			&collection.Name,
			&collection.Color,
			&collection.Icon,
			&collection.CreatedAt); err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}
	return collections, rows.Err()
}

func (r *TaskRepository) DeleteCollectionByUser(id, userID int) error {
	result, err := r.DB.Exec(`DELETE FROM collections WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("collection not found or access denied")
	}
	return nil
}

func (r *TaskRepository) GetTasksByCollection(userID, collectionID int) ([]Task, error) {
	rows, err := r.DB.Query(`
	SELECT id, user_id, collection_id, name, text, complete, create_time, complete_at FROM tasks
	WHERE user_id = $1 AND collection_id = $2
	ORDER BY create_time DESC`, userID, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.CollectionID,
			&task.Name,
			&task.Text,
			&task.Complete,
			&task.CreateTime,
			&task.CompleteAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) CompleteTaskByUser(id, userID int) error {
	result, err := r.DB.Exec(`
    UPDATE tasks 
    SET complete = TRUE,
    complete_at = Now()
    WHERE id = $1 AND user_id = $2 AND complete = FALSE`, id, userID)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task already completed, not found, or access denied")
	}

	return nil
}

func (r *TaskRepository) CompleteTask(id int) error {
	result, err := r.DB.Exec(`
    UPDATE tasks 
    SET complete = TRUE,
    completeat = Now()
    WHERE id = $1 AND complete = FALSE`, id)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task already completed or not found")
	}

	return nil
}
