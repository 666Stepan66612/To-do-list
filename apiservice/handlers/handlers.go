package handlers

import (
	"apiservice/middleware"
	"apiservice/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TaskHandlers struct {
	DBClient      DBClientInterface
	EventProducer EventProducerInterface
}

func NewTaskHandlers(dbClient DBClientInterface, eventProducer EventProducerInterface) *TaskHandlers {
	return &TaskHandlers{
		DBClient:      dbClient,
		EventProducer: eventProducer,
	}
}

func (h *TaskHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, `error: Unauthorized`, http.StatusUnauthorized)
		return
	}

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `error: Invalid JSON`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `error: Name is required`, http.StatusBadRequest)
		return
	}

	task, err := h.DBClient.CreateTask(&req, claims.UserID)
	if err != nil {
		http.Error(w, `error: Failed to create task`, http.StatusInternalServerError)
		h.EventProducer.SendEvent(
			claims.UserID,
			claims.Username,
			"CREATE_TASK",
			fmt.Sprintf("Failed to create task: name=%s", req.Name), "ERROR")
		return
	}

	h.EventProducer.SendEvent(
		claims.UserID,
		claims.Username,
		"CREATE_TASK",
		fmt.Sprintf("Task created: id=%d, name=%s", task.ID, task.Name), "SUCCESS")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, `error: Unauthorized`, http.StatusUnauthorized)
		return
	}

	tasks, err := h.DBClient.GetAllTasks(claims.UserID)
	if err != nil {
		http.Error(w, `error: Failed to get tasks`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, `error: Unauthorized`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	err = h.DBClient.DeleteTask(id, claims.UserID)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete task"}`, http.StatusInternalServerError)
		h.EventProducer.SendEvent(
			claims.UserID,
			claims.Username,
			"DELETE_TASK",
			fmt.Sprintf("Failed to delete task: id=%d", id), "ERROR")
		return
	}

	h.EventProducer.SendEvent(
		claims.UserID,
		claims.Username,
		"DELETE_TASK",
		fmt.Sprintf("Task deleted: id=%d", id), "SUCCESS")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
}

func (h *TaskHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, `error: Unauthorized`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	err = h.DBClient.CompleteTask(id, claims.UserID)
	if err != nil {
		http.Error(w, `{"error": "Failed to complete task"}`, http.StatusInternalServerError)
		h.EventProducer.SendEvent(
			claims.UserID,
			claims.Username,
			"COMPLETE_TASK",
			fmt.Sprintf("Failed to complete task: id=%d", id), "ERROR")
		return
	}

	h.EventProducer.SendEvent(
		claims.UserID,
		claims.Username,
		"COMPLETE_TASK",
		fmt.Sprintf("Task completed: id=%d", id), "SUCCESS")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Task completed successfully",
	})
}

func (h *TaskHandlers) HandleGetCompletedTasks(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, `error: Unauthorized`, http.StatusUnauthorized)
		return
	}

	tasks, err := h.DBClient.GetCompleted(claims.UserID)
	if err != nil {
		http.Error(w, `error: Failed to get tasks`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandlers) HandleGetUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, `error: Unauthorized`, http.StatusUnauthorized)
		return
	}

	tasks, err := h.DBClient.GetUncompleted(claims.UserID)
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
	if err != nil {
		http.Error(w, `error: Invalid task ID`, http.StatusBadRequest)
		return
	}

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
