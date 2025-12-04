package auth

import (
	"testing"
	"time"
)

// ============================================================================
// ТЕСТЫ ДЛЯ HASHING ПАРОЛЕЙ
// ============================================================================

// TestHashPassword проверяет успешное хеширование пароля
func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword() вернул ошибку: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() вернул пустой хеш")
	}

	if hash == password {
		t.Error("HashPassword() вернул незашифрованный пароль!")
	}

	if len(hash) < 4 || hash[:4] != "$2a$" {
		t.Errorf("HashPassword() вернул неправильный формат хеша: %s", hash)
	}

	if len(hash) < 59 || len(hash) > 61 {
		t.Errorf("HashPassword() вернул хеш неправильной длины: %d (ожидается ~60)", len(hash))
	}

	if err := CheckPassword(password, hash); err != nil {
		t.Errorf("CheckPassword() не смог проверить хеш, созданный HashPassword(): %v", err)
	}
}

// TestHashPasswordEmpty проверяет хеширование пустого пароля
func TestHashPasswordEmpty(t *testing.T) {
	password := ""

	hash, err := HashPassword(password)

	// Bcrypt должен обработать пустой пароль без ошибки
	if err != nil {
		t.Fatalf("HashPassword() вернул ошибку для пустого пароля: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() вернул пустой хеш для пустого пароля")
	}
}

// TestHashPasswordLong проверяет хеширование очень длинного пароля
func TestHashPasswordLong(t *testing.T) {
	// Bcrypt имеет лимит в 72 байта и должен вернуть ошибку
	password := "verylongpasswordverylongpasswordverylongpasswordverylongpasswordverylongpassword123456789"

	hash, err := HashPassword(password)

	// Ожидаем ошибку для слишком длинного пароля
	if err == nil {
		t.Error("HashPassword() не вернул ошибку для пароля длиннее 72 байт")
	}

	if hash != "" {
		t.Error("HashPassword() вернул хеш для слишком длинного пароля")
	}
}

