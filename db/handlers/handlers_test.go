package handlers

import (
	"bytes"
	"database/sql"
	"dbservice/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func setupMockRepo(t *testing.T) (*models.TaskRepository, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	repo := models.NewTaskRepository(db)
	return repo, mock, db
}

// ============================================================================
// ТЕСТЫ ДЛЯ NewTaskHandlers
// ============================================================================

func TestNewTaskHandlers(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)
	if handlers == nil {
		t.Fatal("NewTaskHandlers вернул nil")
	}
	if handlers.Repo != repo {
		t.Error("Repo не установлен правильно")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ HandleCreate
// ============================================================================

func TestHandleCreateSuccess(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, 1, "Test Task", "Test Description", false, now, nil)

	mock.ExpectQuery(`INSERT INTO tasks`).
		WithArgs(1, "Test Task", "Test Description").
		WillReturnRows(rows)

	body := `{"name":"Test Task","text":"Test Description"}`
	req := httptest.NewRequest("POST", "/create?user_id=1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.HandleCreate(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Ожидался код 201, получен %d", rr.Code)
	}

	var task models.Task
	if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
		t.Errorf("Ошибка декодирования ответа: %v", err)
	}

	if task.Name != "Test Task" {
		t.Errorf("Ожидалось имя 'Test Task', получено %s", task.Name)
	}
}

func TestHandleCreateMissingUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	body := `{"name":"Test Task","text":"Test Description"}`
	req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handlers.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleCreateInvalidUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	body := `{"name":"Test Task","text":"Test Description"}`
	req := httptest.NewRequest("POST", "/create?user_id=invalid", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handlers.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleCreateInvalidJSON(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	body := `{"name":"Test Task","text":` // невалидный JSON
	req := httptest.NewRequest("POST", "/create?user_id=1", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handlers.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleCreateEmptyName(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	body := `{"name":"","text":"Test Description"}`
	req := httptest.NewRequest("POST", "/create?user_id=1", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handlers.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleCreateDBError(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	mock.ExpectQuery(`INSERT INTO tasks`).
		WithArgs(1, "Test Task", "Test Description").
		WillReturnError(sql.ErrConnDone)

	body := `{"name":"Test Task","text":"Test Description"}`
	req := httptest.NewRequest("POST", "/create?user_id=1", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handlers.HandleCreate(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался код 500, получен %d", rr.Code)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ HandleGetAll
// ============================================================================

func TestHandleGetAllSuccess(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, 1, "Task 1", "Description 1", false, now, nil).
		AddRow(2, 1, "Task 2", "Description 2", true, now, &now)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/get?user_id=1", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetAll(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", rr.Code)
	}

	var tasks []models.Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Errorf("Ошибка декодирования ответа: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Ожидалось 2 задачи, получено %d", len(tasks))
	}
}

func TestHandleGetAllMissingUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("GET", "/get", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetAll(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleGetAllInvalidUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("GET", "/get?user_id=invalid", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetAll(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleGetAllDBError(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("GET", "/get?user_id=1", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetAll(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался код 500, получен %d", rr.Code)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ HandleGetCompleted
// ============================================================================

func TestHandleGetCompletedSuccess(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, 1, "Completed Task", "Description", true, now, &now)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/get?user_id=1&complete=true", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetCompleted(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", rr.Code)
	}
}

func TestHandleGetCompletedMissingUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("GET", "/get?complete=true", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetCompleted(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleGetCompletedInvalidUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("GET", "/get?user_id=invalid&complete=true", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetCompleted(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleGetCompletedDBError(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("GET", "/get?user_id=1&complete=true", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetCompleted(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался код 500, получен %d", rr.Code)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ HandleGetUncompleted
// ============================================================================

func TestHandleGetUncompletedSuccess(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, 1, "Uncompleted Task", "Description", false, now, nil)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/get?user_id=1&complete=false", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetUncompleted(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", rr.Code)
	}
}

func TestHandleGetUncompletedMissingUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("GET", "/get?complete=false", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetUncompleted(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleGetUncompletedInvalidUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("GET", "/get?user_id=invalid&complete=false", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetUncompleted(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleGetUncompletedDBError(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("GET", "/get?user_id=1&complete=false", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetUncompleted(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался код 500, получен %d", rr.Code)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ HandleDelete
// ============================================================================

func TestHandleDeleteSuccess(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	mock.ExpectExec(`DELETE FROM tasks WHERE id`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/delete/1?user_id=1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handlers.HandleDelete(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", rr.Code)
	}
}

func TestHandleDeleteMissingUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("DELETE", "/delete/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handlers.HandleDelete(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleDeleteInvalidUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("DELETE", "/delete/1?user_id=invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handlers.HandleDelete(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleDeleteInvalidTaskID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("DELETE", "/delete/invalid?user_id=1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rr := httptest.NewRecorder()

	handlers.HandleDelete(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleDeleteDBError(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	mock.ExpectExec(`DELETE FROM tasks WHERE id`).
		WithArgs(1, 1).
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("DELETE", "/delete/1?user_id=1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handlers.HandleDelete(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Ожидался код 403, получен %d", rr.Code)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ HandleComplete
// ============================================================================

func TestHandleCompleteSuccess(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	mock.ExpectExec(`UPDATE tasks`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("PUT", "/complete/1?user_id=1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handlers.HandleComplete(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", rr.Code)
	}
}

func TestHandleCompleteMissingUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("PUT", "/complete/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handlers.HandleComplete(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleCompleteInvalidUserID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("PUT", "/complete/1?user_id=invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handlers.HandleComplete(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleCompleteInvalidTaskID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("PUT", "/complete/invalid?user_id=1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rr := httptest.NewRecorder()

	handlers.HandleComplete(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestHandleCompleteDBError(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	mock.ExpectExec(`UPDATE tasks`).
		WithArgs(1, 1).
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("PUT", "/complete/1?user_id=1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handlers.HandleComplete(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Ожидался код 403, получен %d", rr.Code)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ HandleGetByID
// ============================================================================

func TestHandleGetByIDSuccess(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, "Task 1", "Description", false, now, nil)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/getbyid/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handlers.HandleGetByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", rr.Code)
	}
}

func TestHandleGetByIDInvalidID(t *testing.T) {
	repo, _, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	req := httptest.NewRequest("GET", "/getbyid/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rr := httptest.NewRecorder()

	handlers.HandleGetByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ HandleGetByName
// ============================================================================

func TestHandleGetByNameSuccess(t *testing.T) {
	repo, mock, db := setupMockRepo(t)
	defer db.Close()

	handlers := NewTaskHandlers(repo)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, "Task 1", "Description", false, now, nil)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WithArgs("Task 1").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/getbyname/Task%201", nil)
	req = mux.SetURLVars(req, map[string]string{"name": "Task 1"})
	rr := httptest.NewRecorder()

	handlers.HandleGetByName(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", rr.Code)
	}
}
