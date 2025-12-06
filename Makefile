.PHONY: help build up down restart logs clean test test-api test-db test-kafka test-all coverage ps stop start

# Default target
help:
	@echo "Available targets:"
	@echo "  make build          - Build all Docker images"
	@echo "  make up             - Start all services"
	@echo "  make down           - Stop and remove all containers"
	@echo "  make restart        - Restart all services"
	@echo "  make logs           - Show logs from all services"
	@echo "  make logs-api       - Show API service logs"
	@echo "  make logs-db        - Show DB service logs"
	@echo "  make logs-kafka     - Show Kafka service logs"
	@echo "  make clean          - Remove containers, volumes, and images"
	@echo "  make test           - Run tests in all services"
	@echo "  make test-api       - Run API service tests"
	@echo "  make test-db        - Run DB service tests"
	@echo "  make test-kafka     - Run Kafka service tests"
	@echo "  make coverage       - Run tests with coverage for all services"
	@echo "  make coverage-api   - Generate coverage report for API service"
	@echo "  make coverage-db    - Generate coverage report for DB service"
	@echo "  make ps             - List running containers"
	@echo "  make stop           - Stop all services"
	@echo "  make start          - Start existing containers"

# Build all Docker images
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Start all services with logs
up-logs:
	docker-compose up

# Stop and remove all containers
down:
	docker-compose down

# Restart all services
restart: down up

# Show logs from all services
logs:
	docker-compose logs -f

# Show logs from specific services
logs-api:
	docker-compose logs -f api-service

logs-db:
	docker-compose logs -f db-service

logs-kafka:
	docker-compose logs -f kafka-service

logs-frontend:
	docker-compose logs -f frontend

# List running containers
ps:
	docker-compose ps

# Stop services without removing
stop:
	docker-compose stop

# Start existing containers
start:
	docker-compose start

# Clean up everything
clean:
	docker-compose down -v --rmi all --remove-orphans

# Clean volumes only
clean-volumes:
	docker-compose down -v

# Run tests for API service
test-api:
	cd apiservice && go test ./... -v

# Run tests for DB service
test-db:
	cd db && go test ./... -v

# Run tests for Kafka service
test-kafka:
	cd kafkaservice && go test ./... -v

# Run all tests
test-all: test-api test-db test-kafka

# Alias for test-all
test: test-all

# Generate coverage for API service
coverage-api:
	cd apiservice && go test ./... -coverprofile=coverage.out
	cd apiservice && go tool cover -html=coverage.out -o coverage.html
	@echo "API Coverage report generated: apiservice/coverage.html"

# Generate coverage for DB service
coverage-db:
	cd db && go test ./... -coverprofile=coverage.out
	cd db && go tool cover -html=coverage.out -o coverage.html
	@echo "DB Coverage report generated: db/coverage.html"

# Generate coverage for Kafka service
coverage-kafka:
	cd kafkaservice && go test ./... -coverprofile=coverage.out
	cd kafkaservice && go tool cover -html=coverage.out -o coverage.html
	@echo "Kafka Coverage report generated: kafkaservice/coverage.html"

# Generate coverage for all services
coverage: coverage-api coverage-db coverage-kafka

# Install Go dependencies
deps:
	cd apiservice && go mod download
	cd db && go mod download
	cd kafkaservice && go mod download

# Update Go dependencies
deps-update:
	cd apiservice && go get -u ./...
	cd db && go get -u ./...
	cd kafkaservice && go get -u ./...

# Format Go code
fmt:
	cd apiservice && go fmt ./...
	cd db && go fmt ./...
	cd kafkaservice && go fmt ./...

# Lint Go code (requires golangci-lint)
lint:
	cd apiservice && golangci-lint run
	cd db && golangci-lint run
	cd kafkaservice && golangci-lint run

# Create database backup
backup-db:
	docker-compose exec postgres pg_dump -U postgres postgres > backup_$$(date +%Y%m%d_%H%M%S).sql

# Restore database from backup (usage: make restore-db FILE=backup.sql)
restore-db:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make restore-db FILE=backup.sql"; \
		exit 1; \
	fi
	docker-compose exec -T postgres psql -U postgres postgres < $(FILE)

# Show service health status
health:
	@echo "=== Service Health Status ==="
	@docker-compose ps
	@echo "\n=== Kafka Topics ==="
	@docker-compose exec kafka kafka-topics --bootstrap-server localhost:9092 --list || echo "Kafka not ready"

# Development mode - rebuild and restart with logs
dev: down build up-logs

# Production mode - start without build
prod:
	docker-compose up -d

# Quick restart for specific service
restart-api:
	docker-compose restart api-service

restart-db:
	docker-compose restart db-service

restart-kafka:
	docker-compose restart kafka-service

restart-frontend:
	docker-compose restart frontend
