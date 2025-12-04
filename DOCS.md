# To-Do List Application

A full-stack microservices-based task management application with authentication, built with Go, PostgreSQL, Kafka, and Docker.

## Architecture

The application consists of seven main services:

- **Frontend** (port 8080): Modern web UI with authentication, built with HTML/CSS/JavaScript served via nginx
- **API Service** (port 8081): External HTTP API with JWT authentication that handles client requests and sends events to Kafka
- **DB Service** (port 8080 internal): Internal service that manages database operations
- **PostgreSQL**: Database for storing users and tasks
- **Kafka + Zookeeper**: Message broker for event logging and audit trail
- **Kafka Service**: Consumer that logs all events to a file for monitoring and audit

## Features

### Authentication
- **User Registration** - Create account with username (min 3 chars) and password (min 8 chars)
- **User Login** - Secure authentication with JWT tokens
- **Password Security** - Bcrypt hashing with cost factor 12
- **Session Management** - Token stored in localStorage, automatic logout on expiration
- **Protected Routes** - All task operations require valid authentication

### Frontend
- Modern, responsive web interface with auth screens
- User registration and login forms
- Create tasks with name and optional description
- View task details in modal popup
- Filter tasks: All / Active / Completed
- Mark tasks as complete (visual green highlight)
- Delete tasks with confirmation
- Real-time updates
- Gradient purple theme with smooth animations
- User-specific task display

### Backend
- RESTful API design with JWT middleware
- User registration and login with bcrypt password hashing
- Create tasks with name and text description (user-specific)
- View all tasks or filter by completion status (user-specific)
- Mark tasks as complete (prevents duplicate completions)
- Delete tasks (user-specific)
- Search tasks by ID or name
- **Event logging via Kafka** - All actions (create, delete, complete) with user information are logged to `logs/events.log` for audit trail
- CORS support for cross-origin requests

## Tech Stack

### Frontend
- **HTML5 / CSS3 / JavaScript (ES6+)** - Modern web interface
- **Nginx Alpine** - Static file server

### Backend
- **Go 1.25** - Backend microservices
- **PostgreSQL 15** - Relational database with users and tasks tables
- **Bcrypt** - Secure password hashing with cost factor 12
- **JWT** - JSON Web Tokens for authentication (24-hour expiry)
- **Apache Kafka 7.5** - Event streaming and logging
- **Zookeeper 7.5** - Kafka coordination
- **Docker & Docker Compose** - Containerization and orchestration
- **Gorilla Mux** - HTTP routing
- **Sarama** - Kafka client for Go

## Project Structure

```
.
├── frontend/           # Web UI (nginx)
│   ├── index.html     # Single-page application with auth
│   └── Dockerfile     # Nginx container
├── apiservice/        # External API service
│   ├── auth/          # JWT and bcrypt utilities
│   ├── client/        # HTTP client for DB service
│   ├── handlers/      # HTTP request handlers
│   │   ├── authhandlers.go  # Registration and login
│   │   └── handlers.go      # Task operations with event logging
│   ├── kafka/         # Kafka producer for event logging
│   ├── middleware/    # JWT authentication middleware
│   ├── models/        # Data models
│   └── main.go        # API server with CORS support
├── db/                # Database service
│   ├── handlers/      # HTTP request handlers
│   │   ├── authhandlers.go  # User creation and retrieval
│   │   └── handlers.go      # Task CRUD operations
│   ├── models/        # Data models and repository
│   └── main.go        # DB service server with migrations
├── kafkaservice/      # Kafka consumer for event logging
│   └── main.go        # Consumes events and writes to log file
├── logs/              # Event logs (bind-mounted to host)
│   └── events.log     # All task-related events with user info (JSON format)
└── docker-compose.yaml # Service orchestration
```

## Data Model

### User
```go
type User struct {
    ID           int       `json:"id"`
    Username     string    `json:"username"`
    PasswordHash string    `json:"password_hash,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
}
```

### Task
```go
type Task struct {
    ID         int        `json:"id"`
    UserID     int        `json:"user_id"`
    Name       string     `json:"name"`
    Text       string     `json:"text"`
    CreateTime time.Time  `json:"create_time"`
    Complete   bool       `json:"complete"`
    CompleteAt *time.Time `json:"complete_at"`
}
```

### JWT Claims
```go
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

## API Endpoints

All endpoints are available at `http://localhost:8081`

### Authentication Endpoints

#### Register
```http
POST /register
Content-Type: application/json

{
  "username": "user123",
  "password": "password123"
}

Response: 201 Created
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "user123",
  "user_id": 1
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "username": "user123",
  "password": "password123"
}

Response: 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "user123",
  "user_id": 1
}
```

### Task Endpoints (Require Authentication)

**All task endpoints require the `Authorization` header:**
```http
Authorization: Bearer <jwt_token>
```

#### Create Task
```http
POST /create
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "Task name",
  "text": "Task description (optional)"
}
```

#### Get All Tasks (User-Specific)
```http
GET /tasks
Authorization: Bearer <jwt_token>
```

#### Complete Task
```http
POST /complete/{id}
Authorization: Bearer <jwt_token>
```
**Note:** Prevents duplicate completions - returns error if task already completed.

#### Delete Task
```http
DELETE /delete/{id}
Authorization: Bearer <jwt_token>
```

#### Health Check
```http
GET /health
```

**CORS:** All endpoints support CORS with `Access-Control-Allow-Origin: *` for frontend integration.

## Security Features

