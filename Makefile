SHELL := /usr/bin/env bash

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BACKEND_DIR := $(ROOT_DIR)/backend
FRONTEND_DIR := $(ROOT_DIR)/frontend
VERSION := $(shell tr -d '[:space:]' < $(ROOT_DIR)/VERSION)
CONTAINER_CLI ?= podman
VITE_API_URL ?= http://localhost:8080

.PHONY: help version install backend-install frontend-install dev-backend dev-frontend \
	build backend-build frontend-build test backend-test frontend-check lint \
	format images clean

help:
	@printf "\nTargets disponíveis:\n\n"
	@printf "  %-18s %s\n" "help" "Mostra esta ajuda"
	@printf "  %-18s %s\n" "version" "Exibe a versão atual do projeto"
	@printf "  %-18s %s\n" "install" "Instala dependências de backend e frontend"
	@printf "  %-18s %s\n" "backend-install" "Baixa dependências Go"
	@printf "  %-18s %s\n" "frontend-install" "Instala dependências do frontend com Bun"
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
	@printf "  %-18s %s\n" "images" "Builda imagens versionadas com $(CONTAINER_CLI)"
	@printf "  %-18s %s\n\n" "clean" "Remove artefatos de build locais"

version:
	@printf "%s\n" "$(VERSION)"

install: backend-install frontend-install

backend-install:
	cd $(BACKEND_DIR) && go mod download

frontend-install:
	cd $(FRONTEND_DIR) && bun install

dev-backend:
	cd $(BACKEND_DIR) && go run main.go

dev-frontend:
	cd $(FRONTEND_DIR) && bun run dev

build: backend-build frontend-build

backend-build:
	cd $(BACKEND_DIR) && go build -o bin/server ./main.go

frontend-build:
	cd $(FRONTEND_DIR) && bun run build

test: backend-test frontend-check

backend-test:
	cd $(BACKEND_DIR) && go test ./...

frontend-check:
	cd $(FRONTEND_DIR) && bunx svelte-check --tsconfig ./tsconfig.json && bun run lint

lint: frontend-check

format:
	cd $(BACKEND_DIR) && gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

images:
	CONTAINER_CLI=$(CONTAINER_CLI) VITE_API_URL=$(VITE_API_URL) ./scripts/build-images.sh

clean:
	rm -rf $(BACKEND_DIR)/bin
	rm -rf $(FRONTEND_DIR)/build $(FRONTEND_DIR)/.svelte-kit
