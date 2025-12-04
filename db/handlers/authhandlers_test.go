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
// ТЕСТЫ ДЛЯ CreateUser
// ============================================================================

func TestCreateUserSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := CreateUser(db)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	now := time.Now()
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs("testuser", "hash123").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "created_at"}).
			AddRow(1, "testuser", now))

	body := `{"username":"testuser","password_hash":"hash123"}`
	req := httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Ожидался код 201, получен %d", rr.Code)
	}

	var user models.User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Errorf("Ошибка декодирования ответа: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Ожидалось username 'testuser', получено %s", user.Username)
	}
}

func TestCreateUserInvalidJSON(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := CreateUser(db)

	body := `{"username":"testuser","password_hash":` // невалидный JSON
	req := httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestCreateUserEmptyUsername(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := CreateUser(db)

	body := `{"username":"","password_hash":"hash123"}`
	req := httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestCreateUserEmptyPasswordHash(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := CreateUser(db)

	body := `{"username":"testuser","password_hash":""}`
	req := httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestCreateUserAlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := CreateUser(db)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	body := `{"username":"testuser","password_hash":"hash123"}`
	req := httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("Ожидался код 409, получен %d", rr.Code)
	}
}

func TestCreateUserExistsCheckError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := CreateUser(db)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("testuser").
		WillReturnError(sql.ErrConnDone)

	body := `{"username":"testuser","password_hash":"hash123"}`
	req := httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался код 500, получен %d", rr.Code)
	}
}

func TestCreateUserInsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := CreateUser(db)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs("testuser", "hash123").
		WillReturnError(sql.ErrConnDone)

	body := `{"username":"testuser","password_hash":"hash123"}`
	req := httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался код 500, получен %d", rr.Code)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetUserByUsername
// ============================================================================

func TestGetUserByUsernameSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := GetUserByUsername(db)

	now := time.Now()
	mock.ExpectQuery(`SELECT id, username, password_hash, created_at`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at"}).
			AddRow(1, "testuser", "hash123", now))

	req := httptest.NewRequest("GET", "/user/testuser", nil)
	req = mux.SetURLVars(req, map[string]string{"username": "testuser"})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", rr.Code)
	}

	var user models.User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Errorf("Ошибка декодирования ответа: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Ожидалось username 'testuser', получено %s", user.Username)
	}
}

func TestGetUserByUsernameEmpty(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := GetUserByUsername(db)

	req := httptest.NewRequest("GET", "/user/", nil)
	req = mux.SetURLVars(req, map[string]string{"username": ""})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestGetUserByUsernameNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := GetUserByUsername(db)

	mock.ExpectQuery(`SELECT id, username, password_hash, created_at`).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest("GET", "/user/nonexistent", nil)
	req = mux.SetURLVars(req, map[string]string{"username": "nonexistent"})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Ожидался код 404, получен %d", rr.Code)
	}
}

func TestGetUserByUsernameDBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := GetUserByUsername(db)

	mock.ExpectQuery(`SELECT id, username, password_hash, created_at`).
		WithArgs("testuser").
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("GET", "/user/testuser", nil)
	req = mux.SetURLVars(req, map[string]string{"username": "testuser"})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался код 500, получен %d", rr.Code)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetUserByID
// ============================================================================

func TestGetUserByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := GetUserByID(db)

	now := time.Now()
	mock.ExpectQuery(`SELECT id, username, password_hash, created_at`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at"}).
			AddRow(1, "testuser", "hash123", now))

	req := httptest.NewRequest("GET", "/user/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", rr.Code)
	}
}

func TestGetUserByIDEmpty(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := GetUserByID(db)

	req := httptest.NewRequest("GET", "/user/", nil)
	req = mux.SetURLVars(req, map[string]string{"id": ""})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получен %d", rr.Code)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := GetUserByID(db)

	mock.ExpectQuery(`SELECT id, username, password_hash, created_at`).
		WithArgs("999").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest("GET", "/user/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Ожидался код 404, получен %d", rr.Code)
	}
}

func TestGetUserByIDDBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	handler := GetUserByID(db)

	mock.ExpectQuery(`SELECT id, username, password_hash, created_at`).
		WithArgs("1").
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("GET", "/user/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался код 500, получен %d", rr.Code)
	}
}
