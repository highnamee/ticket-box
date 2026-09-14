# Ticket Box

A ticket booking and event management platform.

## Project Structure

```
ticket-box/
├── apps/
│   ├── be/     # Golang backend (Gin framework)
│   └── fe/     # Frontend application
└── docker-compose.yml
```

## Quick Start with Docker Compose

Start the PostgreSQL database and backend API with a single command:

```bash
# Start all services (Database + Backend)
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f backend

# Stop all services
docker compose down
```

### Start Only PostgreSQL (For local Go development)
```bash
docker compose up -d postgres
```
