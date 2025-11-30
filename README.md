# To-Do List Application

A microservices-based task management application built with Go, PostgreSQL, and Docker.

## Architecture

The application consists of three main services:

- **API Service** (port 8081): External HTTP API that handles client requests
- **DB Service** (port 8080): Internal service that manages database operations
- **PostgreSQL**: Database for storing tasks

## Features

- Create new tasks with name and description
- View all tasks or filter by completion status
- Mark tasks as complete
- Delete tasks
- Search tasks by ID or name
- RESTful API design

## Tech Stack

- **Go** - Backend services
- **PostgreSQL 15** - Database
- **Docker & Docker Compose** - Containerization
- **Gorilla Mux** - HTTP routing

## Project Structure

```
.
├── apiservice/          # External API service
│   ├── client/         # HTTP client for DB service
│   ├── handlersForDB/  # HTTP request handlers
│   └── models/         # Data models
├── db/                 # Database service
│   ├── handlers/       # HTTP request handlers
│   └── models/         # Data models and repository
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
- DB Service connects directly to PostgreSQL database
- PostgreSQL uses health checks to ensure it's ready before dependent services start

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
