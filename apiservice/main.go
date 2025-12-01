package main

import (
	"apiservice/client"
	handlers "apiservice/handlersForDB"
	"apiservice/kafka"
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

	router := mux.NewRouter()

	router.Path("/create").Methods("POST").HandlerFunc(taskHandlers.HandleCreateTask)
	router.Path("/get").Methods("GET").Queries("complete", "true").HandlerFunc(taskHandlers.HandleGetCompletedTasks)
	router.Path("/get").Methods("GET").Queries("complete", "false").HandlerFunc(taskHandlers.HandleGetUncompletedTasks)
	router.Path("/get").Methods("GET").HandlerFunc(taskHandlers.HandleGetAllTasks)
	router.Path("/delete/{id}").Methods("DELETE").HandlerFunc(taskHandlers.HandleDeleteTask)
	router.Path("/complete/{id}").Methods("PUT").HandlerFunc(taskHandlers.HandleCompleteTask)
	router.Path("/getbyid/{id}").Methods("GET").HandlerFunc(taskHandlers.HandleGetTasksByID)
	router.Path("/getbyname/{name}").Methods("GET").HandlerFunc(taskHandlers.HandleGetTasksByName)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API Service is healthy"))
	}).Methods("GET")

	log.Println("API Service starting on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal(err)
	}

}
