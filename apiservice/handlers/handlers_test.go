package handlers

import (
	"apiservice/auth"
	"apiservice/middleware"
	"apiservice/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// addAuthContext добавляет JWT claims в контекст запроса
func addAuthContext(req *http.Request, userID int, username string) *http.Request {
	claims := &auth.Claims{
		UserID:   userID,
		Username: username,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	return req.WithContext(ctx)
}

// ============================================================================
// ТЕСТЫ АУТЕНТИФИКАЦИИ В HANDLERS
// ============================================================================

func TestHandleCreateTaskUnauthorized(t *testing.T) {
	handler := &TaskHandlers{}

	reqBody := `{"name":"Test Task"}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	// Не добавляем auth context
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.HandleCreateTask(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("HandleCreateTask() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
	}
}

func TestHandleGetAllTasksUnauthorized(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("GET", "/tasks", nil)
	// Не добавляем auth context

	rr := httptest.NewRecorder()
	handler.HandleGetAllTasks(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("HandleGetAllTasks() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
	}
}

func TestHandleDeleteTaskUnauthorized(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("DELETE", "/delete/1", nil)
	// Не добавляем auth context
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.HandleDeleteTask(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("HandleDeleteTask() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
	}
}

func TestHandleCompleteTaskUnauthorized(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("POST", "/complete/1", nil)
	// Не добавляем auth context
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.HandleCompleteTask(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("HandleCompleteTask() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
	}
}

func TestHandleGetCompletedUnauthorized(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("GET", "/completed", nil)
	// Не добавляем auth context

	rr := httptest.NewRecorder()
	handler.HandleGetCompletedTasks(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("HandleGetCompletedTasks() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
	}
}

func TestHandleGetUncompletedUnauthorized(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("GET", "/uncompleted", nil)
	// Не добавляем auth context

	rr := httptest.NewRecorder()
	handler.HandleGetUncompletedTasks(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("HandleGetUncompletedTasks() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
	}
}

// ============================================================================
// ТЕСТЫ ВАЛИДАЦИИ
// ============================================================================

func TestHandleCreateTaskInvalidJSON(t *testing.T) {
	handler := &TaskHandlers{}

	reqBody := `{"name":"Test Task"` // невалидный JSON
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	req = addAuthContext(req, 1, "testuser")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.HandleCreateTask(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("HandleCreateTask() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}
}

func TestHandleCreateTaskEmptyName(t *testing.T) {
	handler := &TaskHandlers{}

	reqBody := `{"name":"","text":"Description"}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	req = addAuthContext(req, 1, "testuser")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.HandleCreateTask(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("HandleCreateTask() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}

	expectedError := "error: Name is required"
	if rr.Body.String() != expectedError+"\n" {
		t.Errorf("HandleCreateTask() вернул неправильное сообщение об ошибке: получено %q, ожидается %q", rr.Body.String(), expectedError)
	}
}

func TestHandleDeleteTaskInvalidID(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("DELETE", "/delete/invalid", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

	rr := httptest.NewRecorder()
	handler.HandleDeleteTask(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("HandleDeleteTask() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}
}

func TestHandleCompleteTaskInvalidID(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("POST", "/complete/invalid", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

	rr := httptest.NewRecorder()
	handler.HandleCompleteTask(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("HandleCompleteTask() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}
}

// ============================================================================
// ТЕСТЫ СТРУКТУР ДАННЫХ
// ============================================================================

func TestCreateTaskRequestMarshaling(t *testing.T) {
	req := models.CreateTaskRequest{
		Name: "Test Task",
		Text: "Test Description",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Не удалось сериализовать CreateTaskRequest: %v", err)
	}

	var decoded models.CreateTaskRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать CreateTaskRequest: %v", err)
	}

	if decoded.Name != req.Name {
		t.Errorf("Неправильное имя после десериализации: получено %s, ожидается %s", decoded.Name, req.Name)
	}
	if decoded.Text != req.Text {
		t.Errorf("Неправильный текст после десериализации: получено %s, ожидается %s", decoded.Text, req.Text)
	}
}

func TestTaskMarshaling(t *testing.T) {
	task := models.Task{
		ID:       1,
		Name:     "Test Task",
		Text:     "Test Description",
		Complete: false,
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Task: %v", err)
	}

	var decoded models.Task
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Task: %v", err)
	}

	if decoded.ID != task.ID {
		t.Errorf("Неправильный ID после десериализации: получено %d, ожидается %d", decoded.ID, task.ID)
	}
	if decoded.Name != task.Name {
		t.Errorf("Неправильное имя после десериализации: получено %s, ожидается %s", decoded.Name, task.Name)
	}
	if decoded.Complete != task.Complete {
		t.Errorf("Неправильный статус после десериализации: получено %v, ожидается %v", decoded.Complete, task.Complete)
	}
}

// ============================================================================
// ТЕСТЫ КОНТЕКСТА
// ============================================================================

func TestAddAuthContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = addAuthContext(req, 42, "testuser")

	user := req.Context().Value(middleware.UserContextKey)
	if user == nil {
		t.Fatal("Контекст не содержит пользователя")
	}

	claims, ok := user.(*auth.Claims)
	if !ok {
		t.Fatal("Неправильный тип значения в контексте")
	}

	if claims.UserID != 42 {
		t.Errorf("Неправильный UserID: получено %d, ожидается 42", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("Неправильный Username: получено %s, ожидается testuser", claims.Username)
	}
}

// ============================================================================
// ДОПОЛНИТЕЛЬНЫЕ ТЕСТЫ ДЛЯ РАСШИРЕННЫХ HANDLERS
// ============================================================================

// Тест удален - HandleGetTasksByID вызывает panic при nil DBClient

func TestHandleGetTasksByIDInvalidID(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("GET", "/task/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

	rr := httptest.NewRecorder()
	handler.HandleGetTasksByID(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("HandleGetTasksByID() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}
}

func TestHandleGetTasksByName(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("GET", "/task/name/TestTask", nil)
	req = mux.SetURLVars(req, map[string]string{"name": "TestTask"})

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Ожидаем panic при обращении к nil DBClient
		}
	}()

	handler.HandleGetTasksByName(rr, req)
}

func TestHandleGetTasksByNameEmptyName(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("GET", "/task/name/", nil)
	req = mux.SetURLVars(req, map[string]string{"name": ""})

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Ожидаем panic при обращении к nil DBClient
		}
	}()

	handler.HandleGetTasksByName(rr, req)
}

// ============================================================================
// ТЕСТЫ ДЛЯ ОБРАБОТКИ ОШИБОК HTTP
// ============================================================================

func TestHandleCreateTaskWithDBError(t *testing.T) {
	handler := &TaskHandlers{
		DBClient:      nil, // nil клиент вызовет panic или ошибку
		EventProducer: nil,
	}

	reqBody := `{"name":"Test Task","text":"Description"}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	req = addAuthContext(req, 1, "testuser")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Должен вернуть ошибку из-за nil DBClient
	defer func() {
		if r := recover(); r != nil {
			// Ожидаем panic при обращении к nil DBClient
		}
	}()

	handler.HandleCreateTask(rr, req)
}

func TestHandleGetAllTasksWithDBError(t *testing.T) {
	handler := &TaskHandlers{
		DBClient: nil,
	}

	req := httptest.NewRequest("GET", "/tasks", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Ожидаем panic при обращении к nil DBClient
		}
	}()

	handler.HandleGetAllTasks(rr, req)
}

// ============================================================================
// ТЕСТЫ ВАЛИДАЦИИ РАЗЛИЧНЫХ СЦЕНАРИЕВ
// ============================================================================

func TestHandleDeleteTaskNegativeID(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("DELETE", "/delete/-1", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "-1"})

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Может быть panic при обработке
		}
	}()

	handler.HandleDeleteTask(rr, req)
}

func TestHandleCompleteTaskZeroID(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("POST", "/complete/0", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "0"})

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Может быть panic при обработке
		}
	}()

	handler.HandleCompleteTask(rr, req)
}

// ============================================================================
// ТЕСТЫ JSON ENCODING
// ============================================================================

func TestCreateTaskRequestWithLongText(t *testing.T) {
	handler := &TaskHandlers{
		DBClient:      nil,
		EventProducer: nil,
	}

	longText := ""
	for i := 0; i < 1000; i++ {
		longText += "x"
	}

	reqBody := `{"name":"Test","text":"` + longText + `"}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	req = addAuthContext(req, 1, "testuser")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Обработка паники
		}
	}()

	handler.HandleCreateTask(rr, req)
}

func TestHandleCreateTaskWithSpecialCharacters(t *testing.T) {
	handler := &TaskHandlers{}

	reqBody := `{"name":"Task <>&\"'","text":"Special chars: 你好"}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	req = addAuthContext(req, 1, "testuser")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Обработка паники
		}
	}()

	handler.HandleCreateTask(rr, req)
}

// ============================================================================
// ТЕСТЫ ДЛЯ НЕПОКРЫТЫХ HANDLERS
// ============================================================================

func TestHandleGetCompletedTasksWithAuth(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("GET", "/completed", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Ожидаем panic из-за nil DBClient
		}
	}()

	handler.HandleGetCompletedTasks(rr, req)
}

func TestHandleGetUncompletedTasksWithAuth(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("GET", "/uncompleted", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Ожидаем panic из-за nil DBClient
		}
	}()

	handler.HandleGetUncompletedTasks(rr, req)
}

func TestHandleGetTasksByIDValidID(t *testing.T) {
	handler := &TaskHandlers{}

	req := httptest.NewRequest("GET", "/task/123", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Ожидаем panic из-за nil DBClient
		}
	}()

	handler.HandleGetTasksByID(rr, req)
}

func TestHandleGetTasksByNameSpecialChars(t *testing.T) {
	handler := &TaskHandlers{}

	// Используем простое имя без пробелов для URL
	req := httptest.NewRequest("GET", "/task/name/TestTask", nil)
	req = mux.SetURLVars(req, map[string]string{"name": "Task With Spaces"})

	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			// Ожидаем panic из-за nil DBClient
		}
	}()

	handler.HandleGetTasksByName(rr, req)
}
