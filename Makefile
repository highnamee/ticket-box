.PHONY: help up up-build down logs ps

help: ## Show available commands
	@echo "Ticket Box - Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

up: ## Start all services (PostgreSQL, Backend, Frontend)
	docker compose up -d

up-build: ## Rebuild and start all services
	docker compose up --build -d

down: ## Stop all services and remove containers
	docker compose down

logs: ## Follow logs from all containers
	docker compose logs -f

ps: ## Check status of containers
	docker compose ps
