package models

import (
	"encoding/json"
	"testing"
	"time"
)

// ============================================================================
// ТЕСТЫ ДЛЯ Task
// ============================================================================

func TestTaskMarshaling(t *testing.T) {
	now := time.Now()
	completeAt := now.Add(time.Hour)

	task := Task{
		ID:         1,
		Name:       "Test Task",
		Text:       "Test Description",
		CreateTime: now,
		Complete:   true,
		CompleteAt: &completeAt,
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Task: %v", err)
	}

	var decoded Task
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Task: %v", err)
	}

	if decoded.ID != task.ID {
		t.Errorf("Неправильный ID: получено %d, ожидается %d", decoded.ID, task.ID)
	}
	if decoded.Name != task.Name {
		t.Errorf("Неправильный Name: получено %s, ожидается %s", decoded.Name, task.Name)
	}
	if decoded.Complete != task.Complete {
		t.Errorf("Неправильный Complete: получено %v, ожидается %v", decoded.Complete, task.Complete)
	}
}

func TestTaskWithNilCompleteAt(t *testing.T) {
	task := Task{
		ID:         1,
		Name:       "Incomplete Task",
		Text:       "Not completed",
		CreateTime: time.Now(),
		Complete:   false,
		CompleteAt: nil,
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Task с nil CompleteAt: %v", err)
	}

	var decoded Task
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Task с nil CompleteAt: %v", err)
	}

	if decoded.CompleteAt != nil {
		t.Errorf("CompleteAt должен быть nil для незавершенной задачи, получено %v", decoded.CompleteAt)
	}
}

func TestTaskJSONFields(t *testing.T) {
	task := Task{
		ID:         1,
		Name:       "Task",
		Text:       "Description",
		CreateTime: time.Now(),
		Complete:   false,
		CompleteAt: nil,
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Task: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Не удалось десериализовать Task в map: %v", err)
	}

	expectedFields := []string{"id", "name", "text", "create_time", "complete", "complete_at"}
	for _, field := range expectedFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("Отсутствует поле %s в JSON", field)
		}
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ CreateTaskRequest
// ============================================================================

func TestCreateTaskRequestMarshaling(t *testing.T) {
	req := CreateTaskRequest{
		Name: "New Task",
		Text: "Task Description",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Не удалось сериализовать CreateTaskRequest: %v", err)
	}

	var decoded CreateTaskRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать CreateTaskRequest: %v", err)
	}

	if decoded.Name != req.Name {
		t.Errorf("Неправильный Name: получено %s, ожидается %s", decoded.Name, req.Name)
	}
	if decoded.Text != req.Text {
		t.Errorf("Неправильный Text: получено %s, ожидается %s", decoded.Text, req.Text)
	}
}

func TestCreateTaskRequestEmptyFields(t *testing.T) {
	req := CreateTaskRequest{
		Name: "",
		Text: "",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Не удалось сериализовать CreateTaskRequest с пустыми полями: %v", err)
	}

	var decoded CreateTaskRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать CreateTaskRequest с пустыми полями: %v", err)
	}

	if decoded.Name != "" {
		t.Errorf("Пустой Name должен быть '', получено '%s'", decoded.Name)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ User
// ============================================================================

func TestUserMarshaling(t *testing.T) {
	user := User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "$2a$12$hash",
		CreatedAt:    time.Now(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Не удалось сериализовать User: %v", err)
	}

	var decoded User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать User: %v", err)
	}

	if decoded.ID != user.ID {
		t.Errorf("Неправильный ID: получено %d, ожидается %d", decoded.ID, user.ID)
	}
	if decoded.Username != user.Username {
		t.Errorf("Неправильный Username: получено %s, ожидается %s", decoded.Username, user.Username)
	}
}

func TestUserPasswordHashOmitEmpty(t *testing.T) {
	user := User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "", // Пустой хеш
		CreatedAt:    time.Now(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Не удалось сериализовать User: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Не удалось десериализовать User в map: %v", err)
	}

	// С omitempty пустой password_hash не должен быть в JSON
	if _, exists := jsonMap["password_hash"]; exists {
		t.Error("Пустой password_hash не должен присутствовать в JSON (omitempty)")
	}
}

