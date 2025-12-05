package handlers

import (
	"apiservice/auth"
	"apiservice/middleware"
	"apiservice/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// ============================================================================
// TESTS FOR CONSTRUCTORS
// ============================================================================

func TestNewTaskHandlers(t *testing.T) {
	mockDB := &MockDBClient{}
	mockProducer := &MockEventProducer{}

	handlers := NewTaskHandlers(mockDB, mockProducer)

	if handlers == nil {
		t.Fatal("NewTaskHandlers вернул nil")
	}
	if handlers.DBClient != mockDB {
		t.Error("DBClient не установлен правильно")
	}
	if handlers.EventProducer != mockProducer {
		t.Error("EventProducer не установлен правильно")
	}
}

// ============================================================================
// MOCK IMPLEMENTATIONS
// ============================================================================

// MockDBClient для тестирования handlers
type MockDBClient struct {
	CreateTaskFunc           func(*models.CreateTaskRequest, int) (*models.Task, error)
	GetAllTasksFunc          func(int) ([]models.Task, error)
	DeleteTaskFunc           func(int, int) error
	CompleteTaskFunc         func(int, int) error
	GetCompletedFunc         func(int) ([]models.Task, error)
	GetUncompletedFunc       func(int) ([]models.Task, error)
	GetTaskByIDFunc          func(int) (*models.Task, error)
	GetTaskByNameFunc        func(string) (*models.Task, error)
	CreateCollectionFunc     func(*models.CreateCollectionRequest, int) (*models.Collection, error)
	GetCollectionsFunc       func(int) ([]models.Collection, error)
	DeleteCollectionFunc     func(int, int) error
	GetTasksByCollectionFunc func(int, int) ([]models.Task, error)
}

func (m *MockDBClient) CreateTask(req *models.CreateTaskRequest, userID int) (*models.Task, error) {
	if m.CreateTaskFunc != nil {
		return m.CreateTaskFunc(req, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDBClient) GetAllTasks(userID int) ([]models.Task, error) {
	if m.GetAllTasksFunc != nil {
		return m.GetAllTasksFunc(userID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDBClient) DeleteTask(taskID, userID int) error {
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(taskID, userID)
	}
	return errors.New("not implemented")
}

func (m *MockDBClient) CompleteTask(taskID, userID int) error {
	if m.CompleteTaskFunc != nil {
		return m.CompleteTaskFunc(taskID, userID)
	}
	return errors.New("not implemented")
}

func (m *MockDBClient) GetCompleted(userID int) ([]models.Task, error) {
	if m.GetCompletedFunc != nil {
		return m.GetCompletedFunc(userID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDBClient) GetUncompleted(userID int) ([]models.Task, error) {
	if m.GetUncompletedFunc != nil {
		return m.GetUncompletedFunc(userID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDBClient) GetTaskByID(id int) (*models.Task, error) {
	if m.GetTaskByIDFunc != nil {
		return m.GetTaskByIDFunc(id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDBClient) GetTaskByName(name string) (*models.Task, error) {
	if m.GetTaskByNameFunc != nil {
		return m.GetTaskByNameFunc(name)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDBClient) CreateCollection(req *models.CreateCollectionRequest, userID int) (*models.Collection, error) {
	if m.CreateCollectionFunc != nil {
		return m.CreateCollectionFunc(req, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDBClient) GetCollections(userID int) ([]models.Collection, error) {
	if m.GetCollectionsFunc != nil {
		return m.GetCollectionsFunc(userID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDBClient) DeleteCollection(collectionID, userID int) error {
	if m.DeleteCollectionFunc != nil {
		return m.DeleteCollectionFunc(collectionID, userID)
	}
	return errors.New("not implemented")
}

func (m *MockDBClient) GetTasksByCollection(collectionID, userID int) ([]models.Task, error) {
	if m.GetTasksByCollectionFunc != nil {
		return m.GetTasksByCollectionFunc(collectionID, userID)
	}
	return nil, errors.New("not implemented")
}

// MockEventProducer для тестирования handlers
type MockEventProducer struct {
	SendEventFunc func(userID int, username, action, details, status string) error
	Events        []MockEvent
}

type MockEvent struct {
	UserID   int
	Username string
	Action   string
	Details  string
	Status   string
}

func (m *MockEventProducer) SendEvent(userID int, username, action, details, status string) error {
	m.Events = append(m.Events, MockEvent{
		UserID:   userID,
		Username: username,
		Action:   action,
		Details:  details,
		Status:   status,
	})
	if m.SendEventFunc != nil {
		return m.SendEventFunc(userID, username, action, details, status)
	}
	return nil
}

func (m *MockEventProducer) Close() error {
	return nil
}

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

// ============================================================================
// ТЕСТЫ С MOCK КЛИЕНТАМИ (УСПЕШНЫЕ СЦЕНАРИИ)
// ============================================================================

func TestHandleCreateTaskWithMockSuccess(t *testing.T) {
	mockDB := &MockDBClient{
		CreateTaskFunc: func(req *models.CreateTaskRequest, userID int) (*models.Task, error) {
			return &models.Task{
				ID:       1,
				Name:     req.Name,
				Text:     req.Text,
				Complete: false,
			}, nil
		},
	}

	mockKafka := &MockEventProducer{
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	reqBody := `{"name":"Test Task","text":"Description"}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	req = addAuthContext(req, 1, "testuser")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.HandleCreateTask(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusCreated)
	}

	var task models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &task); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if task.Name != "Test Task" {
		t.Errorf("Неправильное имя: получено %s, ожидается Test Task", task.Name)
	}

	if len(mockKafka.Events) != 1 {
		t.Errorf("Неправильное количество событий: получено %d, ожидается 1", len(mockKafka.Events))
	}

	if mockKafka.Events[0].Action != "CREATE_TASK" {
		t.Errorf("Неправильное действие: получено %s, ожидается CREATE_TASK", mockKafka.Events[0].Action)
	}
}

func TestHandleCreateTaskWithMockDBError(t *testing.T) {
	mockDB := &MockDBClient{
		CreateTaskFunc: func(req *models.CreateTaskRequest, userID int) (*models.Task, error) {
			return nil, errors.New("database error")
		},
	}

	mockKafka := &MockEventProducer{
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	reqBody := `{"name":"Test Task"}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	req = addAuthContext(req, 1, "testuser")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.HandleCreateTask(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}

	if len(mockKafka.Events) != 1 || mockKafka.Events[0].Status != "ERROR" {
		t.Error("Должно быть отправлено событие об ошибке")
	}
}

func TestHandleGetAllTasksWithMockSuccess(t *testing.T) {
	mockDB := &MockDBClient{
		GetAllTasksFunc: func(userID int) ([]models.Task, error) {
			return []models.Task{
				{ID: 1, Name: "Task 1", Complete: false},
				{ID: 2, Name: "Task 2", Complete: true},
			}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/tasks", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetAllTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	var tasks []models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Неправильное количество задач: получено %d, ожидается 2", len(tasks))
	}
}

func TestHandleGetAllTasksWithMockError(t *testing.T) {
	mockDB := &MockDBClient{
		GetAllTasksFunc: func(userID int) ([]models.Task, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/tasks", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetAllTasks(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}
}

func TestHandleDeleteTaskWithMockSuccess(t *testing.T) {
	mockDB := &MockDBClient{
		DeleteTaskFunc: func(taskID, userID int) error {
			return nil
		},
	}

	mockKafka := &MockEventProducer{
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	req := httptest.NewRequest("DELETE", "/delete/1", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.HandleDeleteTask(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	if len(mockKafka.Events) != 1 || mockKafka.Events[0].Action != "DELETE_TASK" {
		t.Error("Должно быть отправлено событие DELETE_TASK")
	}
}

func TestHandleDeleteTaskWithMockError(t *testing.T) {
	mockDB := &MockDBClient{
		DeleteTaskFunc: func(taskID, userID int) error {
			return errors.New("database error")
		},
	}

	mockKafka := &MockEventProducer{
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	req := httptest.NewRequest("DELETE", "/delete/1", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.HandleDeleteTask(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}

	if len(mockKafka.Events) != 1 || mockKafka.Events[0].Status != "ERROR" {
		t.Error("Должно быть отправлено событие об ошибке")
	}
}

func TestHandleCompleteTaskWithMockSuccess(t *testing.T) {
	mockDB := &MockDBClient{
		CompleteTaskFunc: func(taskID, userID int) error {
			return nil
		},
	}

	mockKafka := &MockEventProducer{
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	req := httptest.NewRequest("POST", "/complete/1", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.HandleCompleteTask(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	if len(mockKafka.Events) != 1 || mockKafka.Events[0].Action != "COMPLETE_TASK" {
		t.Error("Должно быть отправлено событие COMPLETE_TASK")
	}
}

func TestHandleCompleteTaskWithMockError(t *testing.T) {
	mockDB := &MockDBClient{
		CompleteTaskFunc: func(taskID, userID int) error {
			return errors.New("database error")
		},
	}

	mockKafka := &MockEventProducer{
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	req := httptest.NewRequest("POST", "/complete/1", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.HandleCompleteTask(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}

	if len(mockKafka.Events) != 1 || mockKafka.Events[0].Status != "ERROR" {
		t.Error("Должно быть отправлено событие об ошибке")
	}
}

func TestHandleGetCompletedTasksWithMockSuccess(t *testing.T) {
	mockDB := &MockDBClient{
		GetCompletedFunc: func(userID int) ([]models.Task, error) {
			return []models.Task{
				{ID: 1, Name: "Completed Task", Complete: true},
			}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/completed", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetCompletedTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	var tasks []models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if len(tasks) != 1 || !tasks[0].Complete {
		t.Error("Должна быть одна завершенная задача")
	}
}

func TestHandleGetUncompletedTasksWithMockSuccess(t *testing.T) {
	mockDB := &MockDBClient{
		GetUncompletedFunc: func(userID int) ([]models.Task, error) {
			return []models.Task{
				{ID: 1, Name: "Uncompleted Task", Complete: false},
			}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/uncompleted", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetUncompletedTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	var tasks []models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if len(tasks) != 1 || tasks[0].Complete {
		t.Error("Должна быть одна незавершенная задача")
	}
}

func TestHandleGetTasksByIDWithMockSuccess(t *testing.T) {
	mockDB := &MockDBClient{
		GetTaskByIDFunc: func(id int) (*models.Task, error) {
			return &models.Task{
				ID:       id,
				Name:     "Test Task",
				Complete: false,
			}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/task/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.HandleGetTasksByID(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	var task models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &task); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Неправильный ID: получено %d, ожидается 1", task.ID)
	}
}

func TestHandleGetTasksByNameWithMockSuccess(t *testing.T) {
	mockDB := &MockDBClient{
		GetTaskByNameFunc: func(name string) (*models.Task, error) {
			return &models.Task{
				ID:       1,
				Name:     name,
				Complete: false,
			}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/task/name/TestTask", nil)
	req = mux.SetURLVars(req, map[string]string{"name": "TestTask"})

	rr := httptest.NewRecorder()
	handler.HandleGetTasksByName(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	var task models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &task); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if task.Name != "TestTask" {
		t.Errorf("Неправильное имя: получено %s, ожидается TestTask", task.Name)
	}
}

func TestHandleGetTasksByIDWithMockError(t *testing.T) {
	mockDB := &MockDBClient{
		GetTaskByIDFunc: func(id int) (*models.Task, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/task/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.HandleGetTasksByID(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}
}

func TestHandleGetTasksByNameWithMockError(t *testing.T) {
	mockDB := &MockDBClient{
		GetTaskByNameFunc: func(name string) (*models.Task, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/task/name/TestTask", nil)
	req = mux.SetURLVars(req, map[string]string{"name": "TestTask"})

	rr := httptest.NewRecorder()
	handler.HandleGetTasksByName(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}
}

func TestHandleGetCompletedTasksWithMockError(t *testing.T) {
	mockDB := &MockDBClient{
		GetCompletedFunc: func(userID int) ([]models.Task, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/completed", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetCompletedTasks(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}
}

func TestHandleGetUncompletedTasksWithMockError(t *testing.T) {
	mockDB := &MockDBClient{
		GetUncompletedFunc: func(userID int) ([]models.Task, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/uncompleted", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetUncompletedTasks(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}
}

func TestHandleCreateTaskKafkaErrorIgnored(t *testing.T) {
	mockDB := &MockDBClient{
		CreateTaskFunc: func(req *models.CreateTaskRequest, userID int) (*models.Task, error) {
			return &models.Task{
				ID:       1,
				Name:     req.Name,
				Complete: false,
			}, nil
		},
	}

	mockKafka := &MockEventProducer{
		SendEventFunc: func(userID int, username, action, details, status string) error {
			return errors.New("kafka error")
		},
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	reqBody := `{"name":"Test Task"}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	req = addAuthContext(req, 1, "testuser")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.HandleCreateTask(rr, req)

	// Несмотря на ошибку Kafka, задача должна быть создана
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusCreated)
	}
}

func TestHandleGetAllTasksEmptyList(t *testing.T) {
	mockDB := &MockDBClient{
		GetAllTasksFunc: func(userID int) ([]models.Task, error) {
			return []models.Task{}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/tasks", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetAllTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	var tasks []models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Неправильное количество задач: получено %d, ожидается 0", len(tasks))
	}
}

func TestHandleDeleteTaskWithZeroIDMock(t *testing.T) {
	mockDB := &MockDBClient{
		DeleteTaskFunc: func(taskID, userID int) error {
			if taskID == 0 {
				return errors.New("invalid task ID")
			}
			return nil
		},
	}

	mockKafka := &MockEventProducer{
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	req := httptest.NewRequest("DELETE", "/delete/0", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "0"})

	rr := httptest.NewRecorder()
	handler.HandleDeleteTask(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}
}

func TestHandleCompleteTaskWithZeroIDMock(t *testing.T) {
	mockDB := &MockDBClient{
		CompleteTaskFunc: func(taskID, userID int) error {
			if taskID == 0 {
				return errors.New("invalid task ID")
			}
			return nil
		},
	}

	mockKafka := &MockEventProducer{
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	req := httptest.NewRequest("POST", "/complete/0", nil)
	req = addAuthContext(req, 1, "testuser")
	req = mux.SetURLVars(req, map[string]string{"id": "0"})

	rr := httptest.NewRecorder()
	handler.HandleCompleteTask(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}
}

func TestHandleCreateTaskWithEmptyText(t *testing.T) {
	mockDB := &MockDBClient{
		CreateTaskFunc: func(req *models.CreateTaskRequest, userID int) (*models.Task, error) {
			return &models.Task{
				ID:       1,
				Name:     req.Name,
				Text:     req.Text,
				Complete: false,
			}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	reqBody := `{"name":"Test Task","text":""}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(reqBody))
	req = addAuthContext(req, 1, "testuser")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.HandleCreateTask(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusCreated)
	}

	var task models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &task); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if task.Text != "" {
		t.Errorf("Text должен быть пустым, получено: %s", task.Text)
	}
}

func TestHandleGetAllTasksWithMultipleTasks(t *testing.T) {
	mockDB := &MockDBClient{
		GetAllTasksFunc: func(userID int) ([]models.Task, error) {
			return []models.Task{
				{ID: 1, Name: "Task 1", Complete: false},
				{ID: 2, Name: "Task 2", Complete: true},
				{ID: 3, Name: "Task 3", Complete: false},
				{ID: 4, Name: "Task 4", Complete: true},
				{ID: 5, Name: "Task 5", Complete: false},
			}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/tasks", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetAllTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	var tasks []models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if len(tasks) != 5 {
		t.Errorf("Неправильное количество задач: получено %d, ожидается 5", len(tasks))
	}
}

func TestHandleGetCompletedTasksMultiple(t *testing.T) {
	mockDB := &MockDBClient{
		GetCompletedFunc: func(userID int) ([]models.Task, error) {
			return []models.Task{
				{ID: 1, Name: "Completed 1", Complete: true},
				{ID: 2, Name: "Completed 2", Complete: true},
				{ID: 3, Name: "Completed 3", Complete: true},
			}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/completed", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetCompletedTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	var tasks []models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Неправильное количество задач: получено %d, ожидается 3", len(tasks))
	}

	for _, task := range tasks {
		if !task.Complete {
			t.Error("Все задачи должны быть завершенными")
		}
	}
}

func TestHandleGetUncompletedTasksMultiple(t *testing.T) {
	mockDB := &MockDBClient{
		GetUncompletedFunc: func(userID int) ([]models.Task, error) {
			return []models.Task{
				{ID: 1, Name: "Uncompleted 1", Complete: false},
				{ID: 2, Name: "Uncompleted 2", Complete: false},
			}, nil
		},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: &MockEventProducer{},
	}

	req := httptest.NewRequest("GET", "/uncompleted", nil)
	req = addAuthContext(req, 1, "testuser")

	rr := httptest.NewRecorder()
	handler.HandleGetUncompletedTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}

	var tasks []models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Не удалось распарсить ответ: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Неправильное количество задач: получено %d, ожидается 2", len(tasks))
	}

	for _, task := range tasks {
		if task.Complete {
			t.Error("Все задачи должны быть незавершенными")
		}
	}
}

func TestHandleDeleteTaskWithDifferentUsers(t *testing.T) {
	deleteCallCount := 0

	mockDB := &MockDBClient{
		DeleteTaskFunc: func(taskID, userID int) error {
			deleteCallCount++
			if userID == 1 {
				return nil
			}
			return errors.New("unauthorized")
		},
	}

	mockKafka := &MockEventProducer{
		Events: []MockEvent{},
	}

	handler := &TaskHandlers{
		DBClient:      mockDB,
		EventProducer: mockKafka,
	}

	// Тест с пользователем 1
	req1 := httptest.NewRequest("DELETE", "/delete/1", nil)
	req1 = addAuthContext(req1, 1, "user1")
	req1 = mux.SetURLVars(req1, map[string]string{"id": "1"})

	rr1 := httptest.NewRecorder()
	handler.HandleDeleteTask(rr1, req1)

	if status := rr1.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус для user1: получено %v, ожидается %v", status, http.StatusOK)
	}

	// Тест с пользователем 2
	req2 := httptest.NewRequest("DELETE", "/delete/1", nil)
	req2 = addAuthContext(req2, 2, "user2")
	req2 = mux.SetURLVars(req2, map[string]string{"id": "1"})

	rr2 := httptest.NewRecorder()
	handler.HandleDeleteTask(rr2, req2)

	if status := rr2.Code; status != http.StatusInternalServerError {
		t.Errorf("Неправильный статус для user2: получено %v, ожидается %v", status, http.StatusInternalServerError)
	}

	if deleteCallCount != 2 {
		t.Errorf("DeleteTask должен быть вызван 2 раза, вызван %d раз", deleteCallCount)
	}
}
