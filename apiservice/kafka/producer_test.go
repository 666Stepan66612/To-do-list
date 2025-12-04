package kafka

import (
	"encoding/json"
	"testing"
	"time"
)

// ============================================================================
// ТЕСТЫ ДЛЯ Event СТРУКТУРЫ
// ============================================================================

func TestEventMarshaling(t *testing.T) {
	event := Event{
		Timestamp: "2024-12-04T10:00:00Z",
		UserID:    42,
		Username:  "testuser",
		Action:    "CREATE_TASK",
		Details:   "Task created: Test Task",
		Status:    "SUCCESS",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event: %v", err)
	}

	if decoded.UserID != event.UserID {
		t.Errorf("Неправильный UserID после десериализации: получено %d, ожидается %d", decoded.UserID, event.UserID)
	}
	if decoded.Username != event.Username {
		t.Errorf("Неправильный Username после десериализации: получено %s, ожидается %s", decoded.Username, event.Username)
	}
	if decoded.Action != event.Action {
		t.Errorf("Неправильный Action после десериализации: получено %s, ожидается %s", decoded.Action, event.Action)
	}
	if decoded.Status != event.Status {
		t.Errorf("Неправильный Status после десериализации: получено %s, ожидается %s", decoded.Status, event.Status)
	}
}

func TestEventJSONFields(t *testing.T) {
	event := Event{
		Timestamp: "2024-12-04T10:00:00Z",
		UserID:    1,
		Username:  "user",
		Action:    "TEST_ACTION",
		Details:   "Test details",
		Status:    "SUCCESS",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Не удалось десериализовать Event в map: %v", err)
	}

	// Проверяем наличие всех полей
	expectedFields := []string{"timestamp", "user_id", "username", "action", "details", "status"}
	for _, field := range expectedFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("Отсутствует поле %s в JSON", field)
		}
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ SendEvent (без реального Kafka)
// ============================================================================

func TestSendEventNilProducer(t *testing.T) {
	var ep *EventProducer = nil

	// Должен не падать и вернуть nil
	err := ep.SendEvent(1, "testuser", "TEST_ACTION", "details", "SUCCESS")
	if err != nil {
		t.Errorf("SendEvent() с nil producer вернул ошибку: %v", err)
	}
}

func TestSendEventEmptyProducer(t *testing.T) {
	ep := &EventProducer{
		producer: nil,
		topic:    "test-topic",
	}

	// Должен упасть с panic при попытке отправки (nil pointer dereference)
	// Проверяем, что такая ситуация обрабатывается
	defer func() {
		if r := recover(); r == nil {
			t.Error("SendEvent() с nil producer внутри структуры должен вызвать panic")
		}
	}()

	_ = ep.SendEvent(1, "testuser", "TEST_ACTION", "details", "SUCCESS")
}

// ============================================================================
// ТЕСТЫ ДЛЯ ФОРМАТИРОВАНИЯ ДАННЫХ
// ============================================================================

func TestEventTimestampFormat(t *testing.T) {
	timestamp := time.Now().Format(time.RFC3339)

	event := Event{
		Timestamp: timestamp,
		UserID:    1,
		Username:  "user",
		Action:    "TEST",
		Details:   "test",
		Status:    "SUCCESS",
	}

	// Проверяем, что timestamp можно распарсить обратно
	parsedTime, err := time.Parse(time.RFC3339, event.Timestamp)
	if err != nil {
		t.Errorf("Не удалось распарсить timestamp: %v", err)
	}

	if parsedTime.IsZero() {
		t.Error("Timestamp не должен быть нулевым")
	}
}

func TestEventWithDifferentActions(t *testing.T) {
	actions := []string{
		"CREATE_TASK",
		"DELETE_TASK",
		"COMPLETE_TASK",
		"UPDATE_TASK",
		"GET_TASKS",
	}

	for _, action := range actions {
		event := Event{
			Timestamp: time.Now().Format(time.RFC3339),
			UserID:    1,
			Username:  "testuser",
			Action:    action,
			Details:   "Test details",
			Status:    "SUCCESS",
		}

		data, err := json.Marshal(event)
		if err != nil {
			t.Errorf("Не удалось сериализовать Event с action %s: %v", action, err)
			continue
		}

		var decoded Event
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Errorf("Не удалось десериализовать Event с action %s: %v", action, err)
			continue
		}

		if decoded.Action != action {
			t.Errorf("Action не совпадает для %s: получено %s", action, decoded.Action)
		}
	}
}

