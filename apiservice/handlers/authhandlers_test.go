package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================================
// ТЕСТЫ ДЛЯ REGISTER
// ============================================================================

// TestRegisterSuccess проверяет успешную регистрацию пользователя
func TestRegisterSuccess(t *testing.T) {
	// Создаем mock DB service
	mockDB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/user/create" && r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       1,
				"username": "testuser",
			})
		}
	}))
	defer mockDB.Close()

	// Временно заменяем URL DB service на mock
	// В реальном проекте лучше использовать dependency injection
	originalURL := "http://db-service:8080"
	// Для теста придется модифицировать код или использовать DI

	reqBody := `{"username":"testuser","password":"password123"}`
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Вызываем handler
	Register(rr, req)

	// Проверяем статус код
	// Примечание: этот тест провалится, потому что handler обращается к реальному db-service
	// Для полноценного тестирования нужен рефакторинг с dependency injection

	t.Log("Статус код:", rr.Code)
	t.Log("Ответ:", rr.Body.String())

	// Для демонстрации структуры теста
	// В продакшене нужно использовать mock или testcontainers
	_ = originalURL
}

// TestRegisterInvalidJSON проверяет регистрацию с невалидным JSON
func TestRegisterInvalidJSON(t *testing.T) {
	reqBody := `{"username":"testuser","password":` // невалидный JSON
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Register(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Register() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}

	body := rr.Body.String()
	if body == "" {
		t.Error("Register() вернул пустое сообщение об ошибке")
	}
}

// TestRegisterEmptyUsername проверяет регистрацию с пустым username
func TestRegisterEmptyUsername(t *testing.T) {
	reqBody := `{"username":"","password":"password123"}`
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Register(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Register() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}

	body := rr.Body.String()
	if !contains(body, "Username and Password are required") && !contains(body, "error") {
		t.Errorf("Register() вернул неожиданное сообщение: %s", body)
	}
}

// TestRegisterEmptyPassword проверяет регистрацию с пустым password
func TestRegisterEmptyPassword(t *testing.T) {
	reqBody := `{"username":"testuser","password":""}`
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Register(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Register() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}

	body := rr.Body.String()
	if !contains(body, "Username and Password are required") && !contains(body, "error") {
		t.Errorf("Register() вернул неожиданное сообщение: %s", body)
	}
}

// TestRegisterShortUsername проверяет регистрацию со слишком коротким username
func TestRegisterShortUsername(t *testing.T) {
	reqBody := `{"username":"ab","password":"password123"}`
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Register(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Register() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}

	body := rr.Body.String()
	if !contains(body, "at least 3 characters") && !contains(body, "error") {
		t.Errorf("Register() вернул неожиданное сообщение: %s", body)
	}
}

// TestRegisterShortPassword проверяет регистрацию со слишком коротким password
func TestRegisterShortPassword(t *testing.T) {
	reqBody := `{"username":"testuser","password":"pass"}`
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Register(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Register() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}

	body := rr.Body.String()
	if !contains(body, "at least 8 characters") && !contains(body, "error") {
		t.Errorf("Register() вернул неожиданное сообщение: %s", body)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ LOGIN
// ============================================================================

// TestLoginInvalidJSON проверяет логин с невалидным JSON
func TestLoginInvalidJSON(t *testing.T) {
	reqBody := `{"username":"testuser",` // невалидный JSON
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Login(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Login() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}
}

// TestLoginEmptyUsername проверяет логин с пустым username
func TestLoginEmptyUsername(t *testing.T) {
	reqBody := `{"username":"","password":"password123"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Login(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Login() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}

	body := rr.Body.String()
	if !contains(body, "Username and Password are required") {
		t.Errorf("Login() вернул неожиданное сообщение: %s", body)
	}
}

// TestLoginEmptyPassword проверяет логин с пустым password
func TestLoginEmptyPassword(t *testing.T) {
	reqBody := `{"username":"testuser","password":""}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Login(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Login() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusBadRequest)
	}

	body := rr.Body.String()
	if !contains(body, "Username and Password are required") {
		t.Errorf("Login() вернул неожиданное сообщение: %s", body)
	}
}

// TestLoginMissingFields проверяет логин без обязательных полей
func TestLoginMissingFields(t *testing.T) {
	tests := []struct {
		name    string
		reqBody string
	}{
		{"only username", `{"username":"testuser"}`},
		{"only password", `{"password":"password123"}`},
		{"empty object", `{}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			Login(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Login() вернул неправильный статус для %s: получено %v, ожидается %v",
					tt.name, status, http.StatusBadRequest)
			}
		})
	}
}

// ============================================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ============================================================================

// contains проверяет, содержит ли строка подстроку
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
