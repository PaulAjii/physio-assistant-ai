# ─────────────────────────────────────────────────────────────────────────────
# Developer tasks for physio-assistant-ai.
#
# Container engine is Podman. If you use podman-compose instead of the built-in
# `podman compose`, override it:
#     make up COMPOSE="podman-compose"
# ─────────────────────────────────────────────────────────────────────────────

COMPOSE ?= podman compose
ENV_FILE := .env

# Load .env so recipes can use $(POSTGRES_USER) etc.
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

POSTGRES_USER ?= physio
POSTGRES_DB   ?= physio
APP_DB_USER   ?= physio_app

# Connection string for migrations. Uses the OWNER role (physio) on purpose:
# migrations create objects that must be owned by the owner so RLS applies to
# the runtime role. Built from .env, so the password stays in .env; migrate
# recipes are silenced with @ so it never prints to the terminal.
DB_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(or $(POSTGRES_PORT),5432)/$(POSTGRES_DB)?sslmode=disable
GOOSE  := cd backend-core && go tool goose -dir migrations postgres "$(DB_URL)"

.DEFAULT_GOAL := help
.PHONY: help env-check machine up down restart ps logs psql psql-app minio-console reset \
	migrate migrate-down migrate-status migrate-redo

help: ## Show available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

env-check: ## Fail early if .env is missing
	@test -f $(ENV_FILE) || { \
		echo "Missing $(ENV_FILE). Run: cp .env.example .env  (then edit the passwords)"; \
		exit 1; \
	}

machine: ## Start the Podman VM (macOS/Windows only; no-op if already running)
	@podman machine inspect >/dev/null 2>&1 \
		&& podman machine start 2>/dev/null || true
	@podman info >/dev/null 2>&1 \
		|| { echo "Podman is not reachable. Try: podman machine init && podman machine start"; exit 1; }

up: env-check ## Start Postgres + MinIO in the background
	$(COMPOSE) up -d
	@echo
	@echo "Postgres  : localhost:$(or $(POSTGRES_PORT),5432)  (db=$(POSTGRES_DB) owner=$(POSTGRES_USER) app=$(APP_DB_USER))"
	@echo "MinIO API : localhost:$(or $(MINIO_PORT),9000)"
	@echo "MinIO UI  : http://localhost:$(or $(MINIO_CONSOLE_PORT),9001)"

down: ## Stop the stack, keeping data volumes
	$(COMPOSE) down

restart: down up ## Stop then start the stack

migrate: env-check ## Apply all pending database migrations (as the OWNER role)
	@$(GOOSE) up

migrate-down: env-check ## Roll back the most recent migration
	@$(GOOSE) down

migrate-status: env-check ## Show which migrations have been applied
	@$(GOOSE) status

migrate-redo: env-check ## Roll back and re-apply the latest migration (tests reversibility)
	@$(GOOSE) redo

ps: ## Show container status
	$(COMPOSE) ps

logs: ## Follow logs from all services
	$(COMPOSE) logs -f

psql: ## psql shell as the OWNER role (bypasses RLS — for migrations/inspection)
	$(COMPOSE) exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

psql-app: ## psql shell as the RUNTIME role (RLS enforced — use to verify isolation)
	$(COMPOSE) exec postgres psql -U $(APP_DB_USER) -d $(POSTGRES_DB)

minio-console: ## Open the MinIO web console
	@open "http://localhost:$(or $(MINIO_CONSOLE_PORT),9001)" 2>/dev/null \
		|| echo "Visit http://localhost:$(or $(MINIO_CONSOLE_PORT),9001)"

reset: ## DESTRUCTIVE: delete all containers AND data volumes, then recreate
	@printf "This deletes the Postgres and MinIO volumes. All local data is lost.\nType 'yes' to continue: " \
		&& read ans && [ "$$ans" = "yes" ] || { echo "Aborted."; exit 1; }
	$(COMPOSE) down -v
	$(MAKE) up