- **Password Hashing**: Bcrypt with cost factor 12 (~400ms per hash)
- **JWT Authentication**: 24-hour token expiry, signed with HS256
- **User Isolation**: Each user sees only their own tasks
- **Audit Trail**: All user actions logged with user_id and username
- **SQL Injection Protection**: Parameterized queries throughout
- **HTTPS Ready**: Works with reverse proxy for SSL termination

⚠️ **Production Security Notes:**
- Change default JWT secret key
- Change default database password
- Use environment variables for sensitive data
- Enable HTTPS in production
- Regularly rotate JWT secret
- Monitor event logs for suspicious activity

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

3. Access the application:
   - **Web UI**: http://localhost:8080
   - **API**: http://localhost:8081
   - **Event Logs**: `./logs/events.log`

### Stopping the Application

```bash
docker-compose down
```

To remove volumes as well:
```bash
docker-compose down -v
```

## Service Communication

```
┌──────────┐     HTTP      ┌─────────────┐     HTTP      ┌────────────┐
│ Frontend ├──────────────►│ API Service ├──────────────►│ DB Service │
│ (nginx)  │  localhost:   │   (Go)      │   internal    │   (Go)     │
│  Auth +  │     8081      │  JWT Auth + │               │ PostgreSQL │
│  Tasks   │               │   CORS      │               │   Access   │
└──────────┘               └──────┬──────┘               └─────┬──────┘
                                  │                            │
                                  │ Kafka Events               │ SQL
                                  │ (with user_id)             │
                           ┌──────▼──────┐              ┌─────▼──────┐
                           │    Kafka    │              │ PostgreSQL │
                           │ (topic:     │              │  - users   │
                           │ task-events)│              │  - tasks   │
                           └──────┬──────┘              └────────────┘
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
                           (with user info)
```

- **Frontend** serves static HTML/CSS/JS via nginx and communicates with API Service
- **API Service** handles authentication (register/login) and validates JWT tokens for all task operations
- **API Service** communicates with DB Service via internal HTTP calls
- **API Service** sends all action events with user information to Kafka topic `task-events`
- **Kafka Service** consumes events from Kafka and writes them to `logs/events.log`
- **DB Service** manages PostgreSQL database with users and tasks tables
- **Health checks** ensure PostgreSQL, Kafka, and Zookeeper are ready before dependent services start
- **Logs volume** is bind-mounted to host filesystem for persistence

## Event Logging

All task operations are logged to `logs/events.log` with user information for audit purposes:

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

**Logged Events:**
- `CREATE_TASK` - Task creation with task ID and name
- `DELETE_TASK` - Task deletion with task ID
- `COMPLETE_TASK` - Task completion with task ID

**Event Status:**
- `SUCCESS` - Operation completed successfully
- `ERROR` - Operation failed (with error details)

This provides a complete audit trail of who performed which actions and when.

## Environment Variables

### API Service
- `WAIT_HOSTS=db-service:8080` - Wait for DB Service to be ready
- `KAFKA_BROKERS=kafka:29092` - Kafka broker address for event logging
- JWT Secret: Configured in `apiservice/auth/auth.go` (⚠️ change in production!)

### DB Service
- `DB_HOST=postgres` - PostgreSQL host
- `DB_USER=postgres` - PostgreSQL user
- `DB_PASSWORD=mypostgres` - PostgreSQL password (⚠️ change in production!)
- `DB_NAME=postgres` - PostgreSQL database name
- `WAIT_HOSTS=postgres:5432` - Wait for PostgreSQL to be ready

### Kafka Service
- `KAFKA_BROKERS=kafka:29092` - Kafka broker address for consuming events

## Database Schema

The DB service automatically creates and manages the following tables:

### `users` table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### `tasks` table
```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    text TEXT,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    complete BOOLEAN DEFAULT FALSE,
    complete_at TIMESTAMP
);

CREATE INDEX idx_tasks_user_id ON tasks(user_id);
```

Data is persisted in Docker volume `todo_postgres_data`.

## Deployment

### Local Development
The application runs on `localhost` by default. Each user has their own account and sees only their own tasks.

### Production Deployment Checklist

⚠️ **Security Updates Required:**

1. **Change JWT Secret** in `apiservice/auth/auth.go`:
   ```go
   var jwtSecret = []byte("your-production-secret-here")
   ```

2. **Change Database Password** in `docker-compose.yaml`:
   ```yaml
   POSTGRES_PASSWORD: your-secure-password
   DB_PASSWORD: your-secure-password
   ```

3. **Update Frontend API URL** in `frontend/index.html`:
   ```javascript
   const API_URL = 'https://your-domain.com';
   ```

4. **Configure HTTPS** with reverse proxy (nginx/traefik) for production

5. **Set up monitoring** for event logs and container health

### Production Deployment Options

#### Option 1: Oracle Cloud (Free Tier)
Oracle Cloud offers **Always Free** VMs with:
- 1 VM with 4 CPU + 24 GB RAM (ARM Ampere)
- 200 GB storage
- Public IP address

Setup:
1. Create Oracle Cloud account
2. Launch Ubuntu ARM instance
3. Install Docker and Docker Compose
4. Clone repository and update configuration
5. Configure firewall to allow ports 80, 443
6. Set up SSL with Let's Encrypt

#### Option 2: VPS Hosting
Deploy to any VPS provider (DigitalOcean, Hetzner, Linode):
1. Install Docker and Docker Compose on VPS
2. Clone repository
3. Update configuration (see checklist above)
4. Run `docker-compose up -d`
5. Configure firewall and SSL

#### Option 3: Ngrok (Quick Testing)
For temporary public access:
```bash
ngrok http 8080
```

## License

This project is licensed under the terms specified in the LICENSE file.