package middleware

import (
	"apiservice/auth"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// setUserContext добавляет claims в контекст
func setUserContext(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, UserContextKey, claims)
}

// setUserContextWrongType добавляет неправильный тип в контекст для тестирования
func setUserContextWrongType(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, UserContextKey, value)
}

// ============================================================================
// ТЕСТЫ ДЛЯ AuthMiddleware
// ============================================================================

func TestAuthMiddlewareSuccess(t *testing.T) {
	// Создаем валидный токен
	token, err := auth.GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("Не удалось создать токен: %v", err)
	}

	// Создаем тестовый handler, который проверяет контекст
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r)
		if claims == nil {
			t.Error("Контекст не содержит пользователя")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if claims.Username != "testuser" {
			t.Errorf("Неправильный username: получено %s, ожидается testuser", claims.Username)
		}
		if claims.UserID != 1 {
			t.Errorf("Неправильный userID: получено %d, ожидается 1", claims.UserID)
		}
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем в middleware
	handler := AuthMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("AuthMiddleware() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}
}

func TestAuthMiddlewareMissingHeader(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler не должен был быть вызван")
	})

	handler := AuthMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	// Не устанавливаем Authorization header

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("AuthMiddleware() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
	}

	expectedError := "Authorization header required"
	if rr.Body.String() != expectedError+"\n" {
		t.Errorf("AuthMiddleware() вернул неправильное сообщение: получено %q, ожидается %q", rr.Body.String(), expectedError)
	}
}

func TestAuthMiddlewareInvalidFormat(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"no_bearer_prefix", "sometoken"},
		{"wrong_prefix", "Basic sometoken"},
		{"no_token", "Bearer"},
		{"extra_parts", "Bearer token extra"},
		{"empty_token", "Bearer "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("Next handler не должен был быть вызван")
			})

			handler := AuthMiddleware(nextHandler)

			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.header)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("AuthMiddleware() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
			}
		})
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler не должен был быть вызван")
	})

	handler := AuthMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("AuthMiddleware() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
	}

	expectedError := "Invalid or expired token"
	if rr.Body.String() != expectedError+"\n" {
		t.Errorf("AuthMiddleware() вернул неправильное сообщение: получено %q, ожидается %q", rr.Body.String(), expectedError)
	}
}

func TestAuthMiddlewareExpiredToken(t *testing.T) {
	// Используем заведомо истекший токен (создан в прошлом)
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxfQ.invalidSignature"

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler не должен был быть вызван")
	})

	handler := AuthMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("AuthMiddleware() вернул неправильный статус: получено %v, ожидается %v", status, http.StatusUnauthorized)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetUserFromContext
// ============================================================================

func TestGetUserFromContextSuccess(t *testing.T) {
	claims := &auth.Claims{
		UserID:   42,
		Username: "testuser",
	}

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := req.Context()
	ctx = setUserContext(ctx, claims)
	req = req.WithContext(ctx)

	result := GetUserFromContext(req)
	if result == nil {
		t.Fatal("GetUserFromContext() вернул nil")
	}

	if result.UserID != 42 {
		t.Errorf("Неправильный UserID: получено %d, ожидается 42", result.UserID)
	}
	if result.Username != "testuser" {
		t.Errorf("Неправильный Username: получено %s, ожидается testuser", result.Username)
	}
}

func TestGetUserFromContextMissing(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	// Не добавляем контекст

	result := GetUserFromContext(req)
	if result != nil {
		t.Error("GetUserFromContext() должен вернуть nil для отсутствующего контекста")
	}
}

func TestGetUserFromContextWrongType(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := req.Context()
	// Добавляем неправильный тип в контекст
	ctx = setUserContextWrongType(ctx, "not a claims object")
	req = req.WithContext(ctx)

	result := GetUserFromContext(req)
	if result != nil {
		t.Error("GetUserFromContext() должен вернуть nil для неправильного типа")
	}
}

// ============================================================================
// ИНТЕГРАЦИОННЫЕ ТЕСТЫ
// ============================================================================

func TestAuthMiddlewareChain(t *testing.T) {
	// Создаем валидный токен
	token, err := auth.GenerateToken(99, "chainuser")
	if err != nil {
		t.Fatalf("Не удалось создать токен: %v", err)
	}

	// Создаем цепочку handlers
	var middlewareCalled, handlerCalled bool

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		claims := GetUserFromContext(r)
		if claims == nil {
			t.Error("Контекст не содержит пользователя в финальном handler")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Дополнительный middleware для проверки цепочки
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	handler := loggingMiddleware(AuthMiddleware(finalHandler))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !middlewareCalled {
		t.Error("Logging middleware не был вызван")
	}
	if !handlerCalled {
		t.Error("Финальный handler не был вызван")
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неправильный статус: получено %v, ожидается %v", status, http.StatusOK)
	}
	if rr.Body.String() != "success" {
		t.Errorf("Неправильный ответ: получено %q, ожидается %q", rr.Body.String(), "success")
	}
}