func TestEventWithDifferentStatuses(t *testing.T) {
	statuses := []string{"SUCCESS", "ERROR", "WARNING", "INFO"}

	for _, status := range statuses {
		event := Event{
			Timestamp: time.Now().Format(time.RFC3339),
			UserID:    1,
			Username:  "testuser",
			Action:    "TEST_ACTION",
			Details:   "Test details",
			Status:    status,
		}

		data, err := json.Marshal(event)
		if err != nil {
			t.Errorf("Не удалось сериализовать Event с status %s: %v", status, err)
			continue
		}

		var decoded Event
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Errorf("Не удалось десериализовать Event с status %s: %v", status, err)
			continue
		}

		if decoded.Status != status {
			t.Errorf("Status не совпадает для %s: получено %s", status, decoded.Status)
		}
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ ВАЛИДАЦИИ ДАННЫХ
// ============================================================================

func TestEventWithEmptyFields(t *testing.T) {
	event := Event{
		Timestamp: "",
		UserID:    0,
		Username:  "",
		Action:    "",
		Details:   "",
		Status:    "",
	}

	// Проверяем, что пустые поля сериализуются корректно
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event с пустыми полями: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event с пустыми полями: %v", err)
	}

	if decoded.UserID != 0 {
		t.Errorf("Пустой UserID должен быть 0, получено %d", decoded.UserID)
	}
	if decoded.Username != "" {
		t.Errorf("Пустой Username должен быть '', получено '%s'", decoded.Username)
	}
}

func TestEventWithLongDetails(t *testing.T) {
	longDetails := ""
	for i := 0; i < 200; i++ {
		longDetails += "Very long details string that simulates a large payload. "
	}

	event := Event{
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    1,
		Username:  "testuser",
		Action:    "TEST_ACTION",
		Details:   longDetails,
		Status:    "SUCCESS",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event с длинными details: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event с длинными details: %v", err)
	}

	if len(decoded.Details) == 0 {
		t.Error("Details не должны быть пустыми после десериализации")
	}
}

