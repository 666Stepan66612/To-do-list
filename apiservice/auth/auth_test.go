package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// ============================================================================
// ДОПОЛНИТЕЛЬНЫЕ ГРАНИЧНЫЕ ТЕСТЫ
// ============================================================================

func TestHashPasswordVeryLong(t *testing.T) {
	// Длинный пароль (70 символов - в пределах bcrypt лимита 72 байта)
	longPassword := ""
	for i := 0; i < 70; i++ {
		longPassword += "a"
	}

	hash, err := HashPassword(longPassword)
	if err != nil {
		t.Fatalf("HashPassword() вернул ошибку для длинного пароля: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() вернул пустой хеш")
	}

	// Проверяем, что можем проверить длинный пароль
	err = CheckPassword(longPassword, hash)
	if err != nil {
		t.Errorf("CheckPassword() не смог проверить длинный пароль: %v", err)
	}
}

func TestCheckPasswordWithWrongHash(t *testing.T) {
	password := "correctPassword"
	wrongHash := "not_a_valid_bcrypt_hash"

	err := CheckPassword(password, wrongHash)
	if err == nil {
		t.Error("CheckPassword() должен вернуть ошибку для невалидного хеша")
	}
}

func TestGenerateTokenNegativeUserID(t *testing.T) {
	token, err := GenerateToken(-1, "testuser")
	if err != nil {
		t.Fatalf("GenerateToken() вернул ошибку для отрицательного userID: %v", err)
	}

	// Должен сгенерировать токен даже с отрицательным ID
	if token == "" {
		t.Error("GenerateToken() вернул пустой токен")
	}

	// Проверяем, что токен валиден
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() не смог проверить токен: %v", err)
	}

	if claims.UserID != -1 {
		t.Errorf("Неправильный UserID: получено %d, ожидается -1", claims.UserID)
	}
}

func TestGenerateTokenWithSpecialCharsUsername(t *testing.T) {
	specialUsername := "user@example.com!#$%"
	token, err := GenerateToken(1, specialUsername)
	if err != nil {
		t.Fatalf("GenerateToken() вернул ошибку для username со спецсимволами: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() не смог проверить токен: %v", err)
	}

	if claims.Username != specialUsername {
		t.Errorf("Username не совпадает: получено %s, ожидается %s", claims.Username, specialUsername)
	}
}

func TestValidateTokenWithDifferentSecret(t *testing.T) {
	// Создаем токен
	token, err := GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("GenerateToken() вернул ошибку: %v", err)
	}

	// Пытаемся изменить токен (подменить данные)
	fakeToken := token + "fake"

	_, err = ValidateToken(fakeToken)
	if err == nil {
		t.Error("ValidateToken() должен вернуть ошибку для измененного токена")
	}
}

func TestValidateTokenPartialToken(t *testing.T) {
	token, _ := GenerateToken(1, "testuser")

	// Берем только часть токена
	partialToken := token[:len(token)/2]

	_, err := ValidateToken(partialToken)
	if err == nil {
		t.Error("ValidateToken() должен вернуть ошибку для частичного токена")
	}
}

func TestHashPasswordUnicode(t *testing.T) {
	unicodePassword := "Пароль123!@#密码"

	hash, err := HashPassword(unicodePassword)
	if err != nil {
		t.Fatalf("HashPassword() вернул ошибку для unicode пароля: %v", err)
	}

	err = CheckPassword(unicodePassword, hash)
	if err != nil {
		t.Errorf("CheckPassword() не смог проверить unicode пароль: %v", err)
	}
}

func TestCheckPasswordCaseSensitive(t *testing.T) {
	password := "MyPassword123"
	hash, _ := HashPassword(password)

	// Проверяем с другим регистром
	err := CheckPassword("mypassword123", hash)
	if err == nil {
		t.Error("CheckPassword() должен быть чувствителен к регистру")
	}
}

