package handlers

import (
	"dbservice/models"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type TaskHandlers struct{
	Repo *models.TaskRepository
}

func NewTaskHandlers(repo *models.TaskRepository) *TaskHandlers{
	return &TaskHandlers{Repo: repo}
}

func (h *TaskHandlers)HandleCreate(w http.ResponseWriter, r *http.Request){
	var task struct{
		Name string `json:name`
		Text string `json:text`
	}

	if err := json.NewDecoder(r. Body).Decode(&task); err != nil{
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
        return
	}

	if task.Name == "" {
		http.Error(w, `{"error": "Task name is required"}`, http.StatusBadRequest)
        return
	}

	taskToCreate := &models.Task{
		Name: task.Name,
		Text: task.Text,
	}

	if err := h.Repo.CreateTask(taskToCreate); err != nil{
		http.Error(w, `{"error": "Failed to create task"}`, http.StatusInternalServerError)
        return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(taskToCreate)
}

func (h *TaskHandlers)HandleGetAll(w http.ResponseWriter, r *http.Request){
	tasks, err := h.Repo.GetAllTasks()

	if err != nil{
		http.Error(w, `{"error": "Failed to get tasks"}`, http.StatusInternalServerError)
        return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandlers)HandleGetCompleted(w http.ResponseWriter, r *http.Request){
	tasks, err := h.Repo.GetCompletedTasks()

	if err != nil{
		http.Error(w, `{"error": "Failed to get tasks"}`, http.StatusInternalServerError)
        return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandlers)HandleGetUncompleted(w http.ResponseWriter, r *http.Request){
	tasks, err := h.Repo.GetUncompletedTasks()

	if err != nil{
		http.Error(w, `{"error": "Failed to get tasks"}`, http.StatusInternalServerError)
        return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandlers)HandleGetByID(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)["id"]

	id, err := strconv.Atoi(vars)
	if err != nil{
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
        return
	}

	task, _ := h.Repo.GetTaskByID(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandlers)HandleGetByName(w http.ResponseWriter, r *http.Request){
	name := mux.Vars(r)["name"]

	task, _ := h.Repo.GetIDByName(name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandlers)HandleDelete(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)["id"]

	id, err := strconv.Atoi(vars)
	if err != nil{
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
        return
	}

	h.Repo.DeleteTask(id)

	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandlers)HandleComplete(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)["id"]

	id, err := strconv.Atoi(vars)
	if err != nil{
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
        return
	}

	h.Repo.CompleteTask(id)

	w.WriteHeader(http.StatusOK)
}

/*
func HandleCreate(w http.ResponseWriter, r *http.Request) {
	models.SQLstart()
	var taskDTO dto.TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := dto.ErrorDTO{
			Message: err.Error(),
			Time: time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}


	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := dto.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	models.PostSQL(taskDTO.Name, taskDTO.Text)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	models.SQLstart()
	models.GetSQL()
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	models.SQLstart()
	models.DeleteSQL()
}

func HandleComplete(w http.ResponseWriter, r *http.Request) {
	models.SQLstart()
	models.CompleteSQL()
}
*/