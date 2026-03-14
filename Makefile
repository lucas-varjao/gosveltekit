SHELL := /usr/bin/env bash

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BACKEND_DIR := $(ROOT_DIR)/backend
FRONTEND_DIR := $(ROOT_DIR)/frontend
VERSION := $(shell tr -d '[:space:]' < $(ROOT_DIR)/VERSION)
CONTAINER_CLI ?= podman
VITE_API_URL ?= http://localhost:8080
ADMIN_DISPLAY_NAME ?= Administrator
KUBECTL ?= kubectl
APP_SLUG ?= $(shell awk -F= '/^APP_SLUG=/{gsub(/"/, "", $$2); print $$2}' $(ROOT_DIR)/project.env)
K8S_NAMESPACE ?= $(shell awk -F= '/^K8S_NAMESPACE=/{gsub(/"/, "", $$2); print $$2}' $(ROOT_DIR)/project.env)
K8S_BASE_MANIFEST ?= $(ROOT_DIR)/k8s/$(APP_SLUG)-base.yaml
K8S_MIGRATE_JOB_MANIFEST ?= $(ROOT_DIR)/k8s/$(APP_SLUG)-migrate.job.yaml
K8S_MIGRATE_JOB_NAME ?= $(APP_SLUG)-migrate
K8S_MIGRATE_TIMEOUT ?= 5m

INIT_ARGS :=
ifdef APP_NAME
INIT_ARGS += --app-name "$(APP_NAME)"
endif
ifdef DISPLAY_NAME
INIT_ARGS += --display-name "$(DISPLAY_NAME)"
endif
ifdef GO_MODULE
INIT_ARGS += --go-module "$(GO_MODULE)"
endif
ifdef CONTAINER_REGISTRY
INIT_ARGS += --container-registry "$(CONTAINER_REGISTRY)"
endif
ifdef DOMAIN
INIT_ARGS += --domain "$(DOMAIN)"
endif
ifdef K8S_NAMESPACE
INIT_ARGS += --k8s-namespace "$(K8S_NAMESPACE)"
endif

.PHONY: help version init bootstrap install backend-install frontend-install infra-up \
	infra-down dev-backend dev-frontend build backend-build frontend-build test \
	backend-test frontend-check lint format images migrate-up migrate-down \
	migrate-create seed-admin images k8s-migrate-job clean

help:
	@printf "\nTargets disponíveis:\n\n"
	@printf "  %-18s %s\n" "help" "Mostra esta ajuda"
	@printf "  %-18s %s\n" "version" "Exibe a versão atual do projeto"
	@printf "  %-18s %s\n" "init" "Inicializa um novo projeto a partir do template"
	@printf "  %-18s %s\n" "bootstrap" "Gera .env locais e instala dependências"
	@printf "  %-18s %s\n" "install" "Instala dependências de backend e frontend"
	@printf "  %-18s %s\n" "backend-install" "Baixa dependências Go"
	@printf "  %-18s %s\n" "frontend-install" "Instala dependências do frontend com Bun"
	@printf "  %-18s %s\n" "infra-up" "Sobe PostgreSQL e Mailpit via compose"
	@printf "  %-18s %s\n" "infra-down" "Derruba a infraestrutura local"
	@printf "  %-18s %s\n" "dev-backend" "Inicia o backend localmente"
	@printf "  %-18s %s\n" "dev-frontend" "Inicia o frontend localmente"
	@printf "  %-18s %s\n" "build" "Executa build de backend e frontend"
	@printf "  %-18s %s\n" "backend-build" "Compila o backend"
	@printf "  %-18s %s\n" "frontend-build" "Gera o build do frontend"
	@printf "  %-18s %s\n" "test" "Executa testes do backend e checagens do frontend"
	@printf "  %-18s %s\n" "backend-test" "Executa go test ./..."
	@printf "  %-18s %s\n" "frontend-check" "Executa lint e type-check do frontend"
	@printf "  %-18s %s\n" "lint" "Alias para frontend-check"
	@printf "  %-18s %s\n" "format" "Formata o backend com gofmt"
	@printf "  %-18s %s\n" "migrate-up" "Aplica migrações do banco"
	@printf "  %-18s %s\n" "migrate-down" "Reverte a última migração do banco"
	@printf "  %-18s %s\n" "migrate-create" "Cria um novo arquivo de migração goose"
	@printf "  %-18s %s\n" "seed-admin" "Cria ou atualiza o usuário administrador"
	@printf "  %-18s %s\n" "images" "Builda imagens versionadas com $(CONTAINER_CLI)"
	@printf "  %-18s %s\n" "k8s-migrate-job" "Aplica base, recria e aguarda o Job de migração"
	@printf "  %-18s %s\n\n" "clean" "Remove artefatos de build locais"

