package main

import (
	"apiservice/client"
	"apiservice/handlersForDB"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main(){
	dbClient := client.NewDBClient("http://db-service:8080")

	taskHandlers := handlers.NewTaskHandlers(dbClient)

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