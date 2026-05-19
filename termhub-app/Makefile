.PHONY: help dev dev-infra dev-backend dev-frontend build up down restart logs clean status

# Default target
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Development ────────────────────────────────────────────────

dev-infra: ## Start dev infrastructure (Postgres + Redis)
	docker compose -f docker-compose.dev.yml up -d

dev-backend: ## Run backend locally (requires dev-infra)
	cd backend && go run cmd/server/main.go

dev-frontend: ## Run frontend dev server
	cd frontend && npm run dev

dev: dev-infra ## Start full dev environment (infra + instructions)
	@echo ""
	@echo "✅ Dev infrastructure started (Postgres + Redis)"
	@echo ""
	@echo "Now run in separate terminals:"
	@echo "  make dev-backend   # Go backend on :8080"
	@echo "  make dev-frontend  # Vite dev server on :5173"

# ─── Production (Docker Compose) ───────────────────────────────

build: ## Build all Docker images
	docker compose build

up: ## Start production stack
	docker compose up -d

down: ## Stop production stack
	docker compose down

restart: ## Restart production stack
	docker compose down && docker compose up -d

logs: ## Tail all container logs
	docker compose logs -f

logs-backend: ## Tail backend logs only
	docker compose logs -f backend

status: ## Show container status
	docker compose ps

# ─── Maintenance ────────────────────────────────────────────────

clean: ## Remove containers, volumes, and images
	docker compose down -v --rmi local
	docker compose -f docker-compose.dev.yml down -v

setup-env: ## Create .env from example (won't overwrite existing)
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ Created .env — edit it with your production secrets"; \
	else \
		echo "⚠️  .env already exists, skipping"; \
	fi

# ─── Testing / Linting ─────────────────────────────────────────

test-backend: ## Run Go tests
	cd backend && go test ./...

lint-backend: ## Run go vet
	cd backend && go vet ./...

build-backend: ## Build Go binary locally
	cd backend && go build -o server ./cmd/server/

build-frontend: ## Build frontend locally
	cd frontend && npm run build

# ─── Tailscale Funnel Deployment ────────────────────────────────

deploy: ## Deploy with Tailscale Funnel (run with sudo)
	./scripts/deploy-tailscale.sh

teardown: ## Stop TermHub and disable Tailscale Funnel
	./scripts/teardown-tailscale.sh

funnel-status: ## Show Tailscale Funnel status
	tailscale funnel status

funnel-off: ## Disable public access (keep app running)
	sudo tailscale funnel 443 off

funnel-on: ## Re-enable public access
	sudo tailscale funnel --bg 443 on

typecheck: ## Run TypeScript type check
	cd frontend && npx tsc --noEmit
