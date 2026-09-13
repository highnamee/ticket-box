# Ticket Box Backend (Go + Gin)

Go backend for the **Ticket Box** system, built with [Gin Web Framework](https://github.com/gin-gonic/gin).

---

## 📁 Project Structure

```
apps/be/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point (dependency injection, graceful shutdown)
├── configs/
│   └── .env.example                # Sample environment variables
├── internal/                       # Private application code (Go compiler protected)
│   ├── config/                     # Configuration and environment variable loader
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain/                     # Core domain entities, errors, and interface definitions
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
│   │   └── response/               # Standardized JSON response envelope
│   │       ├── response.go
│   │       └── response_test.go
│   ├── repository/                 # Data access layer (PostgreSQL, Redis, etc.)
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
- **Make**: (Optional, for running Makefile shortcuts)
- **golangci-lint**: (Optional, for code quality & linting)

### Installation

1. Clone the repository and navigate to the backend directory:
   ```bash
   cd apps/be
   ```

2. Download and verify dependencies:
   ```bash
   go mod tidy
   ```

3. (Optional) Set up environment variables:
   ```bash
   cp .env.example .env
   ```

---

## 🛠️ Development & Commands

The project includes a `Makefile` with common tasks:

| Command              | Description                                             |
|----------------------|---------------------------------------------------------|
| `make run`           | Starts the API server locally                           |
| `make build`         | Compiles the binary to `bin/api`                        |
| `make migrate`       | Runs database migrations (`cmd/migrate`)                |
| `make test`          | Runs all unit tests with race condition detection       |
| `make test-coverage` | Runs unit tests and generates HTML coverage report      |
| `make lint`          | Runs `golangci-lint` against all packages               |
| `make fmt`           | Formats Go source code (`go fmt`)                       |
| `make vet`           | Runs Go static analysis (`go vet`)                      |
| `make tidy`          | Downloads dependencies and cleans `go.mod` / `go.sum`   |
| `make clean`         | Removes compiled binaries and coverage output files     |

---

## 🧪 Testing

Run all unit tests with data race detector:

```bash
make test
# or directly:
go test -v -race ./...
```

Generate HTML test coverage:

```bash
make test-coverage
# Opens coverage.html in your browser
```

---

## 🔍 Code Linting

This project uses **[golangci-lint](https://golangci-lint.run/)** configured in `.golangci.yml` with essential linters enabled (`govet`, `errcheck`, `staticcheck`, `unused`, `revive`, `gocyclo`, `misspell`, etc.).

To install `golangci-lint`:
- **macOS (Homebrew)**: `brew install golangci-lint`
- **Go Install**: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

To run linting:
```bash
make lint
# or:
golangci-lint run ./...
```

---

## 📡 API Endpoints

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
