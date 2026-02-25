# SentinelAI API Monitoring SaaS

A production-grade Go backend for the SentinelAI API monitoring platform.

## Current Features

- JWT Authentication
- Worker Pool Monitoring
- AI-based failure analysis using Ollama
- Configurable LLM provider
- Concurrency-safe worker pool
- Panic recovery
- Context-based graceful shutdown
- Structured Logging
- Pluggable Repositories (PostgreSQL & In-Memory options)
- Docker integration

## Architecture

This project strictly follows Clean Architecture principles to separate concerns efficiently.

**Architecture Flow Diagram:**
Auth → Monitor Service → Repository → WorkerPool → LLM → Repository Scheduler

### Module Structure

- **cmd/**: Application entrypoints. `main.go` acts solely as the orchestrator to bootstrap configuration and dependencies.
- **internal/**: Private application code strictly bounded within the module.
  - **handler/**: HTTP request handlers (presentation layer).
  - **service/**: Business logic layer.
  - **repository/**: Data access layer.
  - **auth/**: Independent authentication module (JWT, bcrypt).
  - **monitor/**: Monitoring engine (Scheduler, Worker pool, Handlers).
  - **llm/**: AI integration for failure analysis (`Ollama`).
  - **middleware/**: HTTP middlewares (logging, recovery, token validation, etc.).
  - **server/**: HTTP server setup and Dependency Injection container wiring.
- **pkg/**: Public libraries that can be used by other applications (config, logger).

## Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- Make
- [Docker] (Optional, for containerized DB/Service deployments)
- [Ollama](https://ollama.com/) (Optional, required for LLM analysis features)
- [PostgreSQL] (Optional, for persistent storage)

## Getting Started

1. Copy the example environment configuration:
   ```sh
   cp .env.example .env
   ```
   *Note: .env configures external DB access. Do not commit it to version control!*

2. Install dependencies:
   ```sh
   go mod tidy
   ```

3. Run the development server (Defaults to in-memory storage unless DB_HOST is populated):
   ```sh
   make run
   ```

### Docker Deployment

SentinelAI uses a multi-stage `Dockerfile` and `docker-compose.yml` to effortlessly spin up a unified containerized PostgreSQL 15 database running seamlessly alongside the latest Go backend binary.

Spin up the entire stack seamlessly:
```sh
docker-compose up -d --build
```

## Authentication Module

The system utilizes an independent authentication module leveraging bcrypt for secure password hashing and JWT (JSON Web Tokens) for stateless session management.

Once a user logs in, a signed JWT token is returned. This token must be provided in the `Authorization` header as a Bearer token (`Authorization: Bearer <token>`) for subsequent protected API requests. The included JWT middleware securely decodes, validates expiration, and extracts user claims.

### Auth API Endpoints

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

### Example Auth Usage

**Register a new user:**
```sh
curl -X POST "http://localhost:8080/api/v1/auth/register" \
 -H "Content-Type: application/json" \
 -d '{"email": "engineer@sentinelai.com", "password": "securepassword123"}'
```

**Login to receive JWT token:**
```sh
curl -X POST "http://localhost:8080/api/v1/auth/login" \
 -H "Content-Type: application/json" \
 -d '{"email": "engineer@sentinelai.com", "password": "securepassword123"}'
```

## Monitoring Engine Overview

The monitoring module handles the periodic execution of health checks against user-registered target URLs.

### Worker Pool Architecture

Operations are distributed via a fixed-size worker pool design using channel-based concurrency. Configurable goroutine workers are spawned on application startup to ingest HTTP health check jobs dynamically from a buffered job channel, enforcing resource predictability and isolating request latencies.

### Scheduler Explanation

A dedicated background scheduler runs concurrently on a configurable duration interval. Upon each tick, it loops through active monitors stored in the repository, validates intervals and active status constraints, and places viable jobs onto the unified worker pool channel.

### Concurrency Safety Design

- Context management governs the background loops for seamless runtime shutdowns.
- Internal state variables (`IsRunning`) lock specific health check targets to prevent redundant evaluations or race conditions.
- Deep-copy abstractions act as serialization safeguards over in-memory structs accessed during concurrent read/write operations spanning the worker pool, preventing panics and logic collision.
- Panic recovery mechanism encapsulates jobs bounding fault-tolerance explicitly to broken operations.

### JWT-Protected Monitor Endpoints

- `POST /api/v1/monitor/add`
- `GET /api/v1/monitor/list`

### Example Monitor Usage

**Add a new monitor (Requires JWT):**
```sh
curl -X POST "http://localhost:8080/api/v1/monitor/add" \
 -H "Authorization: Bearer <token>" \
 -H "Content-Type: application/json" \
 -d '{"url": "https://api.example.com/health", "interval": 60}'
```

**List all user monitors (Requires JWT):**
```sh
curl -X GET "http://localhost:8080/api/v1/monitor/list" \
 -H "Authorization: Bearer <token>"
```

## AI-Powered Failure Analysis

### Ollama Integration Overview

The backend incorporates local Large Language Model operations executing isolated anomaly explanation queries completely free from external cloud API dependencies. Through structured abstraction, the unified `llm.Provider` interface abstracts away direct bindings natively resolving HTTP context bridging safely.

### How LLM is triggered on monitor failure

When the background worker pool encounters an implicit connection termination (`err != nil`) or registers a failing HTTP boundary (`status >= 400`), the worker halts positive assertion workflows and transmits the exact request metrics (Timestamp, Latency, Status Code) into an isolated `AnalyzeFailure` execution sandbox. It bounds this contextually up to 10 seconds locally to preserve core concurrency pools without risk of starvation deadlocks bridging the returned text strings safely inside the generic repository mapping.

### Example Monitor Failure with AI Explanation

When hitting the `/list` endpoint following a timeout event on a registered URL, the AI Explanation output mimics this natively:
```json
{
  "success": true,
  "message": "monitors retrieved",
  "data": [
    {
      "id": "20231105000000000",
      "user_id": "user-uuid",
      "url": "https://unreachable.application.local",
      "interval": 60000000000,
      "last_checked": "2023-11-05T12:00:00Z",
      "status_code": 0,
      "response_time": 10000000000,
      "is_healthy": false,
      "is_running": false,
      "ai_explanation": "The failure strongly points to a domain resolution gap or a severe localized routing disruption given the zeroed status code and exact 10-second timeout ceiling termination."
    }
  ]
}
```

## API Response Format

All REST endpoints strictly adhere to the following standard JSON response structure, ensuring explicit application context for client integrations:

```json
{
  "success": true,
  "message": "string",
  "data": {}
}
```

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `PORT` | Port the HTTP server binds to. | `8080` |
| `ENV` | Environment state (`development`, `production`). | `development` |
| `JWT_SECRET` | Secret key used for cryptographic JWT signing and verification. | `super-secret-local-dev-key` |
| `TOKEN_EXPIRATION` | Number of hours before the issued JWT token expires. | `24` |
| `SCHEDULER_INTERVAL` | Interval logic tick evaluation loop duration in seconds. | `1` |
| `OLLAMA_URL` | Local LLM host URL mapping. | `http://localhost:11434/api/generate` |
| `LLM_MODEL` | Machine learning model invoked for failure parsing. | `llama3` |
| `DB_HOST` | PostgreSQL Hostname (triggers postgres injection). | `postgres` / `localhost` |
| `DB_PORT` | PostgreSQL connection port. | `5432` |
| `DB_USER` | PostgreSQL active role context. | `postgres` |
| `DB_PASSWORD` | PostgreSQL active role password. | `postgres` |
| `DB_NAME` | Initialized backend storage database namespace. | `sentinel` |

## CI/CD Workflow

The repository is natively integrated with GitHub Actions. On every push and pull request to `main`, the CI pipeline automatically:
- Validates formats (`gofmt`)
- Builds the application binary
- Runs unit tests (`go test`)
- Executes strict static analysis (`golangci-lint`)