version:
	@printf "%s\n" "$(VERSION)"

init:
	./scripts/init-project.sh $(INIT_ARGS) $(ARGS)

bootstrap:
	./scripts/bootstrap-env.sh
	$(MAKE) install

install: backend-install frontend-install

backend-install:
	cd $(BACKEND_DIR) && go mod download

frontend-install:
	cd $(FRONTEND_DIR) && bun install

infra-up:
	CONTAINER_CLI=$(CONTAINER_CLI) ./scripts/compose.sh up -d

infra-down:
	CONTAINER_CLI=$(CONTAINER_CLI) ./scripts/compose.sh down

dev-backend:
	cd $(BACKEND_DIR) && set -a && if [ -f .env ]; then . ./.env; fi && set +a && go run .

dev-frontend:
	cd $(FRONTEND_DIR) && set -a && if [ -f .env ]; then . ./.env; fi && set +a && bun run dev

build: backend-build frontend-build

backend-build:
	cd $(BACKEND_DIR) && go build -o bin/server .

frontend-build:
	cd $(FRONTEND_DIR) && set -a && if [ -f .env ]; then . ./.env; fi && set +a && bun run build

test: backend-test frontend-check

backend-test:
	cd $(BACKEND_DIR) && go test ./...

frontend-check:
	cd $(FRONTEND_DIR) && bunx svelte-check --tsconfig ./tsconfig.json && bun run lint

lint: frontend-check

format:
	cd $(BACKEND_DIR) && gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

migrate-up:
	cd $(BACKEND_DIR) && set -a && if [ -f .env ]; then . ./.env; fi && set +a && go run ./cmd/migrate up

migrate-down:
	cd $(BACKEND_DIR) && set -a && if [ -f .env ]; then . ./.env; fi && set +a && go run ./cmd/migrate down

migrate-create:
	cd $(BACKEND_DIR) && go run ./cmd/migrate create $(name)

seed-admin:
	cd $(BACKEND_DIR) && set -a && if [ -f .env ]; then . ./.env; fi && set +a && \
	ADMIN_IDENTIFIER="$(ADMIN_IDENTIFIER)" \
	ADMIN_EMAIL="$(ADMIN_EMAIL)" \
	ADMIN_PASSWORD="$(ADMIN_PASSWORD)" \
	ADMIN_DISPLAY_NAME="$(ADMIN_DISPLAY_NAME)" \
	go run ./cmd/seed-admin

images:
	CONTAINER_CLI=$(CONTAINER_CLI) VITE_API_URL=$(VITE_API_URL) ./scripts/build-images.sh

k8s-migrate-job:
	$(KUBECTL) apply -f $(K8S_BASE_MANIFEST)
	$(KUBECTL) delete -f $(K8S_MIGRATE_JOB_MANIFEST) --ignore-not-found
	$(KUBECTL) create -f $(K8S_MIGRATE_JOB_MANIFEST)
	$(KUBECTL) wait --for=condition=complete job/$(K8S_MIGRATE_JOB_NAME) -n $(K8S_NAMESPACE) --timeout=$(K8S_MIGRATE_TIMEOUT)

clean:
	rm -rf $(BACKEND_DIR)/bin
	rm -rf $(FRONTEND_DIR)/build $(FRONTEND_DIR)/.svelte-kit
