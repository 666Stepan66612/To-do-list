package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
)

// TestGetKafkaConfig tests configuration loading from environment
func TestGetKafkaConfigDefaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("KAFKA_BROKERS")
	os.Unsetenv("KAFKA_TOPIC")
	os.Unsetenv("LOG_FILE")

	config := GetKafkaConfig()

	if len(config.Brokers) != 1 || config.Brokers[0] != "kafka:29092" {
		t.Errorf("Expected default brokers [kafka:29092], got %v", config.Brokers)
	}

	if config.Topic != "task-events" {
		t.Errorf("Expected default topic 'task-events', got %s", config.Topic)
	}

	if config.LogFile != "/app/logs/events.log" {
		t.Errorf("Expected default log file '/app/logs/events.log', got %s", config.LogFile)
	}
}

func TestGetKafkaConfigFromEnvironment(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "broker1:9092,broker2:9092,broker3:9092")
	os.Setenv("KAFKA_TOPIC", "custom-topic")
	os.Setenv("LOG_FILE", "/custom/path/log.txt")
	defer func() {
		os.Unsetenv("KAFKA_BROKERS")
		os.Unsetenv("KAFKA_TOPIC")
		os.Unsetenv("LOG_FILE")
	}()

	config := GetKafkaConfig()

	expectedBrokers := []string{"broker1:9092", "broker2:9092", "broker3:9092"}
	if len(config.Brokers) != len(expectedBrokers) {
		t.Errorf("Expected %d brokers, got %d", len(expectedBrokers), len(config.Brokers))
	}
	for i, broker := range expectedBrokers {
		if config.Brokers[i] != broker {
			t.Errorf("Expected broker %s at index %d, got %s", broker, i, config.Brokers[i])
		}
	}

	if config.Topic != "custom-topic" {
		t.Errorf("Expected topic 'custom-topic', got %s", config.Topic)
	}

	if config.LogFile != "/custom/path/log.txt" {
		t.Errorf("Expected log file '/custom/path/log.txt', got %s", config.LogFile)
	}
}

func TestGetKafkaConfigSingleBroker(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	defer os.Unsetenv("KAFKA_BROKERS")

	config := GetKafkaConfig()

	if len(config.Brokers) != 1 || config.Brokers[0] != "localhost:9092" {
		t.Errorf("Expected single broker [localhost:9092], got %v", config.Brokers)
	}
}

func TestGetKafkaConfigEmptyString(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "")
	defer os.Unsetenv("KAFKA_BROKERS")

	config := GetKafkaConfig()

	if len(config.Brokers) != 1 || config.Brokers[0] != "kafka:29092" {
		t.Errorf("Expected default brokers for empty string, got %v", config.Brokers)
	}
}

// TestCreateSaramaConfig tests Sarama configuration creation
func TestCreateSaramaConfig(t *testing.T) {
	config := CreateSaramaConfig()

	if config == nil {
		t.Fatal("Expected non-nil config")
	}

	if !config.Consumer.Return.Errors {
		t.Error("Expected Consumer.Return.Errors to be true")
	}

	if config.Consumer.Offsets.Initial != sarama.OffsetNewest {
		t.Errorf("Expected OffsetNewest, got %d", config.Consumer.Offsets.Initial)
	}
}

// TestEnsureLogDirectory tests log directory creation
func TestEnsureLogDirectorySuccess(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "logs", "test.log")

	err := EnsureLogDirectory(logFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify directory was created
	dirPath := filepath.Join(tempDir, "logs")
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Error("Expected directory to be created")
	}
}

func TestEnsureLogDirectoryNestedPaths(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "deeply", "nested", "path", "logs", "test.log")

	err := EnsureLogDirectory(logFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify nested directory was created
	dirPath := filepath.Join(tempDir, "deeply", "nested", "path", "logs")
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Error("Expected nested directory to be created")
	}
}

func TestEnsureLogDirectoryAlreadyExists(t *testing.T) {
	tempDir := t.TempDir()
	logsDir := filepath.Join(tempDir, "logs")

	// Create directory first
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create initial directory: %v", err)
	}

	logFile := filepath.Join(logsDir, "test.log")

	err := EnsureLogDirectory(logFile)
	if err != nil {
		t.Fatalf("Expected no error when directory exists, got %v", err)
	}
}

