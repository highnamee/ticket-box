# Ticket Box Backend (Go + Gin)

Go backend for the **Ticket Box** system, built with [Gin Web Framework](https://github.com/gin-gonic/gin).

---

## 📁 Project Structure

```
apps/be/
├── cmd/
│   ├── api/main.go                 # Application entry point (Composition Root, DI wiring, graceful shutdown)
│   └── migrate/main.go             # Goose Versioned Database Migration CLI (up, down, status, reset)
├── docs/                           # Auto-generated OpenAPI / Swagger specs (docs.go, swagger.json, swagger.yaml)
├── internal/                       # Private application code
│   ├── admin/                      # GoAdmin web panel setup & entity table generators (tickets, users)
│   ├── config/                     # Type-safe configuration loader (caarlos0/env)
│   ├── domain/                     # Core domain entities, errors, and interface contracts (e.g. ticket, user)
│   ├── handler/http/               # HTTP controllers & DTOs (e.g. auth_handler, ticket_handler)
│   ├── middleware/                 # Gin HTTP middlewares (Auth/JWT, CORS, Request ID, Slog Logger)
│   ├── pkg/                        # Shared packages (token/jwt, hasher, database, response, validator)
│   ├── repository/postgres/        # PostgreSQL data access layer via GORM (e.g. user_repo, ticket_repo)
│   └── service/                    # Business logic layer (e.g. auth_service, ticket_service)
├── migrations/                     # SQL Goose Versioned Migration files (e.g. 00001_*.sql, 00002_*.sql)
├── Dockerfile                      # Multi-stage production container build
├── Makefile                        # Developer task automation
└── README.md
```

---

## 🔐 Authentication & API Endpoints

The API uses **JWT (JSON Web Tokens)** via `golang-jwt/jwt/v5` and **Bcrypt** password hashing.

### Public Auth Endpoints
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Register a new user account |
| `POST` | `/api/v1/auth/login` | Log in and receive Access + Refresh JWT tokens |
| `POST` | `/api/v1/auth/refresh` | Exchange valid Refresh Token for a new token pair |
| `POST` | `/api/v1/auth/forgot-password` | Request password reset link/token (Anti-enumeration protected) |
| `POST` | `/api/v1/auth/reset-password` | Reset password using single-use reset token |

### Protected Endpoints (Requires `Authorization: Bearer <token>`)
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/users/me` | Get current authenticated user profile |
| `GET` | `/api/v1/tickets` | List public tickets (Paginated) |

---

## 🐳 Docker & Docker Compose

### 1. Full Stack with Docker Compose (PostgreSQL + Backend)
From root or via Makefile:
```bash
# Start both PostgreSQL and Backend
make compose-up

# Start ONLY PostgreSQL (for local Go development)
make db-up

# Stop all containers
make compose-down
```

### 2. Standalone Docker Image
```bash
# Build Docker image
make docker-build

# Run Backend container standalone
make docker-run

# Run database migrations via container
docker run --rm --env-file .env --entrypoint /app/migrate ticket-box-be:latest up
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

4. Run database migrations:
   ```bash
   make migrate
   ```

---

## 🛠️ Development & Commands

The project includes a `Makefile` with common tasks:

| Command              | Description                                                          |
|----------------------|----------------------------------------------------------------------|
| `make setup`         | Installs developer tools (`swag`, `golangci-lint`) and dependencies  |
| `make run`           | Starts the API server locally                                        |
| `make build`         | Compiles the binary to `bin/api`                                     |
| `make migrate`       | Runs pending Goose database migrations (`Up`)                        |
| `make migrate-down`  | Rolls back the most recent Goose database migration (`Down`)         |
| `make migrate-status`| Checks current status of all Goose migrations                        |
| `make migrate-create`| Creates a new migration file (`make migrate-create name=name`)       |
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



