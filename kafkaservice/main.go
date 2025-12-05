package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/IBM/sarama"
)

// KafkaConfig holds kafka configuration
type KafkaConfig struct {
	Brokers []string
	Topic   string
	LogFile string
}

// MessageHandler handles consumed messages
type MessageHandler interface {
	HandleMessage(msg *sarama.ConsumerMessage) error
}

// FileMessageHandler writes messages to file
type FileMessageHandler struct {
	writer io.Writer
}

// NewFileMessageHandler creates a new file message handler
func NewFileMessageHandler(writer io.Writer) *FileMessageHandler {
	return &FileMessageHandler{writer: writer}
}

// HandleMessage writes message to file
func (h *FileMessageHandler) HandleMessage(msg *sarama.ConsumerMessage) error {
	logEntry := string(msg.Value)
	log.Printf("Received event: %s", logEntry)
	_, err := h.writer.Write([]byte(logEntry + "\n"))
	return err
}

// GetKafkaConfig reads configuration from environment
func GetKafkaConfig() KafkaConfig {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "kafka:29092"
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "task-events"
	}

	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "/app/logs/events.log"
	}

	return KafkaConfig{
		Brokers: strings.Split(brokers, ","),
		Topic:   topic,
		LogFile: logFile,
	}
}

// CreateSaramaConfig creates Sarama consumer config
func CreateSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	return config
}

// EnsureLogDirectory creates log directory if it doesn't exist
func EnsureLogDirectory(logFile string) error {
	// Extract directory from log file path
	dir := filepath.Dir(logFile)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

// ConsumeMessages consumes messages from partition consumer
func ConsumeMessages(ctx context.Context, pc sarama.PartitionConsumer, handler MessageHandler) {
	for {
		select {
		case msg := <-pc.Messages():
			if err := handler.HandleMessage(msg); err != nil {
				log.Printf("Failed to handle message: %v", err)
			}

		case err := <-pc.Errors():
			log.Printf("Consumer error: %v", err)

		case <-ctx.Done():
			return
		}
	}
}

func main() {
	config := GetKafkaConfig()

	// Ensure log directory exists
	if err := EnsureLogDirectory(config.LogFile); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Open log file
	file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	log.Printf("Kafka consumer starting, brokers: %v, topic: %s", config.Brokers, config.Topic)

	// Configure consumer
	saramaConfig := CreateSaramaConfig()

	// Create consumer
	consumer, err := sarama.NewConsumer(config.Brokers, saramaConfig)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Get partition consumer
	partitionConsumer, err := consumer.ConsumePartition(config.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to create partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Setup signal handling
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Kafka consumer started successfully, waiting for messages...")

	// Create message handler
	handler := NewFileMessageHandler(file)

	// Consume messages
	go ConsumeMessages(ctx, partitionConsumer, handler)

	// Wait for termination signal
	<-signals
	log.Println("Shutting down consumer...")
	cancel()
}
