package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type EventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

type Event struct {
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	Status    string `json:"status"`
}

func NewEventProducer(brokers []string, topic string) (*EventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &EventProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (ep *EventProducer) SendEvent(action, details, status string) error {
	if ep == nil {
		// Silently skip if producer is not initialized
		return nil
	}

	event := Event{
		Timestamp: time.Now().Format(time.RFC3339),
		Action:    action,
		Details:   details,
		Status:    status,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: ep.topic,
		Value: sarama.StringEncoder(eventJSON),
	}

	_, _, err = ep.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send event to Kafka: %v", err)
		return err
	}

	log.Printf("Event sent to Kafka: %s", action)
	return nil
}

func (ep *EventProducer) Close() error {
	return ep.producer.Close()
}
