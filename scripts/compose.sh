#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONTAINER_CLI="${CONTAINER_CLI:-podman}"

if ! command -v "${CONTAINER_CLI}" >/dev/null 2>&1; then
    echo "Container CLI not found: ${CONTAINER_CLI}" >&2
    exit 1
fi

if "${CONTAINER_CLI}" compose version >/dev/null 2>&1; then
    exec "${CONTAINER_CLI}" compose -f "${ROOT_DIR}/compose.yml" "$@"
fi

if [[ "${CONTAINER_CLI}" == "docker" ]] && command -v docker-compose >/dev/null 2>&1; then
    exec docker-compose -f "${ROOT_DIR}/compose.yml" "$@"
fi

if [[ "${CONTAINER_CLI}" == "podman" ]] && command -v podman-compose >/dev/null 2>&1; then
    exec podman-compose -f "${ROOT_DIR}/compose.yml" "$@"
fi

echo "No compose provider found for ${CONTAINER_CLI}. Install ${CONTAINER_CLI} compose support or set CONTAINER_CLI to a runtime with compose available." >&2
exit 1
