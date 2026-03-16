.PHONY: help setup install dev build up down restart logs db-shell clean test debloat links swagger refresh

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- Local Development ---

setup: ## Build images and start the project in Docker
	@$(MAKE) install
	@$(MAKE) build
	@$(MAKE) up

refresh: ## Regenerate swagger, rebuild, and restart services
	@$(MAKE) swagger
	@$(MAKE) build
	@$(MAKE) up

install: ## Install Go dependencies locally (optional)
	@if command -v go >/dev/null 2>&1; then \
		echo "Installing Go dependencies..."; \
		go mod download; \
	else \
		echo "Notice: Go is not installed on host. Skipping local dependency download."; \
		echo "The project will still work via Docker."; \
	fi

dev: ## Run the backend locally
	@if [ ! -f .env ]; then \
		echo "Warning: .env file not found. Running with defaults."; \
	fi
	go run main.go

test: ## Run tests
	go test ./... -v

# --- Docker Operations ---

up: ## Start the backend and database in Docker
	docker compose up -d
	@$(MAKE) links

links: ## Show links to the running services
	@echo "===================================================="
	@echo "Backend API  : http://localhost:5000/api"
	@echo "Swagger UI   : http://localhost:5000/swagger/index.html"
	@echo "Health Check : http://localhost:5000/api/health"
	@echo "Products     : http://localhost:5000/api/products"
	@echo "pgAdmin      : http://localhost:5050"
	@echo "===================================================="

down: ## Stop all Docker services
	docker compose down

build: ## Build or rebuild Docker images
	docker compose build

build-no-cache: ## Rebuild Docker images without cache
	docker compose build --no-cache

restart: ## Restart the backend service in Docker
	docker compose restart backend

logs: ## View backend logs in Docker
	docker compose logs -f backend

db-shell: ## Access the PostgreSQL database shell
	docker compose exec -it db psql -U myuser -d crud_db

# --- Utilities ---

swagger: ## Generate Swagger documentation
	docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go install github.com/swaggo/swag/cmd/swag@latest && swag init --parseDependency"

clean: ## Remove Go build artifacts and tidy modules
	rm -f backend-go
	go mod tidy
	docker system prune -f

debloat: ## Remove all project containers, volumes, and images
	docker compose down -v --rmi all
	docker system prune -a --volumes -f
