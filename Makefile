.PHONY: help setup up down build build-no-cache restart restart-backend restart-frontend logs logs-backend logs-frontend db-shell clean debloat install frontend-dev backend-dev

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- Docker Compose Commands ---

up: ## Start all services with Docker Compose in detached mode
	docker compose up -d
	@$(MAKE) links

down: ## Stop all services
	docker compose down

build: ## Rebuild all Docker images
	docker compose build

build-no-cache: ## Rebuild all Docker images without cache
	docker compose build --no-cache

restart: ## Restart all Docker services
	docker compose restart

restart-backend: ## Restart only the backend service
	docker compose restart backend

restart-frontend: ## Restart only the frontend service
	docker compose restart frontend

logs: ## View logs for all services
	docker compose logs -f

logs-backend: ## View backend logs
	docker compose logs -f backend

logs-frontend: ## View frontend logs
	docker compose logs -f frontend

links: ## Show links to all local services
	@echo "===================================================="
	@echo "                   DASHBOARD                        "
	@echo "===================================================="
	@echo "Landing Page           : http://localhost:5173"
	@echo "Admin Panel            : http://localhost:5173/admin"
	@echo "pgAdmin                : http://localhost:5050"
	@echo "===================================================="

db-shell: ## Access the PostgreSQL database shell
	docker compose exec -it db psql -U myuser -d crud_db

# --- Local Development Commands ---

setup: ## Install dependencies, build images, and start the project
	@$(MAKE) install
	@$(MAKE) build
	@$(MAKE) up

install: ## Install both frontend (NPM) and backend (Go) dependencies locally
	cd frontend && npm install
	@if command -v go >/dev/null 2>&1; then \
		cd backend && go mod download; \
	else \
		echo "Go is not installed or not in PATH. Skipping backend local dependency installation."; \
	fi

frontend-dev: ## Run the frontend locally (requires Node.js)
	cd frontend && npm run dev

backend-dev: ## Run the backend locally (requires Go)
	cd backend && go run main.go

clean: ## Remove node_modules, dist, and clean Docker unused data
	rm -rf frontend/node_modules
	rm -rf frontend/dist
	docker system prune -f

debloat: ## Remove all project containers, volumes, and images
	docker compose down -v --rmi all
	docker system prune -a --volumes -f
