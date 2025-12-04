package client

import (
	"apiservice/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ============================================================================
// ТЕСТЫ ДЛЯ NewDBClient
// ============================================================================

func TestNewDBClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewDBClient(baseURL)

	if client == nil {
		t.Fatal("NewDBClient() вернул nil")
	}

	if client.BaseURL != baseURL {
		t.Errorf("Неправильный BaseURL: получено %s, ожидается %s", client.BaseURL, baseURL)
	}

	if client.Client == nil {
		t.Error("HTTP Client не должен быть nil")
	}
}

func TestNewDBClientEmptyURL(t *testing.T) {
	client := NewDBClient("")

	if client == nil {
		t.Fatal("NewDBClient() вернул nil")
	}

	if client.BaseURL != "" {
		t.Errorf("BaseURL должен быть пустым, получено %s", client.BaseURL)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ CreateTask
// ============================================================================

func TestCreateTaskSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Неправильный метод: получено %s, ожидается POST", r.Method)
		}

		var req models.CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Не удалось декодировать запрос: %v", err)
		}

		response := models.Task{
			ID:         1,
			Name:       req.Name,
			Text:       req.Text,
			CreateTime: time.Now(),
			Complete:   false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	req := &models.CreateTaskRequest{
		Name: "Test Task",
		Text: "Test Description",
	}

	task, err := client.CreateTask(req, 1)
	if err != nil {
		t.Fatalf("CreateTask() вернул ошибку: %v", err)
	}

	if task.Name != req.Name {
		t.Errorf("Неправильное имя задачи: получено %s, ожидается %s", task.Name, req.Name)
	}
}

func TestCreateTaskInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	req := &models.CreateTaskRequest{Name: "Test"}

	_, err := client.CreateTask(req, 1)
	if err == nil {
		t.Error("CreateTask() должен вернуть ошибку при невалидном JSON")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetAllTasks
// ============================================================================

func TestGetAllTasksSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Неправильный метод: получено %s, ожидается GET", r.Method)
		}

		tasks := []models.Task{
			{ID: 1, Name: "Task 1", Complete: false},
			{ID: 2, Name: "Task 2", Complete: true},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	tasks, err := client.GetAllTasks(1)
	if err != nil {
		t.Fatalf("GetAllTasks() вернул ошибку: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Неправильное количество задач: получено %d, ожидается 2", len(tasks))
	}
}