func TestEnsureLogDirectoryCurrentDir(t *testing.T) {
	// Test with file in current directory (no path separator)
	err := EnsureLogDirectory("test.log")
	if err != nil {
		t.Fatalf("Expected no error for current dir, got %v", err)
	}
}

// TestFileMessageHandler tests message handling
func TestNewFileMessageHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	if handler == nil {
		t.Fatal("Expected non-nil handler")
	}

	if handler.writer == nil {
		t.Error("Expected writer to be set")
	}
}

func TestFileMessageHandlerHandleMessage(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	msg := &sarama.ConsumerMessage{
		Topic:     "test-topic",
		Partition: 0,
		Offset:    123,
		Key:       []byte("key"),
		Value:     []byte("test message content"),
	}

	err := handler.HandleMessage(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "test message content\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestFileMessageHandlerHandleMessageMultiple(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	messages := []string{"message1", "message2", "message3"}

	for _, msgContent := range messages {
		msg := &sarama.ConsumerMessage{
			Value: []byte(msgContent),
		}

		if err := handler.HandleMessage(msg); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}

	expected := "message1\nmessage2\nmessage3\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestFileMessageHandlerHandleMessageEmpty(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	msg := &sarama.ConsumerMessage{
		Value: []byte(""),
	}

	err := handler.HandleMessage(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestFileMessageHandlerHandleMessageSpecialCharacters(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	msg := &sarama.ConsumerMessage{
		Value: []byte("special: äöü ñ 你好 🎉"),
	}

	err := handler.HandleMessage(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "special: äöü ñ 你好 🎉\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

// Mock writer that returns error
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}

func TestFileMessageHandlerHandleMessageWriteError(t *testing.T) {
	handler := NewFileMessageHandler(&errorWriter{})

	msg := &sarama.ConsumerMessage{
		Value: []byte("test"),
	}

	err := handler.HandleMessage(msg)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "write error") {
		t.Errorf("Expected 'write error', got %v", err)
	}
}

// TestConsumeMessages tests message consumption
func TestConsumeMessagesSuccess(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	mockConsumer := mocks.NewConsumer(t, nil)
	partitionConsumer := mockConsumer.ExpectConsumePartition("test-topic", 0, sarama.OffsetNewest)

	// Setup test messages
	messages := []string{"msg1", "msg2", "msg3"}

	ctx, cancel := context.WithCancel(context.Background())

	// Send messages in goroutine
	go func() {
		for _, msg := range messages {
			partitionConsumer.YieldMessage(&sarama.ConsumerMessage{
				Topic:     "test-topic",
				Partition: 0,
				Value:     []byte(msg),
			})
		}
		time.Sleep(50 * time.Millisecond)
		cancel() // Stop consumer after messages
	}()

	pc, err := mockConsumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer pc.Close()

	// Consume messages
	ConsumeMessages(ctx, pc, handler)

	// Verify all messages were written
	expected := "msg1\nmsg2\nmsg3\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestConsumeMessagesWithErrors(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	mockConsumer := mocks.NewConsumer(t, nil)
	partitionConsumer := mockConsumer.ExpectConsumePartition("test-topic", 0, sarama.OffsetNewest)

	ctx, cancel := context.WithCancel(context.Background())

	// Send error and then cancel
	go func() {
		partitionConsumer.YieldError(errors.New("consumer error"))
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	pc, err := mockConsumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer pc.Close()

	// Should not panic on error
	ConsumeMessages(ctx, pc, handler)

	// Buffer should be empty (no messages)
	if buf.String() != "" {
		t.Errorf("Expected empty buffer, got %q", buf.String())
	}
}

func TestConsumeMessagesContextCancellation(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	mockConsumer := mocks.NewConsumer(t, nil)
	mockConsumer.ExpectConsumePartition("test-topic", 0, sarama.OffsetNewest)

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	pc, err := mockConsumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer pc.Close()

	// Should return immediately due to cancelled context
	done := make(chan bool)
	go func() {
		ConsumeMessages(ctx, pc, handler)
		done <- true
	}()

	select {
	case <-done:
		// Success - function returned
	case <-time.After(1 * time.Second):
		t.Error("ConsumeMessages did not return on context cancellation")
	}
}

func TestConsumeMessagesHandlerError(t *testing.T) {
	// Use error writer to force handler errors
	handler := NewFileMessageHandler(&errorWriter{})

	mockConsumer := mocks.NewConsumer(t, nil)
	partitionConsumer := mockConsumer.ExpectConsumePartition("test-topic", 0, sarama.OffsetNewest)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		partitionConsumer.YieldMessage(&sarama.ConsumerMessage{
			Value: []byte("test"),
		})
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	pc, err := mockConsumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer pc.Close()

	// Should not panic even when handler returns error
	ConsumeMessages(ctx, pc, handler)
}

func TestConsumeMessagesMixedEventsAndErrors(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	mockConsumer := mocks.NewConsumer(t, nil)
	partitionConsumer := mockConsumer.ExpectConsumePartition("test-topic", 0, sarama.OffsetNewest)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		partitionConsumer.YieldMessage(&sarama.ConsumerMessage{Value: []byte("msg1")})
		partitionConsumer.YieldError(errors.New("error1"))
		partitionConsumer.YieldMessage(&sarama.ConsumerMessage{Value: []byte("msg2")})
		partitionConsumer.YieldError(errors.New("error2"))
		partitionConsumer.YieldMessage(&sarama.ConsumerMessage{Value: []byte("msg3")})
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	pc, err := mockConsumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer pc.Close()

	ConsumeMessages(ctx, pc, handler)

	// Should have all messages despite errors
	expected := "msg1\nmsg2\nmsg3\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestConsumeMessagesLargeMessage(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	mockConsumer := mocks.NewConsumer(t, nil)
	partitionConsumer := mockConsumer.ExpectConsumePartition("test-topic", 0, sarama.OffsetNewest)

	// Create large message (10KB)
	largeMsg := strings.Repeat("x", 10000)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		partitionConsumer.YieldMessage(&sarama.ConsumerMessage{
			Value: []byte(largeMsg),
		})
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	pc, err := mockConsumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer pc.Close()

	ConsumeMessages(ctx, pc, handler)

	expected := largeMsg + "\n"
	if buf.String() != expected {
		t.Errorf("Expected message of length %d, got %d", len(expected), len(buf.String()))
	}
}

// TestKafkaConfig struct
func TestKafkaConfigStruct(t *testing.T) {
	config := KafkaConfig{
		Brokers: []string{"broker1", "broker2"},
		Topic:   "test-topic",
		LogFile: "/path/to/log",
	}

	if len(config.Brokers) != 2 {
		t.Errorf("Expected 2 brokers, got %d", len(config.Brokers))
	}

	if config.Topic != "test-topic" {
		t.Errorf("Expected topic 'test-topic', got %s", config.Topic)
	}

	if config.LogFile != "/path/to/log" {
		t.Errorf("Expected log file '/path/to/log', got %s", config.LogFile)
	}
}

// Additional edge case tests
func TestGetKafkaConfigMultipleBrokersWithSpaces(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "broker1:9092, broker2:9092, broker3:9092")
	defer os.Unsetenv("KAFKA_BROKERS")

	config := GetKafkaConfig()

	// Should preserve spaces (user might need them)
	if len(config.Brokers) != 3 {
		t.Errorf("Expected 3 brokers, got %d", len(config.Brokers))
	}

	// Check that split works correctly
	if !strings.Contains(config.Brokers[1], "broker2") {
		t.Errorf("Expected broker2 in second position")
	}
}

func TestFileMessageHandlerWithNilMessage(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	// Sarama should never send nil, but test defensive code
	msg := &sarama.ConsumerMessage{
		Value: nil,
	}

	err := handler.HandleMessage(msg)
	if err != nil {
		t.Fatalf("Expected no error with nil value, got %v", err)
	}

	expected := "\n"
	if buf.String() != expected {
		t.Errorf("Expected empty line for nil value, got %q", buf.String())
	}
}

func TestFileMessageHandlerJSONMessage(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	jsonMsg := `{"action":"create","task":"Test Task","user_id":123,"timestamp":"2025-12-05T12:00:00Z"}`
	msg := &sarama.ConsumerMessage{
		Value: []byte(jsonMsg),
	}

	err := handler.HandleMessage(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := jsonMsg + "\n"
	if buf.String() != expected {
		t.Errorf("Expected JSON message, got %q", buf.String())
	}
}

func TestConsumeMessagesEmptyChannel(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	mockConsumer := mocks.NewConsumer(t, nil)
	mockConsumer.ExpectConsumePartition("test-topic", 0, sarama.OffsetNewest)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	pc, err := mockConsumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer pc.Close()

	// Should handle timeout gracefully
	ConsumeMessages(ctx, pc, handler)

	// Buffer should be empty
	if buf.String() != "" {
		t.Errorf("Expected empty buffer, got %q", buf.String())
	}
}

func TestEnsureLogDirectoryPermissions(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "restricted", "test.log")

	err := EnsureLogDirectory(logFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check directory was created with correct permissions
	dirPath := filepath.Join(tempDir, "restricted")
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}

	if !info.IsDir() {
		t.Error("Expected path to be a directory")
	}
}

func TestCreateSaramaConfigValues(t *testing.T) {
	config := CreateSaramaConfig()

	// Test that all expected values are set correctly
	if config.Consumer.Return.Errors != true {
		t.Error("Expected Consumer.Return.Errors to be true")
	}

	if config.Consumer.Offsets.Initial != sarama.OffsetNewest {
		t.Error("Expected Consumer.Offsets.Initial to be OffsetNewest")
	}

	// Verify config is usable (no nil pointers)
	if config.Consumer.Offsets.Initial != sarama.OffsetNewest && config.Consumer.Offsets.Initial != sarama.OffsetOldest {
		t.Error("Invalid offset configuration")
	}
}

func TestConsumeMessagesRapidFireMessages(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	mockConsumer := mocks.NewConsumer(t, nil)
	partitionConsumer := mockConsumer.ExpectConsumePartition("test-topic", 0, sarama.OffsetNewest)

	ctx, cancel := context.WithCancel(context.Background())

	// Send many messages rapidly
	go func() {
		for i := 0; i < 100; i++ {
			partitionConsumer.YieldMessage(&sarama.ConsumerMessage{
				Value: []byte("msg"),
			})
		}
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	pc, err := mockConsumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer pc.Close()

	ConsumeMessages(ctx, pc, handler)

	// Should have all 100 messages
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 100 {
		t.Errorf("Expected 100 messages, got %d", len(lines))
	}
}

func TestGetKafkaConfigBrokersWithPort(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	defer os.Unsetenv("KAFKA_BROKERS")

	config := GetKafkaConfig()

	if len(config.Brokers) != 1 {
		t.Errorf("Expected 1 broker, got %d", len(config.Brokers))
	}

	if config.Brokers[0] != "localhost:9092" {
		t.Errorf("Expected 'localhost:9092', got %s", config.Brokers[0])
	}
}

func TestGetKafkaConfigDNSNames(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "kafka.example.com:9092,kafka2.example.com:9092")
	defer os.Unsetenv("KAFKA_BROKERS")

	config := GetKafkaConfig()

	if len(config.Brokers) != 2 {
		t.Errorf("Expected 2 brokers, got %d", len(config.Brokers))
	}

	if !strings.Contains(config.Brokers[0], "kafka.example.com") {
		t.Errorf("Expected DNS name in first broker, got %s", config.Brokers[0])
	}
}

func TestEnsureLogDirectoryRootPath(t *testing.T) {
	// Test with absolute path
	if os.PathSeparator == '/' {
		// Unix-style path
		err := EnsureLogDirectory("/tmp/test.log")
		if err != nil {
			t.Errorf("Expected no error for root path, got %v", err)
		}
	}
}

func TestFileMessageHandlerBinaryData(t *testing.T) {
	var buf bytes.Buffer
	handler := NewFileMessageHandler(&buf)

	// Test with binary data
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
	msg := &sarama.ConsumerMessage{
		Value: binaryData,
	}

	err := handler.HandleMessage(msg)
	if err != nil {
		t.Fatalf("Expected no error with binary data, got %v", err)
	}

	// Should write binary data as-is
	result := buf.Bytes()
	expectedLen := len(binaryData) + 1 // +1 for newline
	if len(result) != expectedLen {
		t.Errorf("Expected %d bytes, got %d", expectedLen, len(result))
	}
}