func TestEventWithSpecialCharacters(t *testing.T) {
	event := Event{
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    1,
		Username:  "test_user-123",
		Action:    "CREATE_TASK",
		Details:   "Task with special chars: <>&\"'",
		Status:    "SUCCESS",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event со спецсимволами: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event со спецсимволами: %v", err)
	}

	if decoded.Details != event.Details {
		t.Errorf("Details со спецсимволами не совпадают: получено %s, ожидается %s", decoded.Details, event.Details)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ РУССКИХ СИМВОЛОВ (UTF-8)
// ============================================================================

func TestEventWithCyrillicCharacters(t *testing.T) {
	event := Event{
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    1,
		Username:  "Пользователь",
		Action:    "CREATE_TASK",
		Details:   "Задача создана: Тестовая задача",
		Status:    "УСПЕХ",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event с кириллицей: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event с кириллицей: %v", err)
	}

	if decoded.Username != event.Username {
		t.Errorf("Username с кириллицей не совпадает: получено %s, ожидается %s", decoded.Username, event.Username)
	}
	if decoded.Details != event.Details {
		t.Errorf("Details с кириллицей не совпадают: получено %s, ожидается %s", decoded.Details, event.Details)
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ NewEventProducer
// ============================================================================

func TestNewEventProducerWithEmptyBrokers(t *testing.T) {
	_, err := NewEventProducer([]string{}, "test-topic")
	if err == nil {
		t.Error("NewEventProducer() должен вернуть ошибку с пустым списком брокеров")
	}
}

func TestNewEventProducerWithNilBrokers(t *testing.T) {
	_, err := NewEventProducer(nil, "test-topic")
	if err == nil {
		t.Error("NewEventProducer() должен вернуть ошибку с nil брокерами")
	}
}

func TestNewEventProducerWithInvalidBroker(t *testing.T) {
	// Используем недопустимый адрес брокера
	_, err := NewEventProducer([]string{"invalid:99999"}, "test-topic")
	if err == nil {
		t.Error("NewEventProducer() должен вернуть ошибку с недопустимым адресом брокера")
	}
}

func TestNewEventProducerWithEmptyTopic(t *testing.T) {
	// Попытка создать продюсера с пустым топиком
	// Это технически должно работать на этапе создания, но может вызвать проблемы при отправке
	producer, err := NewEventProducer([]string{"localhost:9092"}, "")
	if err == nil && producer != nil {
		// Продюсер создан, но топик пустой - это может быть проблемой при SendEvent
		// Проверяем, что топик сохранен
		if producer.topic != "" {
			t.Errorf("Топик должен быть пустым, получен: %s", producer.topic)
		}
		// Закрываем producer только если он создался
		producer.Close()
	}
	// Если ошибка - это нормально, так как брокер недоступен в тестовой среде
}

func TestNewEventProducerWithInvalidPort(t *testing.T) {
	_, err := NewEventProducer([]string{"localhost:abc"}, "test-topic")
	if err == nil {
		t.Error("NewEventProducer() должен вернуть ошибку с недопустимым портом")
	}
}

func TestNewEventProducerWithMultipleBrokers(t *testing.T) {
	// Тест с несколькими недоступными брокерами
	_, err := NewEventProducer([]string{"broker1:9092", "broker2:9092", "broker3:9092"}, "test-topic")
	// Должна быть ошибка, так как брокеры недоступны
	if err == nil {
		t.Error("NewEventProducer() должен вернуть ошибку с недоступными брокерами")
	}
}

// ============================================================================
// ДОПОЛНИТЕЛЬНЫЕ ТЕСТЫ ДЛЯ Event
// ============================================================================

func TestEventWithZeroUserID(t *testing.T) {
	event := Event{
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    0,
		Username:  "guest",
		Action:    "VIEW",
		Details:   "Viewed page",
		Status:    "SUCCESS",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event с нулевым UserID: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event: %v", err)
	}

	if decoded.UserID != 0 {
		t.Errorf("UserID должен быть 0, получено: %d", decoded.UserID)
	}
}

func TestEventWithNegativeUserID(t *testing.T) {
	event := Event{
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    -1,
		Username:  "invalid",
		Action:    "ERROR",
		Details:   "Invalid user",
		Status:    "ERROR",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event с отрицательным UserID: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event: %v", err)
	}

	if decoded.UserID != -1 {
		t.Errorf("UserID должен быть -1, получено: %d", decoded.UserID)
	}
}

func TestEventWithVeryLongDetails(t *testing.T) {
	// Создаем очень длинную строку деталей
	longDetails := make([]byte, 10000)
	for i := range longDetails {
		longDetails[i] = 'A'
	}

	event := Event{
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    1,
		Username:  "testuser",
		Action:    "BULK_OPERATION",
		Details:   string(longDetails),
		Status:    "SUCCESS",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event с длинными деталями: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event с длинными деталями: %v", err)
	}

	if len(decoded.Details) != len(event.Details) {
		t.Errorf("Длина Details не совпадает: получено %d, ожидается %d", len(decoded.Details), len(event.Details))
	}
}

func TestEventWithAllEmptyFields(t *testing.T) {
	event := Event{
		Timestamp: "",
		UserID:    0,
		Username:  "",
		Action:    "",
		Details:   "",
		Status:    "",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event с пустыми полями: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event с пустыми полями: %v", err)
	}

	if decoded.Username != "" || decoded.Action != "" {
		t.Error("Пустые поля должны оставаться пустыми после десериализации")
	}
}

func TestEventWithSpecialJSONCharacters(t *testing.T) {
	event := Event{
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    1,
		Username:  `user"with"quotes`,
		Action:    "TEST\nNEWLINE\tTAB",
		Details:   `{"nested": "json", "value": 123}`,
		Status:    "SUCCESS",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Не удалось сериализовать Event со специальными символами: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Не удалось десериализовать Event со специальными символами: %v", err)
	}

	if decoded.Username != event.Username {
		t.Errorf("Username со спецсимволами не совпадает: получено %s, ожидается %s", decoded.Username, event.Username)
	}
	if decoded.Details != event.Details {
		t.Errorf("Details со спецсимволами не совпадают: получено %s, ожидается %s", decoded.Details, event.Details)
	}
}