func TestGetAllTasksEmptyList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Task{})
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	tasks, err := client.GetAllTasks(1)
	if err != nil {
		t.Fatalf("GetAllTasks() вернул ошибку: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Должен вернуть пустой список, получено %d задач", len(tasks))
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetCompleted
// ============================================================================

func TestGetCompletedSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		complete := r.URL.Query().Get("complete")
		if complete != "true" {
			t.Errorf("Неправильный параметр complete: получено %s, ожидается true", complete)
		}

		tasks := []models.Task{
			{ID: 1, Name: "Completed Task", Complete: true},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	tasks, err := client.GetCompleted(1)
	if err != nil {
		t.Fatalf("GetCompleted() вернул ошибку: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Неправильное количество задач: получено %d, ожидается 1", len(tasks))
	}

	if !tasks[0].Complete {
		t.Error("Задача должна быть завершенной")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetUncompleted
// ============================================================================

func TestGetUncompletedSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		complete := r.URL.Query().Get("complete")
		if complete != "false" {
			t.Errorf("Неправильный параметр complete: получено %s, ожидается false", complete)
		}

		tasks := []models.Task{
			{ID: 1, Name: "Uncompleted Task", Complete: false},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	tasks, err := client.GetUncompleted(1)
	if err != nil {
		t.Fatalf("GetUncompleted() вернул ошибку: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Неправильное количество задач: получено %d, ожидается 1", len(tasks))
	}

	if tasks[0].Complete {
		t.Error("Задача не должна быть завершенной")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ DeleteTask
// ============================================================================

func TestDeleteTaskSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Неправильный метод: получено %s, ожидается DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	err := client.DeleteTask(1, 1)
	if err != nil {
		t.Fatalf("DeleteTask() вернул ошибку: %v", err)
	}
}

func TestDeleteTaskServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	err := client.DeleteTask(1, 1)
	// DeleteTask не проверяет статус ответа, поэтому не возвращает ошибку
	if err != nil {
		t.Errorf("DeleteTask() вернул неожиданную ошибку: %v", err)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ CompleteTask
// ============================================================================

func TestCompleteTaskSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Неправильный метод: получено %s, ожидается PUT", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	err := client.CompleteTask(1, 1)
	if err != nil {
		t.Fatalf("CompleteTask() вернул ошибку: %v", err)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetTaskByID
// ============================================================================

func TestGetTaskByIDSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := models.Task{
			ID:       1,
			Name:     "Test Task",
			Complete: false,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	task, err := client.GetTaskByID(1)
	if err != nil {
		t.Fatalf("GetTaskByID() вернул ошибку: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Неправильный ID задачи: получено %d, ожидается 1", task.ID)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetTaskByName
// ============================================================================

func TestGetTaskByNameSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := models.Task{
			ID:       1,
			Name:     "Test Task",
			Complete: false,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	task, err := client.GetTaskByName("Test Task")
	if err != nil {
		t.Fatalf("GetTaskByName() вернул ошибку: %v", err)
	}

	if task.Name != "Test Task" {
		t.Errorf("Неправильное имя задачи: получено %s, ожидается Test Task", task.Name)
	}
}

// ============================================================================
// ТЕСТЫ ОБРАБОТКИ ОШИБОК СЕТИ
// ============================================================================

func TestCreateTaskNetworkError(t *testing.T) {
	client := NewDBClient("http://invalid-host-that-does-not-exist:9999")
	req := &models.CreateTaskRequest{Name: "Test"}

	_, err := client.CreateTask(req, 1)
	if err == nil {
		t.Error("CreateTask() должен вернуть ошибку при сетевой ошибке")
	}
}

func TestGetAllTasksNetworkError(t *testing.T) {
	client := NewDBClient("http://invalid-host:9999")

	_, err := client.GetAllTasks(1)
	if err == nil {
		t.Error("GetAllTasks() должен вернуть ошибку при сетевой ошибке")
	}
}

// ============================================================================
// ТЕСТЫ ПАРАМЕТРОВ URL
// ============================================================================

func TestCreateTaskWithUserID(t *testing.T) {
	userIDReceived := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "42" {
			userIDReceived = 42
		}

		response := models.Task{ID: 1, Name: "Test"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	req := &models.CreateTaskRequest{Name: "Test"}

	_, err := client.CreateTask(req, 42)
	if err != nil {
		t.Fatalf("CreateTask() вернул ошибку: %v", err)
	}

	if userIDReceived != 42 {
		t.Errorf("Неправильный user_id: получено %d, ожидается 42", userIDReceived)
	}
}

// ============================================================================
// ДОПОЛНИТЕЛЬНЫЕ ТЕСТЫ ДЛЯ УВЕЛИЧЕНИЯ ПОКРЫТИЯ
// ============================================================================

func TestDeleteTaskWithZeroID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	err := client.DeleteTask(0, 1)
	if err != nil {
		t.Fatalf("DeleteTask() вернул ошибку: %v", err)
	}
}

func TestCompleteTaskWithNegativeID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	err := client.CompleteTask(-1, 1)
	if err != nil {
		t.Fatalf("CompleteTask() вернул ошибку: %v", err)
	}
}

func TestGetTaskByIDNetworkError(t *testing.T) {
	client := NewDBClient("http://invalid-host:9999")

	_, err := client.GetTaskByID(1)
	if err == nil {
		t.Error("GetTaskByID() должен вернуть ошибку при сетевой ошибке")
	}
}

func TestGetTaskByNameNetworkError(t *testing.T) {
	client := NewDBClient("http://invalid-host:9999")

	_, err := client.GetTaskByName("test")
	if err == nil {
		t.Error("GetTaskByName() должен вернуть ошибку при сетевой ошибке")
	}
}

func TestGetCompletedNetworkError(t *testing.T) {
	client := NewDBClient("http://invalid-host:9999")

	_, err := client.GetCompleted(1)
	if err == nil {
		t.Error("GetCompleted() должен вернуть ошибку при сетевой ошибке")
	}
}

func TestGetUncompletedNetworkError(t *testing.T) {
	client := NewDBClient("http://invalid-host:9999")

	_, err := client.GetUncompleted(1)
	if err == nil {
		t.Error("GetUncompleted() должен вернуть ошибку при сетевой ошибке")
	}
}

func TestDeleteTaskNetworkError(t *testing.T) {
	client := NewDBClient("http://invalid-host:9999")

	err := client.DeleteTask(1, 1)
	if err == nil {
		t.Error("DeleteTask() должен вернуть ошибку при сетевой ошибке")
	}
}

func TestCompleteTaskNetworkError(t *testing.T) {
	client := NewDBClient("http://invalid-host:9999")

	err := client.CompleteTask(1, 1)
	if err == nil {
		t.Error("CompleteTask() должен вернуть ошибку при сетевой ошибке")
	}
}

func TestCreateTaskInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("{invalid json}"))
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	req := &models.CreateTaskRequest{Name: "Test"}

	_, err := client.CreateTask(req, 1)
	if err == nil {
		t.Error("CreateTask() должен вернуть ошибку при невалидном JSON ответе")
	}
}

func TestGetAllTasksInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	_, err := client.GetAllTasks(1)
	if err == nil {
		t.Error("GetAllTasks() должен вернуть ошибку при невалидном JSON")
	}
}

func TestGetCompletedInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[{invalid}]"))
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	_, err := client.GetCompleted(1)
	if err == nil {
		t.Error("GetCompleted() должен вернуть ошибку при невалидном JSON")
	}
}

func TestGetUncompletedInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bad json"))
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	_, err := client.GetUncompleted(1)
	if err == nil {
		t.Error("GetUncompleted() должен вернуть ошибку при невалидном JSON")
	}
}

func TestGetTaskByIDInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{not valid json"))
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	_, err := client.GetTaskByID(1)
	if err == nil {
		t.Error("GetTaskByID() должен вернуть ошибку при невалидном JSON")
	}
}

func TestGetTaskByNameInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("}}}}"))
	}))
	defer server.Close()

	client := NewDBClient(server.URL)
	_, err := client.GetTaskByName("test")
	if err == nil {
		t.Error("GetTaskByName() должен вернуть ошибку при невалидном JSON")
	}
}
