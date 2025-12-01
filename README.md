# To-Do List Application

A full-stack microservices-based task management application with a modern web interface, built with Go, PostgreSQL, Kafka, and Docker.

## Architecture

The application consists of six main services:

- **Frontend** (port 8080): Modern web UI built with HTML/CSS/JavaScript served via nginx
- **API Service** (port 8081): External HTTP API that handles client requests and sends events to Kafka
- **DB Service** (port 8080 internal): Internal service that manages database operations
- **PostgreSQL**: Database for storing tasks
- **Kafka + Zookeeper**: Message broker for event logging and audit trail
- **Kafka Service**: Consumer that logs all events to a file for monitoring and audit

## Features

### Frontend
- Modern, responsive web interface
- Create tasks with name and optional description
- View task details in modal popup
- Filter tasks: All / Active / Completed
- Mark tasks as complete (visual green highlight)
- Delete tasks with confirmation
- Real-time updates
- Gradient purple theme with smooth animations

### Backend
- RESTful API design
- Create tasks with name and text description
- View all tasks or filter by completion status
- Mark tasks as complete (prevents duplicate completions)
- Delete tasks
- Search tasks by ID or name
- **Event logging via Kafka** - All actions (create, delete, complete) are logged to `logs/events.log` for audit trail
- CORS support for cross-origin requests

## Tech Stack

### Frontend
- **HTML5 / CSS3 / JavaScript (ES6+)** - Modern web interface
- **Nginx Alpine** - Static file server

### Backend
- **Go 1.25** - Backend microservices
- **PostgreSQL 15** - Relational database
- **Apache Kafka 7.5** - Event streaming and logging
- **Zookeeper 7.5** - Kafka coordination
- **Docker & Docker Compose** - Containerization and orchestration
- **Gorilla Mux** - HTTP routing
- **Sarama** - Kafka client for Go

## Project Structure

```
.
├── frontend/           # Web UI (nginx)
│   ├── index.html     # Single-page application
│   └── Dockerfile     # Nginx container
├── apiservice/        # External API service
│   ├── client/        # HTTP client for DB service
│   ├── handlersForDB/ # HTTP request handlers with event emission
│   ├── kafka/         # Kafka producer for event logging
│   ├── models/        # Data models
│   └── main.go        # API server with CORS support
├── db/                # Database service
│   ├── handlers/      # HTTP request handlers
│   ├── models/        # Data models and repository
│   └── main.go        # DB service server
├── kafkaservice/      # Kafka consumer for event logging
│   └── main.go        # Consumes events and writes to log file
├── logs/              # Event logs (bind-mounted to host)
│   └── events.log     # All task-related events (JSON format)
├── docker-compose.yaml # Service orchestration
├── Makefile           # Convenient shortcuts for docker-compose
└── README.md          # This file
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
```http
POST /create
Content-Type: application/json

{
  "name": "Task name",
  "text": "Task description (optional)"
}
```

### Get All Tasks
```http
GET /tasks
GET /get
```

### Get Completed Tasks
```http
GET /get?complete=true
```

### Get Uncompleted Tasks
```http
GET /get?complete=false
```

### Get Task by ID
```http
GET /getbyid/{id}
```

### Get Task by Name
```http
GET /getbyname/{name}
```

### Complete Task
```http
POST /complete/{id}
PUT /complete/{id}
```
**Note:** Prevents duplicate completions - returns error if task already completed.

### Delete Task
```http
DELETE /delete/{id}
```

### Health Check
```http
GET /health
```

**CORS:** All endpoints support CORS with `Access-Control-Allow-Origin: *` for frontend integration.

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

Or use the Makefile:
```bash
make up
```

3. Access the application:
   - **Web UI**: http://localhost:8080
   - **API**: http://localhost:8081
   - **Event Logs**: `./logs/events.log`

### Makefile Commands

```bash
make up          # Start all services
make down        # Stop all services
make restart     # Restart all services
make logs        # View logs from all services
make logs-api    # View API service logs
make logs-db     # View DB service logs
make logs-kafka  # View Kafka consumer logs
make ps          # Show running containers
make clean       # Stop and remove all containers, networks, volumes
make rebuild     # Rebuild and restart all services
```

### Stopping the Application

```bash
docker-compose down
# or
make down
```

To remove volumes as well:
```bash
docker-compose down -v
# or
make clean
```

## Service Communication

```
┌──────────┐     HTTP      ┌─────────────┐     HTTP      ┌────────────┐
│ Frontend ├──────────────►│ API Service ├──────────────►│ DB Service │
│ (nginx)  │  localhost:   │   (Go)      │   internal    │   (Go)     │
│          │     8081      │             │               │            │
└──────────┘               └──────┬──────┘               └─────┬──────┘
                                  │                            │
                                  │ Kafka Events               │ SQL
                                  │                            │
                           ┌──────▼──────┐              ┌─────▼──────┐
                           │    Kafka    │              │ PostgreSQL │
                           │ (topic:     │              │            │
                           │ task-events)│              └────────────┘
                           └──────┬──────┘
                                  │
                                  │ Consume
                                  │
                           ┌──────▼────────┐
                           │ Kafka Service │
                           │   (Go)        │
                           │  Write logs   │
                           └───────────────┘
                                  │
                                  ▼
                           logs/events.log
```

- **Frontend** serves static HTML/CSS/JS via nginx and communicates with API Service
- **API Service** communicates with DB Service via internal HTTP calls
- **API Service** sends all action events to Kafka topic `task-events`
- **Kafka Service** consumes events from Kafka and writes them to `logs/events.log`
- **DB Service** connects directly to PostgreSQL database
- **Health checks** ensure PostgreSQL, Kafka, and Zookeeper are ready before dependent services start
- **Logs volume** is bind-mounted to host filesystem for persistence

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

### API Service
- `WAIT_HOSTS=db-service:8080` - Wait for DB Service to be ready
- `KAFKA_BROKERS=kafka:29092` - Kafka broker address for event logging

### DB Service
- `WAIT_HOSTS=postgres:5432` - Wait for PostgreSQL to be ready

### Kafka Service
- `KAFKA_BROKERS=kafka:29092` - Kafka broker address for consuming events

## Database Connection

The DB service connects to PostgreSQL using:
```
postgres://postgres:mypostgres@postgres:5432/postgres?sslmode=disable
```

Data is persisted in Docker volume `todo_postgres_data`.

## Deployment

### Local Development
The application runs on `localhost` by default.

### Production Deployment Options

#### Option 1: Oracle Cloud (Free Tier)
Oracle Cloud offers **Always Free** VMs with:
- 1 VM with 4 CPU + 24 GB RAM (ARM Ampere)
- 200 GB storage
- Public IP address

See deployment guide in the repository wiki.

#### Option 2: VPS Hosting
Deploy to any VPS provider (DigitalOcean, Hetzner, etc.):
1. Install Docker and Docker Compose on VPS
2. Clone repository
3. Update `frontend/index.html` with your VPS IP:
   ```javascript
   const API_URL = 'http://your-vps-ip:8081';
   ```
4. Run `docker-compose up -d`
5. Configure firewall to allow ports 8080, 8081

#### Option 3: Ngrok (Quick Testing)
For temporary public access:
```bash
ngrok http 8080
```

**Note:** Current configuration uses a shared database - all users see the same tasks. For multi-user support, implement authentication and user-specific task filtering.

## License

This project is licensed under the terms specified in the LICENSE file.