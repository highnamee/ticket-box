# Ticket Box

A modern event ticketing and booking platform built with **Go (Gin)** and **Next.js 16 (App Router)**.

---

## 📁 Repository Overview

```text
ticket-box/
├── apps/
│   ├── be/         # Backend API (Go, Gin, GORM, PostgreSQL)
│   └── fe/         # Frontend Web (Next.js 16, Tailwind v4, shadcn/ui)
├── docker-compose.yml
└── Makefile
```

| Component | Directory | Documentation |
| :--- | :--- | :--- |
| **Backend API** | [apps/be](apps/be) | [README](apps/be/README.md) • [AGENTS](apps/be/AGENTS.md) |
| **Frontend Web** | [apps/fe](apps/fe) | [README](apps/fe/README.md) • [AGENTS](apps/fe/AGENTS.md) • [DESIGN](apps/fe/DESIGN.md) |

---

## 🚀 Quick Start (Makefile)

Manage the entire platform using root `make` commands:

```bash
# Start all services in Docker (PostgreSQL, Backend, Frontend)
make up

# Rebuild and start all services
make up-build

# Check status of running containers
make ps

# Follow logs across all services
make logs

# Stop all services and remove containers
make down

# View all available make targets
make help
```

For detailed setup, environment variables, and architecture rules, refer to [apps/be/README.md](apps/be/README.md) and [apps/fe/README.md](apps/fe/README.md).
