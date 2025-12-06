# To-Do List Application

A microservices-based task management application with JWT authentication, event logging via Kafka, and modern web interface.

## Quick Start

### Using Make (Recommended)

```bash
# Clone repository
git clone https://github.com/666Stepan66612/To-do-list.git
cd To-do-list

# View all available commands
make help

# Start all services
make up

# View logs
make logs

# Access the app
# Web UI: http://localhost:8080
# API: http://localhost:8081
```

### Using Docker Compose

```bash
# Clone repository
git clone https://github.com/666Stepan66612/To-do-list.git
cd To-do-list

# Start all services
docker-compose up --build

# Access the app
# Web UI: http://localhost:8080
# API: http://localhost:8081
```

## Features

- **User Authentication** - JWT-based auth with bcrypt password hashing
- **Task Management** - Create, complete, delete tasks with descriptions
- **User Isolation** - Each user sees only their own tasks
- **Event Logging** - All actions logged to Kafka with user information
- **Modern UI** - Responsive web interface with auth screens
- **Task Filtering** - View all, active, or completed tasks

## Architecture

**7 Microservices:**
- **Frontend** (nginx) - Web UI on port 8080
- **API Service** (Go) - REST API with JWT auth on port 8081
- **DB Service** (Go) - Database operations (internal)
- **PostgreSQL** - User and task storage
- **Kafka + Zookeeper** - Event streaming
- **Kafka Service** (Go) - Event consumer and logger

## Tech Stack

- **Backend:** Go 1.25, PostgreSQL 15, Kafka 7.5, JWT, Bcrypt
- **Frontend:** HTML5/CSS3/JavaScript, Nginx
- **Infrastructure:** Docker, Docker Compose, Make

## API Endpoints

### Authentication
```http
POST /register   # Register new user
POST /login      # Login user
```

### Tasks (Require JWT Token)
```http
GET    /tasks         # Get all user tasks
POST   /create        # Create new task
POST   /complete/{id} # Mark task complete
DELETE /delete/{id}   # Delete task
```

## Security Features

- Bcrypt password hashing (cost factor 12)
- JWT tokens with 24-hour expiry
- User-specific task isolation
- Complete audit trail with user information
- SQL injection protection via parameterized queries

## Event Logging

All operations logged to `logs/events.log`:
```json
{
  "timestamp": "2025-12-04T13:35:31Z",
  "user_id": 3,
  "username": "MrKhm",
  "action": "DELETE_TASK",
  "details": "Task deleted: id=3",
  "status": "SUCCESS"
}
```

## Development Commands

### Make Commands

| Command | Description |
|---------|-------------|
| `make up` | Start all services |
| `make down` | Stop all services |
| `make logs` | View logs from all services |
| `make restart` | Restart all services |
| `make test` | Run all tests |
| `make coverage` | Generate coverage reports |
| `make clean` | Clean up containers and volumes |
| `make dev` | Development mode with logs |
| `make health` | Check service health |

Run `make help` to see all available commands.

### WSL Setup (Windows)

1. Enable Docker Desktop WSL integration:
   - Settings → Resources → WSL Integration
   - Enable your WSL distribution
   - Apply & Restart

2. Navigate to project in WSL:
```bash
cd /mnt/c/Users/YourUsername/path/to/toDo
make up
```

## Production Notes

Before deploying to production:
1. Change JWT secret in `apiservice/auth/auth.go`
2. Change database password in `docker-compose.yaml`
3. Update API URL in `frontend/index.html`
4. Enable HTTPS with reverse proxy
5. Set up monitoring for event logs

## Full Documentation

See [DOCS.md](DOCS.md) for complete documentation including:
- Detailed architecture
- API specifications
- Database schema
- Deployment guides
- Environment variables

## License

This project is licensed under the terms specified in the LICENSE file.