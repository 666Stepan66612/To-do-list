package handlers

import (
	"apiservice/client"
	"apiservice/models"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/gorilla/mux"
)

type TaskHandlers struct{
	DBClient *client.DBClient
}

func NewTaskHandlers(dbClient *client.DBClient) *TaskHandlers {
	return &TaskHandlers{
		DBClient: dbClient,
	}
}

func (h *TaskHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `error: Invalid JSON`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `error: Name is required`, http.StatusBadRequest)
		return
	}

	task, err := h.DBClient.CreateTask(&req)
	if err != nil {
		http.Error(w, `error: Failed to create task`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.DBClient.GetAllTasks()
	if err != nil {
		http.Error(w, `error: Failed to get tasks`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	err = h.DBClient.DeleteTask(id)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
}

func (h *TaskHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
        return
    }

    err = h.DBClient.CompleteTask(id)
    if err != nil {
        http.Error(w, `{"error": "Failed to complete task"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Task completed successfully",
    })
}

func (h *TaskHandlers) HandleGetCompletedTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.DBClient.GetCompleted()
	if err != nil {
		http.Error(w, `error: Failed to get tasks`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
	
func (h *TaskHandlers) HandleGetUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.DBClient.GetUncompleted()
	if err != nil {
		http.Error(w, `error: Failed to get tasks`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandlers) HandleGetTasksByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	tasks, err := h.DBClient.GetTaskByID(id)
	if err != nil {
		http.Error(w, `error: Failed to get tasks`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandlers) HandleGetTasksByName(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	tasks, err := h.DBClient.GetTaskByName(name)
	if err != nil {
		http.Error(w, `error: Failed to get tasks`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}