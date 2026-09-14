# Backend Development Guidelines for AI Agents

Guidelines and architectural rules for AI assistants working on the **Ticket Box Backend**.

---

## 🛠️ Tech Stack & Constraints

- **Language**: Go 1.23+
- **HTTP Framework**: Gin (`github.com/gin-gonic/gin`)
- **Database / ORM**: PostgreSQL via GORM (`gorm.io/driver/postgres`, `gorm.io/gorm`)
- **Testing**: Native Go testing + `github.com/stretchr/testify`
- **Database Engine**: **Pure PostgreSQL only**. Do NOT add SQLite or CGO dependencies.

---

## 🏗️ Architectural Rules (Clean Architecture)

Follow the established layered pattern strictly:

```
[HTTP Request] ➔ [Handler] (internal/handler/http)
                     │
                     ▼
                 [Service] (internal/service)
                     │
                     ▼
           [Repository Interface] (internal/domain)
                     │
                     ▼
        [PostgreSQL Repository] (internal/repository/postgres)
```

1. **Domain Layer (`internal/domain/`)**:
   - Contains entity structs, value objects, domain errors, and interface contracts (`Repository`, `Service`).
   - Must remain pure Go: **No framework imports (Gin, HTTP, net/http) or database driver dependencies**.
   - Define sentinel errors with standard `errors.New` (e.g. `ErrTicketNotFound = errors.New("ticket not found")`).
   - Primary Keys: Always use **UUID v7 (Time-Ordered)** with `ID uuid.UUID` (`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`).
   - Hook: Include `BeforeCreate(_ *gorm.DB) error` generating `uuid.NewV7()` when `t.ID == uuid.Nil`:
     ```go
     func (t *Entity) BeforeCreate(_ *gorm.DB) error {
         if t.ID == uuid.Nil {
             id, err := uuid.NewV7()
             if err != nil {
                 return err
             }
             t.ID = id
         }
         return nil
     }
     ```

2. **Repository Layer (`internal/repository/postgres/`)**:
   - Implements domain repository interfaces using GORM and PostgreSQL.
   - For concurrency-sensitive operations (e.g. ticket booking/reservation), use Pessimistic Locking (`clause.Locking{Strength: "UPDATE"}`) inside database transactions to prevent race conditions and overbooking.

3. **Service Layer (`internal/service/`)**:
   - Implements business use cases and orchestrates repository calls.

4. **Handler Layer (`internal/handler/http/`)**:
   - Handles HTTP request binding, validation, controller flow, and response formatting.
   - **File Convention**: Separate handler logic from request/response DTOs:
     - `<feature>_handler.go`: Controller methods, routing, service/repo calls, and HTTP responses.
     - `<feature>_dto.go`: Request queries/bodies with Gin binding tags (`form:`, `json:`), Response Serializers/DTOs, and mapping functions.
     - `<feature>_handler_test.go`: Unit tests for endpoints and serialization.
   - Returns standardized API responses via `internal/pkg/response`:
     - Regular endpoints: `response.Success(c, http.StatusOK, "message", data)`
     - Paginated endpoints: `response.SuccessWithPagination(c, http.StatusOK, "message", items, paginationMeta)` (wraps `{ items, pagination }` inside `data`)
     - Error handling: `response.BadRequest(c, message, err)`, `response.NotFound(c, message, err)`, `response.InternalServerError(c, message, err)`.
   - Never write database queries directly inside handlers.

---

## 🗄️ Database & Migrations

- All database models must be registered in `internal/pkg/database/migration.go` inside the `models` slice in `AutoMigrate(db)`.
- Migrations are run via CLI entrypoint: `cmd/migrate/main.go` or `make migrate`.

---

## 🧪 Testing Guidelines

1. **Unit Tests**:
   - Mock repository interfaces when testing service logic.
2. **Integration Tests**:
   - Must run against real PostgreSQL using `testutil.SetupTestDB(t)`.
   - `testutil.SetupTestDB(t)` automatically wraps each test in an isolated transaction and calls `tx.Rollback()` on cleanup.
   - Do NOT use SQLite for tests.
3. **Running Tests**:
   - Always run tests targeting internal packages: `make test` (or `go test -v -race ./internal/...`).
   - Do NOT run coverage on `cmd/...` to avoid `covdata` toolchain issues.

---

## 🛠️ Admin Panel (`internal/admin/`)

- GoAdmin web interface is mounted at `/admin` (accessible with `admin` / `admin`).
- To add a new entity to the Admin panel:
  1. Define a `Get<Entity>Table(ctx *context.Context) table.Table` generator in `internal/admin/tables.go`.
  2. Register the table generator in the `Generators` map in `internal/admin/admin.go`.

---

## ⚡ Useful Commands

```bash
make setup          # Install developer CLI tools (swag, golangci-lint) & dependencies
make run            # Start API server locally
make migrate        # Run database migrations
make swagger        # Generate OpenAPI/Swagger docs into docs/
make test           # Run unit & PostgreSQL integration tests
make test-coverage  # Run tests with HTML coverage report
make lint           # Run golangci-lint
make vet            # Run go vet
make fmt            # Format Go code
make tidy           # Sync go.mod / go.sum
```
