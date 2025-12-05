package main

import (
	"database/sql"
	"dbservice/handlers"
	"dbservice/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Параметры из docker-compose.yaml
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "mypostgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "postgres"
	}

	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)

	log.Printf("Connecting to database: host=%s, user=%s, dbname=%s", dbHost, dbUser, dbName)

	var db *sql.DB
	var err error

	for i := 0; i < 30; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to database")
				break
			}
		}
		log.Printf("Failed to connect to database, retrying... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	repo := models.NewTaskRepository(db)
	taskHandlers := handlers.NewTaskHandlers(repo)

	router := mux.NewRouter()

	router.HandleFunc("/user/create", handlers.CreateUser(db)).Methods("POST")
	router.HandleFunc("/user/{username}", handlers.GetUserByUsername(db)).Methods("GET")

	router.Path("/create").Methods("POST").HandlerFunc(taskHandlers.HandleCreate)
	router.Path("/get").Methods("GET").Queries("complete", "true").HandlerFunc(taskHandlers.HandleGetCompleted)
	router.Path("/get").Methods("GET").Queries("complete", "false").HandlerFunc(taskHandlers.HandleGetUncompleted)
	router.Path("/get").Methods("GET").HandlerFunc(taskHandlers.HandleGetAll)
	router.Path("/delete/{id}").Methods("DELETE").HandlerFunc(taskHandlers.HandleDelete)
	router.Path("/complete/{id}").Methods("PUT").HandlerFunc(taskHandlers.HandleComplete)
	router.Path("/getbyid/{id}").Methods("GET").HandlerFunc(taskHandlers.HandleGetByID)
	router.Path("/getbyname/{name}").Methods("GET").HandlerFunc(taskHandlers.HandleGetByName)

	// Collection routes
	router.Path("/collections").Methods("POST").HandlerFunc(taskHandlers.HandleCreateCollection)
	router.Path("/collections").Methods("GET").HandlerFunc(taskHandlers.HandleGetCollections)
	router.Path("/collections/{id}").Methods("DELETE").HandlerFunc(taskHandlers.HandleDeleteCollection)
	router.Path("/collections/{id}/tasks").Methods("GET").HandlerFunc(taskHandlers.HandleGetTasksByCollection)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func runMigrations(db *sql.DB) error {
	//Создаём таблицу users (сначала, т.к. tasks ссылается на неё)
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	//Создаём таблицу tasks
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
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

	//Добавляем колонку user_id (если её нет)
	_, err = db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'tasks' AND column_name = 'user_id'
			) THEN
				ALTER TABLE tasks ADD COLUMN user_id INT REFERENCES users(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("failed to add user_id column: %w", err)
	}

	//Создаём таблицу collections
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS collections (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			color VARCHAR(7) DEFAULT '#2564cf',
			icon VARCHAR(50) DEFAULT '📁',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create collections table: %w", err)
	}

	//Добавляем колонку collection_id (если её нет)
	_, err = db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'tasks' AND column_name = 'collection_id'
			) THEN
				ALTER TABLE tasks ADD COLUMN collection_id INT REFERENCES collections(id) ON DELETE SET NULL;
			END IF;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("failed to add collection_id column: %w", err)
	}

	//Создаём индекс для быстрого поиска задач по user_id
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	log.Println("Database migrations ran successfully")
	return nil
}
