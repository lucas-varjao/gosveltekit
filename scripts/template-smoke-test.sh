#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONTAINER_CLI="${CONTAINER_CLI:-docker}"
TEMP_DIR="$(mktemp -d)"
BACKEND_PID=""

cleanup() {
    if [[ -n "${BACKEND_PID}" ]] && kill -0 "${BACKEND_PID}" >/dev/null 2>&1; then
        kill "${BACKEND_PID}" >/dev/null 2>&1 || true
        wait "${BACKEND_PID}" 2>/dev/null || true
    fi

    if [[ -f "${TEMP_DIR}/compose.yml" ]]; then
        (
            cd "${TEMP_DIR}"
            CONTAINER_CLI="${CONTAINER_CLI}" ./scripts/compose.sh down -v >/dev/null 2>&1 || true
        )
    fi

    rm -rf "${TEMP_DIR}"
}
trap cleanup EXIT

tar \
    --exclude=.git \
    --exclude=frontend/node_modules \
    --exclude=frontend/build \
    --exclude=frontend/.svelte-kit \
    --exclude=backend/bin \
    -cf - -C "${ROOT_DIR}" . | tar -xf - -C "${TEMP_DIR}"

cd "${TEMP_DIR}"

make init \
    APP_NAME=acme-starter \
    DISPLAY_NAME='Acme Starter' \
    GO_MODULE=github.com/acme/acme-starter \
    CONTAINER_REGISTRY=ghcr.io/acme \
    DOMAIN=acme-starter.local \
    K8S_NAMESPACE=acme-starter

if rg -n 'gosveltekit|GoSvelteKit|gosveltekit.local' \
    README.md backend frontend k8s compose.yml .env.example backend/.env.example frontend/.env.example Makefile project.env AGENTS.md >/dev/null; then
    echo "found stale template references after init" >&2
    exit 1
fi

make bootstrap
make infra-up CONTAINER_CLI="${CONTAINER_CLI}"

source .env

for attempt in {1..30}; do
    if "${CONTAINER_CLI}" compose -f compose.yml exec -T postgres \
        pg_isready -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" >/dev/null 2>&1; then
        break
    fi

    if [[ "${attempt}" -eq 30 ]]; then
        echo "postgres did not become ready" >&2
        exit 1
    fi

    sleep 2
done

make migrate-up
make seed-admin \
    ADMIN_IDENTIFIER=admin \
    ADMIN_EMAIL=admin@acme-starter.local \
    ADMIN_PASSWORD='Starter123!' \
    ADMIN_DISPLAY_NAME='Platform Admin'
make build
CONTAINER_CLI="${CONTAINER_CLI}" VITE_API_URL='' ./scripts/build-images.sh >/dev/null

make dev-backend >"${TEMP_DIR}/backend.log" 2>&1 &
BACKEND_PID=$!

for attempt in {1..30}; do
    if curl --fail --silent http://localhost:8080/health >/dev/null 2>&1; then
        break
    fi

    if [[ "${attempt}" -eq 30 ]]; then
        echo "backend did not become healthy" >&2
        cat "${TEMP_DIR}/backend.log" >&2 || true
        exit 1
    fi

    sleep 2
done

LOGIN_RESPONSE="$(curl --fail --silent \
    -X POST http://localhost:8080/auth/login \
    -H 'Content-Type: application/json' \
    -d '{"username":"admin","password":"Starter123!"}')"

if [[ "${LOGIN_RESPONSE}" != *'"session_id"'* ]]; then
    echo "login smoke test failed" >&2
    echo "${LOGIN_RESPONSE}" >&2
    exit 1
fi

echo "template smoke test passed"
