package auth

import (
	"testing"
)

// TestHashPassword проверяет, что пароль правильно хешируется
func TestHashPassword(t *testing.T) {
	// Arrange (подготовка)
	password := "mySecurePassword123"

	// Act (действие)
	hash, err := HashPassword(password)

	// Assert (проверки)

	// 1. Проверяем, что нет ошибки
	if err != nil {
		t.Fatalf("HashPassword() вернул ошибку: %v", err)
	}

	// 2. Проверяем, что хеш не пустой
	if hash == "" {
		t.Error("HashPassword() вернул пустой хеш")
	}

	// 3. Проверяем, что хеш НЕ равен оригинальному паролю
	if hash == password {
		t.Error("HashPassword() вернул незашифрованный пароль!")
	}

	// 4. Проверяем, что хеш начинается с префикса bcrypt "$2a$"
	if len(hash) < 4 || hash[:4] != "$2a$" {
		t.Errorf("HashPassword() вернул неправильный формат хеша: %s", hash)
	}

	// 5. Проверяем, что хеш имеет правильную длину (около 60 символов)
	if len(hash) < 59 || len(hash) > 61 {
		t.Errorf("HashPassword() вернул хеш неправильной длины: %d (ожидается ~60)", len(hash))
	}

	// 6. Проверяем, что хеш можно проверить с помощью CheckPassword
	if err := CheckPassword(password, hash); err != nil {
		t.Errorf("CheckPassword() не смог проверить хеш, созданный HashPassword(): %v", err)
	}
}
