package main

import (
	"database/sql"
	"dbservice/handlers"
	"dbservice/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main(){
	connStr := "postgres://postgres:mypostgres@postgres:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
	defer db.Close()

	repo := models.NewTaskRepository(db)
	if err := repo.CreateTable(); err != nil{
		log.Fatal(err)
	}

	taskHandlers := handlers.NewTaskHandlers(repo)

	router := mux.NewRouter()

	router.Path("/create").Methods("POST").HandlerFunc(taskHandlers.HandleCreate)
	router.Path("/get").Methods("GET").Queries("complete", "true").HandlerFunc(taskHandlers.HandleGetCompleted)
	router.Path("/get").Methods("GET").Queries("complete", "false").HandlerFunc(taskHandlers.HandleGetUncompleted)
	router.Path("/get").Methods("GET").HandlerFunc(taskHandlers.HandleGetAll)
	router.Path("/delete/{id}").Methods("DELETE").HandlerFunc(taskHandlers.HandleDelete)
	router.Path("/complete/{id}").Methods("PUT").HandlerFunc(taskHandlers.HandleComplete)
	router.Path("/getbyid/{id}").Methods("GET").HandlerFunc(taskHandlers.HandleGetByID)
	router.Path("/getbyname/{name}").Methods("GET").HandlerFunc(taskHandlers.HandleGetByName)

    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal(err)
    }
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			text TEXT,
			create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			complete BOOLEAN DEFAULT FALSE,
			complete_at TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(60) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)

	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Database migrations ran successfully")
	return nil
}