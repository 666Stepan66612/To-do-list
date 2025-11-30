package models

import (
	"database/sql"
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

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) CreateTable() error {
	_, err := r.DB.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		text TEXT,
		complete BOOLEAN DEFAULT FALSE,
		createtime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		completeat TIMESTAMP
	)`)
	return err
}

func (r *TaskRepository) CreateTask(task *Task) error {
	return r.DB.QueryRow(`
	INSERT INTO tasks (name, text, complete, createtime) 
	VALUES ($1, $2, FALSE, Now()) 
	RETURNING id, name, text, complete, createtime, completeat`,
		task.Name, task.Text).Scan(
		&task.ID,
		&task.Name,
		&task.Text,
		&task.Complete,
		&task.CreateTime,
		&task.CompleteAt)
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

func (r *TaskRepository) DeleteTask(id int) error {
	_, err := r.DB.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	return err
}

func (r *TaskRepository) CompleteTask(id int) error {
	_, err := r.DB.Exec(`
	UPDATE tasks 
	SET complete = TRUE,
	completeat = Now()
	WHERE id = $1`, id)
	return err
}