// TestHashPasswordDeterministic проверяет, что один пароль создает разные хеши
func TestHashPasswordDeterministic(t *testing.T) {
	password := "samePassword123"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Хеши должны быть разными из-за соли
	if hash1 == hash2 {
		t.Error("HashPassword() создает одинаковые хеши для одного пароля (плохо!)")
	}

	// Но оба должны проходить проверку
	if err := CheckPassword(password, hash1); err != nil {
		t.Error("CheckPassword() не прошел для hash1")
	}
	if err := CheckPassword(password, hash2); err != nil {
		t.Error("CheckPassword() не прошел для hash2")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ ПРОВЕРКИ ПАРОЛЕЙ
// ============================================================================

// TestCheckPasswordCorrect проверяет валидацию правильного пароля
func TestCheckPasswordCorrect(t *testing.T) {
	password := "correctPassword123"
	hash, _ := HashPassword(password)

	err := CheckPassword(password, hash)

	if err != nil {
		t.Errorf("CheckPassword() вернул ошибку для правильного пароля: %v", err)
	}
}

// TestCheckPasswordWrong проверяет валидацию неправильного пароля
func TestCheckPasswordWrong(t *testing.T) {
	password := "correctPassword123"
	wrongPassword := "wrongPassword123"
	hash, _ := HashPassword(password)

	err := CheckPassword(wrongPassword, hash)

	if err == nil {
		t.Error("CheckPassword() не вернул ошибку для неправильного пароля")
	}
}

// TestCheckPasswordEmpty проверяет проверку с пустым паролем
func TestCheckPasswordEmpty(t *testing.T) {
	password := "password123"
	hash, _ := HashPassword(password)

	err := CheckPassword("", hash)

	if err == nil {
		t.Error("CheckPassword() не вернул ошибку для пустого пароля")
	}
}

// TestCheckPasswordInvalidHash проверяет проверку с невалидным хешем
func TestCheckPasswordInvalidHash(t *testing.T) {
	password := "password123"
	invalidHash := "invalid_hash"

	err := CheckPassword(password, invalidHash)

	if err == nil {
		t.Error("CheckPassword() не вернул ошибку для невалидного хеша")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ ГЕНЕРАЦИИ JWT ТОКЕНОВ
// ============================================================================

// TestGenerateTokenSuccess проверяет успешную генерацию токена
func TestGenerateTokenSuccess(t *testing.T) {
	userID := 123
	username := "testuser"

	token, err := GenerateToken(userID, username)

	if err != nil {
		t.Fatalf("GenerateToken() вернул ошибку: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() вернул пустой токен")
	}

	// Токен должен состоять из 3 частей, разделенных точками
	// header.payload.signature
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("GenerateToken() вернул токен неправильного формата (ожидается 2 точки, получено %d)", parts)
	}
}

// TestGenerateTokenEmptyUsername проверяет генерацию с пустым username
func TestGenerateTokenEmptyUsername(t *testing.T) {
	userID := 123
	username := ""

	token, err := GenerateToken(userID, username)

	// Должен сгенерировать токен даже с пустым username
	if err != nil {
		t.Fatalf("GenerateToken() вернул ошибку для пустого username: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() вернул пустой токен")
	}
}

// TestGenerateTokenZeroUserID проверяет генерацию с userID = 0
func TestGenerateTokenZeroUserID(t *testing.T) {
	userID := 0
	username := "testuser"

	token, err := GenerateToken(userID, username)

	// Должен сгенерировать токен даже с userID = 0
	if err != nil {
		t.Fatalf("GenerateToken() вернул ошибку для userID = 0: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() вернул пустой токен")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ ВАЛИДАЦИИ JWT ТОКЕНОВ
// ============================================================================

// TestValidateTokenValid проверяет валидацию правильного токена
func TestValidateTokenValid(t *testing.T) {
	userID := 123
	username := "testuser"

	token, _ := GenerateToken(userID, username)

	claims, err := ValidateToken(token)

	if err != nil {
		t.Fatalf("ValidateToken() вернул ошибку для валидного токена: %v", err)
	}

	if claims == nil {
		t.Fatal("ValidateToken() вернул nil claims")
	}

	if claims.UserID != userID {
		t.Errorf("ValidateToken() вернул неправильный UserID: получено %d, ожидается %d", claims.UserID, userID)
	}

	if claims.Username != username {
		t.Errorf("ValidateToken() вернул неправильный Username: получено %s, ожидается %s", claims.Username, username)
	}
}

// TestValidateTokenInvalid проверяет валидацию невалидного токена
func TestValidateTokenInvalid(t *testing.T) {
	invalidToken := "invalid.token.here"

	claims, err := ValidateToken(invalidToken)

	if err == nil {
		t.Error("ValidateToken() не вернул ошибку для невалидного токена")
	}

	if claims != nil {
		t.Error("ValidateToken() вернул claims для невалидного токена")
	}
}

// TestValidateTokenEmpty проверяет валидацию пустого токена
func TestValidateTokenEmpty(t *testing.T) {
	emptyToken := ""

	claims, err := ValidateToken(emptyToken)

	if err == nil {
		t.Error("ValidateToken() не вернул ошибку для пустого токена")
	}

	if claims != nil {
		t.Error("ValidateToken() вернул claims для пустого токена")
	}
}

// TestValidateTokenMalformed проверяет валидацию некорректно сформированного токена
func TestValidateTokenMalformed(t *testing.T) {
	malformedToken := "not-a-jwt-token"

	claims, err := ValidateToken(malformedToken)

	if err == nil {
		t.Error("ValidateToken() не вернул ошибку для некорректного токена")
	}

	if claims != nil {
		t.Error("ValidateToken() вернул claims для некорректного токена")
	}
}

// TestValidateTokenExpiry проверяет, что токен содержит время истечения
func TestValidateTokenExpiry(t *testing.T) {
	userID := 123
	username := "testuser"

	token, _ := GenerateToken(userID, username)
	claims, _ := ValidateToken(token)

	// Проверяем, что ExpiresAt установлен
	if claims.ExpiresAt == nil {
		t.Error("ValidateToken() вернул claims без ExpiresAt")
	}

	// Проверяем, что ExpiresAt в будущем
	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("ValidateToken() вернул claims с ExpiresAt в прошлом")
	}

	// Проверяем, что ExpiresAt примерно через 24 часа (±1 минута)
	expectedExpiry := time.Now().Add(24 * time.Hour)
	diff := claims.ExpiresAt.Time.Sub(expectedExpiry)
	if diff < -time.Minute || diff > time.Minute {
		t.Errorf("ValidateToken() вернул claims с неправильным ExpiresAt: разница %v", diff)
	}
}
