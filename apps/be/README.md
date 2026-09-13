# Ticket Box Backend (Go + Gin)

Go backend for the **Ticket Box** system, built with [Gin Web Framework](https://github.com/gin-gonic/gin).

---

## 📁 Project Structure

```
apps/be/
├── cmd/
│   ├── api/
│   │   └── main.go                 # Application entry point (dependency injection, graceful shutdown)
│   └── migrate/
│       └── main.go                 # Database auto-migration CLI
├── configs/
│   └── .env.example                # Sample environment variables
├── docs/                           # Auto-generated OpenAPI / Swagger specs (docs.go, swagger.json, swagger.yaml)
├── internal/                       # Private application code (Go compiler protected)
│   ├── config/                     # Configuration and environment variable loader
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain/                     # Core domain entities, errors, and interface definitions
│   │   ├── errors.go
│   │   ├── ticket.go
│   │   └── ticket_test.go
│   ├── handler/                    # Transport layer (HTTP / Gin handlers)
│   │   └── http/
│   │       ├── health_handler.go
│   │       ├── health_handler_test.go
│   │       ├── router.go
│   │       └── router_test.go
│   ├── middleware/                 # Gin HTTP middlewares (CORS, Logger, Recovery)
│   │   ├── cors.go
│   │   └── cors_test.go
│   ├── pkg/                        # Internal shared libraries and helpers
│   │   ├── database/               # Database connection and AutoMigrate runner
│   │   ├── response/               # Standardized JSON response envelope
│   │   └── testutil/               # Test isolation DB helpers (transaction rollback)
│   ├── repository/                 # Data access layer (PostgreSQL GORM implementations)
│   │   └── postgres/
│   │       ├── ticket_repo.go
│   │       └── ticket_repo_test.go
│   └── service/                    # Business logic and use-cases layer
├── .golangci.yml                   # Linter configuration (golangci-lint)
├── Makefile                        # Common developer task automation
├── go.mod                          # Go module definitions
├── go.sum                          # Go module checksums
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites
- **Go**: Version `1.22+` (tested on `1.25.x`)
- **PostgreSQL**: Version `16+` (for local development and integration tests)
- **Make**: (For running Makefile automation)

### Installation & Setup

1. Clone the repository and navigate to the backend directory:
   ```bash
   cd apps/be
   ```

2. Run automated setup to install dependencies, developer CLI tools (`swag`, `golangci-lint`), and generate docs:
   ```bash
   make setup
   ```

3. (Optional) Set up environment variables:
   ```bash
   cp .env.example .env
   ```

---

## 🛠️ Development & Commands

The project includes a `Makefile` with common tasks:

| Command              | Description                                                          |
|----------------------|----------------------------------------------------------------------|
| `make setup`         | Installs developer tools (`swag`, `golangci-lint`) and dependencies  |
| `make run`           | Starts the API server locally                                        |
| `make build`         | Compiles the binary to `bin/api`                                     |
| `make migrate`       | Runs database migrations (`cmd/migrate`)                             |
| `make swagger`       | Generates OpenAPI / Swagger docs into `docs/` using `swag`           |
| `make test`          | Runs unit and integration tests with race condition detection        |
| `make test-coverage` | Runs unit tests and generates HTML coverage report                   |
| `make lint`          | Runs `golangci-lint` against all packages                            |
| `make fmt`           | Formats Go source code (`go fmt`)                                    |
| `make vet`           | Runs Go static analysis (`go vet`)                                   |
| `make tidy`          | Downloads dependencies and cleans `go.mod` / `go.sum`                |
| `make clean`         | Removes compiled binaries and coverage output files                  |

---

## 🧪 Testing

Run all unit and integration tests with data race detector:

```bash
make test
# or directly:
go test -v -race ./internal/...
```

Generate HTML test coverage:

```bash
make test-coverage
# Opens coverage.html in your browser
```

---

## 🔍 Code Linting

This project uses **[golangci-lint](https://golangci-lint.run/)** configured in `.golangci.yml`.

Run linting:
```bash
make lint
```

---

## 📡 API Endpoints & Documentation

### Swagger OpenAPI Documentation (Web UI)

- **URL**: `http://localhost:8080/swagger/index.html`
- **Method**: `GET`
- **Description**: Interactive Swagger UI for exploring and testing API endpoints.
- **Environment**: **Enabled in Development / Staging only** (`APP_ENV != "production"`). Automatically disabled (returns 404) in Production for security.
- **Code generation**: The `docs/` directory is **auto-generated**. Run `make swagger` to re-generate docs whenever API annotations change.

### Health Check

- **URL**: `/health` or `/api/v1/ping`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "success": true,
    "message": "Ticket Box API is healthy",
    "data": {
      "status": "UP",
      "timestamp": "2026-09-13T11:46:19.467472Z"
    }
  }
  ```

