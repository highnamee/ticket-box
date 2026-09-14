# Ticket Box Backend (Go + Gin)

Go backend for the **Ticket Box** system, built with [Gin Web Framework](https://github.com/gin-gonic/gin).

---

## 📁 Project Structure

```
apps/be/
├── cmd/
│   ├── api/
│   │   └── main.go                 # Application entry point (Composition Root, DI wiring, graceful shutdown)
│   └── migrate/
│       └── main.go                 # Database auto-migration CLI
├── docs/                           # Auto-generated OpenAPI / Swagger specs (docs.go, swagger.json, swagger.yaml)
├── internal/                       # Private application code (Go compiler protected)
│   ├── admin/                      # GoAdmin web panel setup & model table generators
│   │   ├── admin.go
│   │   ├── admin_test.go
│   │   ├── schema.go
│   │   ├── schema.sql
│   │   └── tables.go
│   ├── config/                     # Type-safe configuration loader (caarlos0/env)
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain/                     # Core domain entities, errors, and interface contracts
│   │   ├── errors.go
│   │   ├── ticket.go
│   │   └── ticket_test.go
│   ├── handler/                    # Transport layer (HTTP / Gin controllers)
│   │   └── http/
│   │       ├── health_handler.go
│   │       ├── health_handler_test.go
│   │       ├── router.go
│   │       ├── router_test.go
│   │       ├── ticket_dto.go
│   │       ├── ticket_handler.go
│   │       └── ticket_handler_test.go
│   ├── middleware/                 # Gin HTTP middlewares (CORS, Request ID, Slog Logger, Recovery)
│   │   ├── cors.go
│   │   ├── cors_test.go
│   │   ├── logger.go
│   │   ├── logger_test.go
│   │   ├── request_id.go
│   │   └── request_id_test.go
│   ├── pkg/                        # Shared libraries and packages
│   │   ├── database/               # Database connection, pool config, migration, and teardown
│   │   ├── logger/                 # Structured logging with Go log/slog and context tracing
│   │   ├── response/               # Standardized JSON response envelope & pagination
│   │   ├── testutil/               # Test isolation DB helpers (transaction rollback)
│   │   └── validator/              # Custom validation tags for Gin (uuidv7, notblank)
│   ├── repository/                 # Data access layer (PostgreSQL GORM implementations)
│   │   └── postgres/
│   │       ├── ticket_repo.go
│   │       └── ticket_repo_test.go
│   └── service/                    # Business logic and use-cases layer
│       ├── ticket_service.go
│       └── ticket_service_test.go
├── .env.example                    # Sample environment variables
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

## 📡 API Documentation & Admin Panel

### Swagger OpenAPI Documentation (Web UI)

- **URL**: `http://localhost:8080/swagger/index.html`
- **Description**: Interactive Swagger UI for exploring and testing API endpoints.
- **Environment**: **Enabled in Development / Staging only** (`APP_ENV != "production"`). Automatically disabled (returns 404) in Production for security.
- **Code generation**: The `docs/` directory is **auto-generated**. Run `make swagger` to re-generate docs whenever API annotations change.

### GoAdmin Management Panel (Web UI)

- **URL**: `http://localhost:8080/admin`
- **Credentials (Configurable)**: Defaults to `admin` / `admin` in development via `ADMIN_USERNAME` and `ADMIN_PASSWORD` env variables.
- **Environment**: **Disabled by default in Production** (`APP_ENV=production`) unless `ENABLE_ADMIN=true` is explicitly configured with strong credentials.
- **Description**: Full-featured database & entity management dashboard (AdminLTE theme) for managing Tickets, inventory, roles, permissions, and audit logs.



