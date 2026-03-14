#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION_FILE="${ROOT_DIR}/VERSION"
CONTAINER_CLI="${CONTAINER_CLI:-podman}"
VITE_API_URL="${VITE_API_URL:-http://localhost:8080}"

if [[ ! -f "${VERSION_FILE}" ]]; then
    echo "VERSION file not found at ${VERSION_FILE}" >&2
    exit 1
fi

APP_VERSION="$(tr -d '[:space:]' < "${VERSION_FILE}")"

if [[ -z "${APP_VERSION}" ]]; then
    echo "VERSION file is empty" >&2
    exit 1
fi

if [[ ! "${APP_VERSION}" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "VERSION must use MAJOR.MINOR.PATCH format, got: ${APP_VERSION}" >&2
    exit 1
fi

if ! command -v "${CONTAINER_CLI}" >/dev/null 2>&1; then
    echo "Container CLI not found: ${CONTAINER_CLI}" >&2
    exit 1
fi

echo "Building images with ${CONTAINER_CLI} for version ${APP_VERSION}"

"${CONTAINER_CLI}" build \
    --build-arg APP_VERSION="${APP_VERSION}" \
    -t "gosveltekit-backend:${APP_VERSION}" \
    -t "gosveltekit-backend:latest" \
    "${ROOT_DIR}/backend"

"${CONTAINER_CLI}" build \
    --build-arg APP_VERSION="${APP_VERSION}" \
    --build-arg VITE_API_URL="${VITE_API_URL}" \
    -t "gosveltekit-frontend:${APP_VERSION}" \
    -t "gosveltekit-frontend:latest" \
    "${ROOT_DIR}/frontend"

echo "Built images:"
echo "  gosveltekit-backend:${APP_VERSION}"
echo "  gosveltekit-backend:latest"
echo "  gosveltekit-frontend:${APP_VERSION}"
echo "  gosveltekit-frontend:latest"
