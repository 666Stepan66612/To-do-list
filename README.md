# To-Do List Application

A microservices-based task management application built with Go, PostgreSQL, and Docker.

## Architecture

The application consists of five main services:

- **API Service** (port 8081): External HTTP API that handles client requests and sends events to Kafka
- **DB Service** (port 8080): Internal service that manages database operations
- **PostgreSQL**: Database for storing tasks
- **Kafka + Zookeeper**: Message broker for event logging
- **Kafka Service**: Consumer that logs all events to a file for audit and monitoring

## Features

- Create new tasks with name and description
- View all tasks or filter by completion status
- Mark tasks as complete
- Delete tasks
- Search tasks by ID or name
- RESTful API design
- **Event logging via Kafka** - All actions (create, update, delete, complete) are logged to `logs/events.log` for audit and monitoring

## Tech Stack

- **Go** - Backend services
- **PostgreSQL 15** - Database
- **Apache Kafka** - Event streaming and logging
- **Zookeeper** - Kafka coordination
- **Docker & Docker Compose** - Containerization
- **Gorilla Mux** - HTTP routing
- **Sarama** - Kafka client for Go

## Project Structure

```
.
├── apiservice/          # External API service
│   ├── client/         # HTTP client for DB service
│   ├── handlersForDB/  # HTTP request handlers
│   ├── kafka/          # Kafka producer for event logging
│   └── models/         # Data models
├── db/                 # Database service
│   ├── handlers/       # HTTP request handlers
│   └── models/         # Data models and repository
├── kafkaservice/       # Kafka consumer for event logging
│   └── main.go         # Consumes events and writes to log file
├── logs/               # Event logs (created at runtime)
│   └── events.log      # All task-related events
└── docker-compose.yaml # Service orchestration
```

## Data Model

```go
type Task struct {
    ID         int        `json:"id"`
    Name       string     `json:"name"`
    Text       string     `json:"text"`
    CreateTime time.Time  `json:"create_time"`
    Complete   bool       `json:"complete"`
    CompleteAt *time.Time `json:"complete_at"`
}
```

## API Endpoints

All endpoints are available at `http://localhost:8081`

### Create Task
```
POST /create
Content-Type: application/json

{
  "name": "Task name",
  "text": "Task description"
}
```

### Get All Tasks
```
GET /get
```

### Get Completed Tasks
```
GET /get?complete=true
```

### Get Uncompleted Tasks
```
GET /get?complete=false
```

### Get Task by ID
```
GET /getbyid/{id}
```

### Get Task by Name
```
GET /getbyname/{name}
```

### Complete Task
```
PUT /complete/{id}
```

### Delete Task
```
DELETE /delete/{id}
```

### Health Check
```
GET /health
```

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Running the Application

1. Clone the repository:
```bash
git clone https://github.com/666Stepan66612/To-do-list.git
cd To-do-list
```

2. Start all services:
```bash
docker-compose up --build
```

3. The API will be available at `http://localhost:8081`

### Stopping the Application

```bash
docker-compose down
```

To remove volumes as well:
```bash
docker-compose down -v
```

## Service Communication

- API Service communicates with DB Service via internal HTTP calls
- API Service sends all action events to Kafka topic `task-events`
- Kafka Service consumes events from Kafka and writes them to `logs/events.log`
- DB Service connects directly to PostgreSQL database
- PostgreSQL, Kafka, and Zookeeper use health checks to ensure readiness before dependent services start

## Event Logging

All task operations are logged to `logs/events.log` with the following format:

```json
{
  "timestamp": "2025-12-01T10:23:45Z",
  "action": "CREATE_TASK",
  "details": "Task created: id=1, name=Buy milk",
  "status": "SUCCESS"
}
```

Supported events: `CREATE_TASK`, `DELETE_TASK`, `COMPLETE_TASK`

## Environment Variables

### DB Service
- `WAIT_HOSTS=postgres:5432` - Wait for PostgreSQL to be ready

### API Service
- `WAIT_HOSTS=db-service:8080` - Wait for DB Service to be ready

## Database Connection

The DB service connects to PostgreSQL using:
```
postgres://postgres:mypostgres@postgres:5432/postgres?sslmode=disable
```

## License

This project is licensed under the terms specified in the LICENSE file.