# SentinelAI API Monitoring SaaS

A production-grade Go backend for the SentinelAI platform.

## Architecture

This project strictly follows Clean Architecture principles:

- **cmd/**: Application entrypoints. `main.go` acts solely as the orchestrator to bootstrap configuration and dependencies.
- **internal/**: Private application code strictly bounded within the module.
  - **handler/**: HTTP request handlers (presentation layer).
  - **service/**: Business logic layer.
  - **repository/**: Data access layer.
  - **auth/**: Independent authentication module (JWT, bcrypt).
  - **middleware/**: HTTP middlewares (logging, recovery, token validation, etc.).
  - **server/**: HTTP server setup and Dependency Injection container wiring.
- **pkg/**: Public libraries that can be used by other applications (config, logger).

## Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- Make

## Getting Started

1. Copy the example environment configuration:
   ```sh
   cp .env.example .env
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

3. Run the development server:
   ```sh
   make run
   ```

## Authentication Module

The system utilizes an independent authentication module leveraging **bcrypt** for secure password hashing and **JWT (JSON Web Tokens)** for stateless session management.

Once a user logs in, a signed JWT token is returned. This token must be provided in the `Authorization` header as a Bearer token (`Authorization: Bearer <token>`) for subsequent protected API requests. The included JWT middleware securely decodes, validates expiration, and extracts user claims.

### API Endpoints

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

### Example API Usage

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

| Variable | Description | Default / Example |
|---|---|---|
| `PORT` | Port the HTTP server binds to. | `8080` |
| `ENV` | Environment state (`development`, `production`). | `development` |
| `JWT_SECRET` | Secret key used for cryptographic JWT signing and verification. | `super-secret-local-dev-key` |
| `TOKEN_EXPIRATION` | Number of hours before the issued JWT token expires. | `24` |

## CI/CD Workflow

The repository is natively integrated with GitHub Actions. On every push and pull request to `main`, the CI pipeline automatically:
- Validates formats (`gofmt`)
- Builds the application binary
- Runs unit tests (`go test`)
- Executes strict static analysis (`golangci-lint`)
