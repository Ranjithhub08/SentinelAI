# SentinelAI API Monitoring SaaS

This is the core backend skeleton for the SentinelAI API Monitoring SaaS platform.

## Architecture

This project strictly follows Clean Architecture principles:
- **cmd/**: Application entrypoints. `main.go` only acts as the orchestrator to bootstrap configuration and dependencies.
- **internal/**: Private application code strictly bounded within the module.
  - **handler/**: HTTP request handlers (presentation layer).
  - **service/**: Business logic layer.
  - **repository/**: Data access layer.
  - **auth/**: Independent authentication module (JWT, bcrypt).
  - **server/**: HTTP server and dependency injection setup.
  - **middleware/**: HTTP middlewares (logging, recovery, etc.).
- **pkg/**: Public libraries that can be used by other applications (config, logger).

## Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- Make

## Getting Started

1. **Set up the environment:**
   Copy the example environment configuration to establish your local `.env`.
   ```sh
   cp .env.example .env
   ```
   *(Note: The `.env` file is excluded from version control to protect sensitive configurations/secrets).*

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Run the server locally:**
   ```sh
   make run
   ```

4. **Build the binary:**
   ```sh
   make build
   ```

5. **Run tests:**
   ```sh
   make test
   ```

## Example API Usage

**1. Register a new user:**
```sh
curl -X POST "http://localhost:8080/api/v1/auth/register" \
 -H "Content-Type: application/json" \
 -d '{"email": "engineer@sentinelai.com", "password": "securepassword123"}'
```

**2. Login to receive JWT token:**
```sh
curl -X POST "http://localhost:8080/api/v1/auth/login" \
 -H "Content-Type: application/json" \
 -d '{"email": "engineer@sentinelai.com", "password": "securepassword123"}'
```

## Development & CI/CD Workflow

The repository is fully configured with a GitHub Actions workflow (`.github/workflows/ci.yml`) ensuring robust code quality. On every push and pull request to `main`, the CI pipeline automatically:
- Sets up the Go environment
- Installs dependencies and verifies formats (`gofmt`)
- Builds the application (`make build`)
- Runs your test suites (`make test`)
- Enforces strict code quality scanning with `golangci-lint`

No API Keys, tokens, or environment files (`.env`) should ever be committed to git.