func TestUserWithPasswordHash(t *testing.T) {
	user := User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "$2a$12$somehash",
		CreatedAt:    time.Now(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Не удалось сериализовать User: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Не удалось десериализовать User в map: %v", err)
	}

	// Непустой password_hash должен быть в JSON
	if _, exists := jsonMap["password_hash"]; !exists {
		t.Error("Непустой password_hash должен присутствовать в JSON")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ RegisterRequest
// ============================================================================

func TestRegisterRequestMarshaling(t *testing.T) {
	req := RegisterRequest{
		Username: "newuser",
		Password: "password123",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Не удалось сериализовать RegisterRequest: %v", err)
	}

	var decoded RegisterRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать RegisterRequest: %v", err)
	}

	if decoded.Username != req.Username {
		t.Errorf("Неправильный Username: получено %s, ожидается %s", decoded.Username, req.Username)
	}
	if decoded.Password != req.Password {
		t.Errorf("Неправильный Password: получено %s, ожидается %s", decoded.Password, req.Password)
	}
}

func TestRegisterRequestJSONFields(t *testing.T) {
	req := RegisterRequest{
		Username: "user",
		Password: "pass",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Не удалось сериализовать RegisterRequest: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Не удалось десериализовать RegisterRequest в map: %v", err)
	}

	expectedFields := []string{"username", "password"}
	for _, field := range expectedFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("Отсутствует поле %s в JSON", field)
		}
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ LoginRequest
// ============================================================================

func TestLoginRequestMarshaling(t *testing.T) {
	req := LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Не удалось сериализовать LoginRequest: %v", err)
	}

	var decoded LoginRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать LoginRequest: %v", err)
	}

	if decoded.Username != req.Username {
		t.Errorf("Неправильный Username: получено %s, ожидается %s", decoded.Username, req.Username)
	}
	if decoded.Password != req.Password {
		t.Errorf("Неправильный Password: получено %s, ожидается %s", decoded.Password, req.Password)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ AuthResponse
// ============================================================================

func TestAuthResponseMarshaling(t *testing.T) {
	resp := AuthResponse{
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		Username: "testuser",
		UserID:   42,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Не удалось сериализовать AuthResponse: %v", err)
	}

	var decoded AuthResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать AuthResponse: %v", err)
	}

	if decoded.Token != resp.Token {
		t.Errorf("Неправильный Token: получено %s, ожидается %s", decoded.Token, resp.Token)
	}
	if decoded.Username != resp.Username {
		t.Errorf("Неправильный Username: получено %s, ожидается %s", decoded.Username, resp.Username)
	}
	if decoded.UserID != resp.UserID {
		t.Errorf("Неправильный UserID: получено %d, ожидается %d", decoded.UserID, resp.UserID)
	}
}

func TestAuthResponseJSONFields(t *testing.T) {
	resp := AuthResponse{
		Token:    "token",
		Username: "user",
		UserID:   1,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Не удалось сериализовать AuthResponse: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Не удалось десериализовать AuthResponse в map: %v", err)
	}

	expectedFields := []string{"token", "username", "user_id"}
	for _, field := range expectedFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("Отсутствует поле %s в JSON", field)
		}
	}
}

// ============================================================================
// ТЕСТЫ СО СПЕЦИАЛЬНЫМИ СИМВОЛАМИ
// ============================================================================

func TestModelsWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		username string
	}{
		{"cyrillic", "Пользователь123"},
		{"special_chars", "user_name-123"},
		{"mixed", "user@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := RegisterRequest{
				Username: tt.username,
				Password: "pass123",
			}

			data, err := json.Marshal(req)
			if err != nil {
				t.Fatalf("Не удалось сериализовать RegisterRequest: %v", err)
			}

			var decoded RegisterRequest
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("Не удалось десериализовать RegisterRequest: %v", err)
			}

			if decoded.Username != req.Username {
				t.Errorf("Username не совпадает: получено %s, ожидается %s", decoded.Username, req.Username)
			}
		})
	}
}

// ============================================================================
// ТЕСТЫ СОВМЕСТИМОСТИ СТРУКТУР
// ============================================================================

func TestTaskLifecycle(t *testing.T) {
	// Создаем запрос на создание задачи
	createReq := CreateTaskRequest{
		Name: "Test Task",
		Text: "Description",
	}

	// Имитируем создание задачи
	task := Task{
		ID:         1,
		Name:       createReq.Name,
		Text:       createReq.Text,
		CreateTime: time.Now(),
		Complete:   false,
		CompleteAt: nil,
	}

	if task.Name != createReq.Name {
		t.Errorf("Имя задачи не совпадает с запросом: получено %s, ожидается %s", task.Name, createReq.Name)
	}
	if task.Text != createReq.Text {
		t.Errorf("Текст задачи не совпадает с запросом: получено %s, ожидается %s", task.Text, createReq.Text)
	}

	// Завершаем задачу
	completeTime := time.Now()
	task.Complete = true
	task.CompleteAt = &completeTime

	if !task.Complete {
		t.Error("Задача должна быть завершена")
	}
	if task.CompleteAt == nil {
		t.Error("CompleteAt не должен быть nil для завершенной задачи")
	}
}

func TestAuthWorkflow(t *testing.T) {
	// Регистрация
	registerReq := RegisterRequest{
		Username: "newuser",
		Password: "password123",
	}

	// Создаем пользователя (имитация)
	user := User{
		ID:           1,
		Username:     registerReq.Username,
		PasswordHash: "$2a$12$hashedpassword",
		CreatedAt:    time.Now(),
	}

	// Логин
	loginReq := LoginRequest{
		Username: registerReq.Username,
		Password: registerReq.Password,
	}

	if loginReq.Username != user.Username {
		t.Errorf("Username не совпадает: получено %s, ожидается %s", loginReq.Username, user.Username)
	}

	// Ответ авторизации
	authResp := AuthResponse{
		Token:    "jwt.token.here",
		Username: user.Username,
		UserID:   user.ID,
	}

	if authResp.Username != user.Username {
		t.Errorf("Username в ответе не совпадает: получено %s, ожидается %s", authResp.Username, user.Username)
	}
	if authResp.UserID != user.ID {
		t.Errorf("UserID в ответе не совпадает: получено %d, ожидается %d", authResp.UserID, user.ID)
	}
}
