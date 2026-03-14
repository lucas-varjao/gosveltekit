#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION_FILE="${ROOT_DIR}/VERSION"
PROJECT_ENV_FILE="${ROOT_DIR}/project.env"
CONTAINER_CLI="${CONTAINER_CLI:-podman}"
VITE_API_URL="${VITE_API_URL:-http://localhost:8080}"
PUSH_IMAGES="${PUSH_IMAGES:-false}"

if [[ -f "${PROJECT_ENV_FILE}" ]]; then
    # shellcheck disable=SC1090
    source "${PROJECT_ENV_FILE}"
fi

BACKEND_IMAGE_NAME="${BACKEND_IMAGE_NAME:-gosveltekit-backend}"
FRONTEND_IMAGE_NAME="${FRONTEND_IMAGE_NAME:-gosveltekit-frontend}"
MIGRATOR_IMAGE_NAME="${MIGRATOR_IMAGE_NAME:-gosveltekit-migrator}"
BACKEND_IMAGE_REF="${BACKEND_IMAGE_REF:-${BACKEND_IMAGE_NAME}}"
FRONTEND_IMAGE_REF="${FRONTEND_IMAGE_REF:-${FRONTEND_IMAGE_NAME}}"
MIGRATOR_IMAGE_REF="${MIGRATOR_IMAGE_REF:-${MIGRATOR_IMAGE_NAME}}"

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

if [[ ! "${PUSH_IMAGES}" =~ ^(true|false)$ ]]; then
    echo "PUSH_IMAGES must be true or false, got: ${PUSH_IMAGES}" >&2
    exit 1
fi

push_image_ref() {
    local image_ref="$1"

    echo "Pushing ${image_ref}:${APP_VERSION}"
    "${CONTAINER_CLI}" push "${image_ref}:${APP_VERSION}"

    echo "Pushing ${image_ref}:latest"
    "${CONTAINER_CLI}" push "${image_ref}:latest"
}

echo "Building images with ${CONTAINER_CLI} for version ${APP_VERSION}"

"${CONTAINER_CLI}" build \
    --target backend \
    --build-arg APP_VERSION="${APP_VERSION}" \
    -t "${BACKEND_IMAGE_NAME}:${APP_VERSION}" \
    -t "${BACKEND_IMAGE_NAME}:latest" \
    -t "${BACKEND_IMAGE_REF}:${APP_VERSION}" \
    -t "${BACKEND_IMAGE_REF}:latest" \
    "${ROOT_DIR}/backend"

"${CONTAINER_CLI}" build \
    --target migrator \
    --build-arg APP_VERSION="${APP_VERSION}" \
    -t "${MIGRATOR_IMAGE_NAME}:${APP_VERSION}" \
    -t "${MIGRATOR_IMAGE_NAME}:latest" \
    -t "${MIGRATOR_IMAGE_REF}:${APP_VERSION}" \
    -t "${MIGRATOR_IMAGE_REF}:latest" \
    "${ROOT_DIR}/backend"

"${CONTAINER_CLI}" build \
    --build-arg APP_VERSION="${APP_VERSION}" \
    --build-arg VITE_API_URL="${VITE_API_URL}" \
    -t "${FRONTEND_IMAGE_NAME}:${APP_VERSION}" \
    -t "${FRONTEND_IMAGE_NAME}:latest" \
    -t "${FRONTEND_IMAGE_REF}:${APP_VERSION}" \
    -t "${FRONTEND_IMAGE_REF}:latest" \
    "${ROOT_DIR}/frontend"

echo "Built images:"
echo "  ${BACKEND_IMAGE_NAME}:${APP_VERSION}"
echo "  ${BACKEND_IMAGE_NAME}:latest"
echo "  ${BACKEND_IMAGE_REF}:${APP_VERSION}"
echo "  ${BACKEND_IMAGE_REF}:latest"
echo "  ${MIGRATOR_IMAGE_NAME}:${APP_VERSION}"
echo "  ${MIGRATOR_IMAGE_NAME}:latest"
echo "  ${MIGRATOR_IMAGE_REF}:${APP_VERSION}"
echo "  ${MIGRATOR_IMAGE_REF}:latest"
echo "  ${FRONTEND_IMAGE_NAME}:${APP_VERSION}"
echo "  ${FRONTEND_IMAGE_NAME}:latest"
echo "  ${FRONTEND_IMAGE_REF}:${APP_VERSION}"
echo "  ${FRONTEND_IMAGE_REF}:latest"

if [[ "${PUSH_IMAGES}" == "true" ]]; then
    echo "Publishing image refs"
    push_image_ref "${BACKEND_IMAGE_REF}"
    push_image_ref "${MIGRATOR_IMAGE_REF}"
    push_image_ref "${FRONTEND_IMAGE_REF}"
fi