func TestValidateTokenWithInvalidSigningMethod(t *testing.T) {
	// Создаем токен с другим методом подписи (RS256 вместо HS256)
	claims := &Claims{
		UserID:   1,
		Username: "test",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Используем неправильный метод подписи
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	_, err := ValidateToken(tokenString)
	if err == nil {
		t.Error("ValidateToken() должен отклонить токен с неправильным методом подписи")
	}
}

func TestHashPasswordWithSpecialCharacters(t *testing.T) {
	passwords := []string{
		"!@#$%^&*()",
		"pass<>word",
		"pass\"word",
		"pass'word",
		"pass\nword",
		"pass\tword",
	}

	for _, pwd := range passwords {
		hash, err := HashPassword(pwd)
		if err != nil {
			t.Errorf("HashPassword() не смог хешировать пароль со спецсимволами %q: %v", pwd, err)
		}

		if err := CheckPassword(pwd, hash); err != nil {
			t.Errorf("CheckPassword() не смог проверить пароль со спецсимволами %q: %v", pwd, err)
		}
	}
}

func TestGenerateTokenWithVeryLongUsername(t *testing.T) {
	longUsername := string(make([]byte, 10000))
	for i := range longUsername {
		longUsername = longUsername[:i] + "a" + longUsername[i:]
	}
	longUsername = longUsername[:10000]

	token, err := GenerateToken(1, longUsername[:10000])
	if err != nil {
		t.Errorf("GenerateToken() не смог создать токен с длинным username: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() вернул пустой токен")
	}
}

func TestValidateTokenClaimsExtraction(t *testing.T) {
	userID := 42
	username := "testuser"

	token, err := GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("Не удалось создать токен: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() вернул ошибку: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID не совпадает: получено %d, ожидается %d", claims.UserID, userID)
	}

	if claims.Username != username {
		t.Errorf("Username не совпадает: получено %s, ожидается %s", claims.Username, username)
	}

	if claims.ExpiresAt == nil {
		t.Error("ExpiresAt не установлен")
	}

	if claims.IssuedAt == nil {
		t.Error("IssuedAt не установлен")
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "testpass"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("Ошибка хеширования: %v, %v", err1, err2)
	}

	// Хеши должны быть разными (bcrypt использует соль)
	if hash1 == hash2 {
		t.Error("Хеши одинаковы, должны отличаться из-за соли")
	}

	// Но оба должны проходить проверку
	if err := CheckPassword(password, hash1); err != nil {
		t.Errorf("CheckPassword не прошла для hash1: %v", err)
	}

	if err := CheckPassword(password, hash2); err != nil {
		t.Errorf("CheckPassword не прошла для hash2: %v", err)
	}
}

func TestCheckPasswordWithDifferentPasswords(t *testing.T) {
	password1 := "password1"
	password2 := "password2"

	hash1, _ := HashPassword(password1)

	// Проверка правильного пароля
	if err := CheckPassword(password1, hash1); err != nil {
		t.Error("CheckPassword должна пройти для правильного пароля")
	}

	// Проверка неправильного пароля
	if err := CheckPassword(password2, hash1); err == nil {
		t.Error("CheckPassword должна вернуть ошибку для неправильного пароля")
	}
}

func TestGenerateTokenWithEmptyStrings(t *testing.T) {
	token, err := GenerateToken(0, "")
	if err != nil {
		t.Fatalf("GenerateToken не смог создать токен: %v", err)
	}

	if token == "" {
		t.Error("Token не должен быть пустым")
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken вернул ошибку: %v", err)
	}

	if claims.Username != "" {
		t.Errorf("Username должен быть пустым, получено: %s", claims.Username)
	}
}

func TestValidateTokenWithManipulatedClaims(t *testing.T) {
	// Создаем валидный токен
	token, _ := GenerateToken(1, "user")

	// Пытаемся изменить токен (добавляем символ)
	manipulated := token + "x"

	_, err := ValidateToken(manipulated)
	if err == nil {
		t.Error("ValidateToken должен отклонить измененный токен")
	}
}

func TestCheckPasswordEmptyPassword(t *testing.T) {
	hash, err := HashPassword("validpassword")
	if err != nil {
		t.Fatalf("Не удалось создать хеш: %v", err)
	}

	err = CheckPassword("", hash)
	if err == nil {
		t.Error("CheckPassword должен вернуть ошибку для пустого пароля")
	}
}

func TestValidateTokenExpired(t *testing.T) {
	// Создаем токен с истекшим временем
	claims := &Claims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)), // истек вчера
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)

	_, err := ValidateToken(tokenString)
	if err == nil {
		t.Error("ValidateToken должен отклонить истекший токен")
	}
}

func TestValidateTokenWithWrongClaimsType(t *testing.T) {
	// Создаем невалидный токен с правильным форматом но неправильными данными
	// Токен в формате header.payload.signature но с пустым payload
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.invalid"

	_, err := ValidateToken(invalidToken)
	if err == nil {
		t.Error("ValidateToken должен вернуть ошибку для невалидного токена")
	}
}
