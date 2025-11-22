package models

import (
	"database/sql"
	"time"
	_ "github.com/lib/pq"
)

type TaskRepository struct {
	DB *sql.DB
}

type Task struct {
	ID         int
	Name       string
	Text       string
	CreateTime time.Time
	Complete   bool
	CompleteAt time.Time
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) CreateTable() error {
	_, err := r.DB.Exec(`
		CREATE TABLE IF NOT EXISTS Tasks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		text VARCHAR(100),
		create_time TIMESTAMP,
		complete BOOLEAN DEFAULT FALSE,
		complete_at TIMESTAMP)`)
	return err
}

func (r *TaskRepository) CreateTask(t *Task) error {
	err := r.DB.QueryRow(`
        INSERT INTO tasks (name, text, create_time, complete) 
        VALUES ($1, $2, Now(), FALSE) 
        RETURNING id, create_time
    `, t.Name, t.Text).Scan(&t.ID, &t.CreateTime)

	return err
}

func (r *TaskRepository) GetAllTasks() ([]Task, error) {
	var rows *sql.Rows
	var err error
	rows, err = r.DB.Query(`SELECT * FROM Tasks`)
	defer rows.Close()

	var tasks []Task
	var task Task
	for rows.Next() {
		rows.Scan(
			&task.ID,
			&task.Name,
			&task.Text,
			&task.CreateTime,
			&task.Complete,
			&task.CompleteAt)

		tasks = append(tasks, task)
	}

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetTaskByID(id int) (*Task, error) {
	var rows *sql.Rows
	var err error
	rows, err = r.DB.Query(`SELECT * FROM Tasks WHERE id = $1`, id)
	defer rows.Close()

	var task Task
	for rows.Next() {
		rows.Scan(
			&task.ID,
			&task.Name,
			&task.Text,
			&task.CreateTime,
			&task.Complete,
			&task.CompleteAt)
		
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) GetIDByName(name string) (*Task, error) {
	var rows *sql.Rows
	var err error
	rows, err = r.DB.Query(`SELECT * FROM Tasks WHERE name = $1`, name)
	defer rows.Close()

	var task Task
	for rows.Next() {
		rows.Scan(
			&task.ID,
			&task.Name,
			&task.Text,
			&task.CreateTime,
			&task.Complete,
			&task.CompleteAt)
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) GetCompletedTasks() ([]Task, error) {
	var rows *sql.Rows
	var err error
	rows, err = r.DB.Query(`SELECT * FROM Tasks WHERE complete = TRUE`)
	defer rows.Close()

	var tasks []Task
	var task Task
	for rows.Next() {
		rows.Scan(
			&task.ID,
			&task.Name,
			&task.Text,
			&task.CreateTime,
			&task.Complete,
			&task.CompleteAt)
		
		tasks = append(tasks, task)
	}

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetUncompletedTasks() ([]Task, error) {
	var rows *sql.Rows
	var err error
	rows, err = r.DB.Query(`SELECT * FROM Tasks WHERE complete = FALSE`)
	defer rows.Close()

	var tasks []Task
	var task Task
	for rows.Next() {
		rows.Scan(
			&task.ID,
			&task.Name,
			&task.Text,
			&task.CreateTime,
			&task.Complete,
			&task.CompleteAt)

		tasks = append(tasks, task)
	}

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteTask(id int) error {
	_, err := r.DB.Exec(`
        DELETE FROM tasks WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) CompleteTask(id int) error {
	_, err := r.DB.Exec(`
        UPDATE tasks 
        SET complete = TRUE, complete_at = NOW() 
        WHERE id = $1
    `, id)

	if err != nil {
		return err
	}

	return nil
}

/*
var GetSliceTaskTransfer []Task

func SQLstart(){
	connStr := "user=postgres password=mypostgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
	defer db.Close()

	_, err = db.Exec(`
		CREATE DATABASE IF NOT EXISTS TODO`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Tasks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		text VARCHAR(100),
		time TIME,
		complete BOOLEAN DEFAULT FALSE)
		completeat TIME`)

	if err != nil {
		log.Fatal(err)
	}
}

func PostSQL(name string, text string){
	connStr := "user=postgres password= dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO Tasks(name, text, time, complete)
		VALUES ($1, $2, Now(), 'False')`, name, text)

	if err != nil {
		log.Fatal(err)
	}
}

func GetSQL(GetSliceTaskTransfer *[]Task){
	connStr := "user=postgres password= dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
	defer db.Close()

	*GetSliceTaskTransfer = nil
	var rows *sql.Rows
	rows, err = db.Query(`SELECT * FROM Tasks`)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	*GetSliceTaskTransfer = nil

	for rows.Next() {
		var id int
		var name, text string
		var complete bool
		var taskTime,completeAt time.Time
		err := rows.Scan(&id, &name, &text, &taskTime, &complete, &completeAt)
		if err != nil {
			log.Fatal(err)
		}

		*GetSliceTaskTransfer = append(*GetSliceTaskTransfer, Task{id, name, text, taskTime, complete, completeAt})
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func FilterGetSQL(GetSliceTaskTransfer *[]Task, flag bool){
	connStr := "user=postgres password= dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
	defer db.Close()

	if flag{
		var rows *sql.Rows
		rows, err = db.Query(`SELECT * FROM Tasks WHERE complete = True`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
		var id int
		var name, text string
		var complete bool
		var taskTime, completeat time.Time
		err := rows.Scan(&id, &name, &text, &taskTime, &complete, &completeat)
		if err != nil {
			log.Fatal(err)
		}
		*GetSliceTaskTransfer = append(*GetSliceTaskTransfer, Task{id, name, text, taskTime, complete, completeat})
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
	}else{
		var rows *sql.Rows
		rows, err = db.Query(`SELECT * FROM Tasks WHERE complete = False`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var name, text string
			var complete bool
			var taskTime, completeAt time.Time
			err := rows.Scan(&id, &name, &text, &taskTime, &complete, &completeAt)
			if err != nil {
				log.Fatal(err)
			}
			*GetSliceTaskTransfer = append(*GetSliceTaskTransfer, Task{id, name, text, taskTime, complete, completeAt})
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func DeleteSQL(id int){
	connStr := "user=postgres password= dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM Tasks
		WHERE id = $1`, id)

	if err != nil {
		log.Fatal(err)
	}
}

func CompleteSQL(id int){
	connStr := "user=postgres password=90800022 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil{
		log.Fatal(err)
	}

	_, err = tx.Exec(`UPDATE Tasks SET complete = True
	WHERE id = $1`, id)

	if err != nil {
    	tx.Rollback()
    	log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
    	log.Fatal(err)
	}
}
*/
