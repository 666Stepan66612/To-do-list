package main

import (
	"apiservice/client"
	"apiservice/handlers"
	"apiservice/kafka"
	"apiservice/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	dbClient := client.NewDBClient("http://db-service:8080")

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:29092"
	}

	eventProducer, err := kafka.NewEventProducer([]string{kafkaBrokers}, "task-events")
	if err != nil {
		log.Printf("Warning: Failed to initialize Kafka producer: %v. Events will not be logged.", err)

		eventProducer = nil
	}
	defer func() {
		if eventProducer != nil {
			eventProducer.Close()
		}
	}()

	taskHandlers := handlers.NewTaskHandlers(dbClient, eventProducer)

	//Без JWT
	router := mux.NewRouter()

	router.Use(corsMiddleware)

	router.HandleFunc("/register", handlers.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")

	//C JWT
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Enable CORS
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	protected.Path("/create").Methods("POST", "OPTIONS").HandlerFunc(taskHandlers.HandleCreateTask)
	protected.Path("/get").Methods("GET", "OPTIONS").Queries("complete", "true").HandlerFunc(taskHandlers.HandleGetCompletedTasks)
	protected.Path("/get").Methods("GET", "OPTIONS").Queries("complete", "false").HandlerFunc(taskHandlers.HandleGetUncompletedTasks)
	protected.Path("/get").Methods("GET", "OPTIONS").HandlerFunc(taskHandlers.HandleGetAllTasks)
	protected.Path("/tasks").Methods("GET", "OPTIONS").HandlerFunc(taskHandlers.HandleGetAllTasks)
	protected.Path("/delete/{id}").Methods("DELETE", "OPTIONS").HandlerFunc(taskHandlers.HandleDeleteTask)
	protected.Path("/complete/{id}").Methods("PUT", "POST", "OPTIONS").HandlerFunc(taskHandlers.HandleCompleteTask)
	protected.Path("/getbyid/{id}").Methods("GET", "OPTIONS").HandlerFunc(taskHandlers.HandleGetTasksByID)
	protected.Path("/getbyname/{name}").Methods("GET", "OPTIONS").HandlerFunc(taskHandlers.HandleGetTasksByName)

	// Collection routes
	protected.Path("/collections").Methods("POST", "OPTIONS").HandlerFunc(taskHandlers.HandleCreateCollection)
	protected.Path("/collections").Methods("GET", "OPTIONS").HandlerFunc(taskHandlers.HandleGetCollections)
	protected.Path("/collections/{id}").Methods("DELETE", "OPTIONS").HandlerFunc(taskHandlers.HandleDeleteCollection)
	protected.Path("/collections/{id}/tasks").Methods("GET", "OPTIONS").HandlerFunc(taskHandlers.HandleGetTasksByCollection)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API Service is healthy"))
	}).Methods("GET")

	log.Println("API Service starting on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal(err)
	}

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
