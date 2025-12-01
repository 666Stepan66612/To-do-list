package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "kafka:29092"
	}

	brokerList := strings.Split(brokers, ",")
	topic := "task-events"
	logFile := "/app/logs/events.log"

	// Ensure log directory exists
	if err := os.MkdirAll("/app/logs", 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Open log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	log.Printf("Kafka consumer starting, brokers: %v, topic: %s", brokerList, topic)

	// Configure consumer
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Create consumer
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Get partition consumer
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
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

	// Consume messages
	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				logEntry := string(msg.Value)
				log.Printf("Received event: %s", logEntry)

				// Write to file
				if _, err := file.WriteString(logEntry + "\n"); err != nil {
					log.Printf("Failed to write to log file: %v", err)
				}

			case err := <-partitionConsumer.Errors():
				log.Printf("Consumer error: %v", err)

			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for termination signal
	<-signals
	log.Println("Shutting down consumer...")
	cancel()
}
