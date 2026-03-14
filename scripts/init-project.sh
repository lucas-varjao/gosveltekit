#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROJECT_ENV_FILE="${ROOT_DIR}/project.env"

if [[ ! -f "${PROJECT_ENV_FILE}" ]]; then
    echo "project.env not found at ${PROJECT_ENV_FILE}" >&2
    exit 1
fi

# shellcheck disable=SC1090
source "${PROJECT_ENV_FILE}"

CURRENT_APP_SLUG="${APP_SLUG}"
CURRENT_APP_DISPLAY_NAME="${APP_DISPLAY_NAME}"
CURRENT_APP_DESCRIPTION="${APP_DESCRIPTION}"
CURRENT_GO_MODULE="${GO_MODULE}"
CURRENT_CONTAINER_REGISTRY="${CONTAINER_REGISTRY}"
CURRENT_APP_DOMAIN="${APP_DOMAIN}"
CURRENT_K8S_NAMESPACE="${K8S_NAMESPACE}"
CURRENT_BACKEND_IMAGE_NAME="${BACKEND_IMAGE_NAME}"
CURRENT_FRONTEND_IMAGE_NAME="${FRONTEND_IMAGE_NAME}"
CURRENT_BACKEND_IMAGE_REF="${BACKEND_IMAGE_REF}"
CURRENT_FRONTEND_IMAGE_REF="${FRONTEND_IMAGE_REF}"
CURRENT_APP_SUPPORT_EMAIL="${APP_SUPPORT_EMAIL}"
CURRENT_APP_RESET_URL="${APP_RESET_URL}"

TARGET_APP_SLUG=""
TARGET_APP_DISPLAY_NAME=""
TARGET_GO_MODULE=""
TARGET_CONTAINER_REGISTRY=""
TARGET_APP_DOMAIN=""
TARGET_K8S_NAMESPACE=""

prompt_value() {
    local label="$1"
    local default_value="$2"
    local result=""

    if [[ -t 0 ]]; then
        read -r -p "${label} [${default_value}]: " result
        if [[ -z "${result}" ]]; then
            result="${default_value}"
        fi
    else
        result="${default_value}"
    fi

    printf '%s' "${result}"
}

normalize_registry() {
    local registry="$1"
    registry="${registry%/}"
    printf '%s' "${registry}"
}

has_stale_project_refs() {
    if command -v rg >/dev/null 2>&1; then
        rg -F -n \
            -e "${CURRENT_APP_SLUG}" \
            -e "${CURRENT_APP_DISPLAY_NAME}" \
            -e "${CURRENT_APP_DOMAIN}" \
            README.md backend frontend k8s compose.yml .env.example backend/.env.example frontend/.env.example Makefile project.env AGENTS.md >/dev/null 2>&1
        return
    fi

    grep -RIF -n \
        -e "${CURRENT_APP_SLUG}" \
        -e "${CURRENT_APP_DISPLAY_NAME}" \
        -e "${CURRENT_APP_DOMAIN}" \
        README.md backend frontend k8s compose.yml .env.example backend/.env.example frontend/.env.example Makefile project.env AGENTS.md >/dev/null 2>&1
}

replace_in_file() {
    local old_value="$1"
    local new_value="$2"
    local file="$3"

    if [[ "${old_value}" == "${new_value}" || -z "${old_value}" ]]; then
        return
    fi

    OLD_VALUE="${old_value}" NEW_VALUE="${new_value}" perl -0pi -e 's/\Q$ENV{OLD_VALUE}\E/$ENV{NEW_VALUE}/g' "${file}"
}

text_files() {
    find "${ROOT_DIR}" \
        -type f \
        ! -path "${ROOT_DIR}/.git/*" \
        ! -path "${ROOT_DIR}/frontend/node_modules/*" \
        ! -path "${ROOT_DIR}/frontend/build/*" \
        ! -path "${ROOT_DIR}/frontend/.svelte-kit/*" \
        ! -path "${ROOT_DIR}/backend/bin/*" \
        "$@"
}

replace_in_text_files() {
    local old_value="$1"
    local new_value="$2"

    while IFS= read -r -d '' file; do
        if grep -Iq . "${file}"; then
            replace_in_file "${old_value}" "${new_value}" "${file}"
        fi
    done < <(text_files -print0)
}

rewrite_project_env() {
    cat > "${PROJECT_ENV_FILE}" <<EOF
APP_SLUG=${TARGET_APP_SLUG}
APP_DISPLAY_NAME=${TARGET_APP_DISPLAY_NAME}
APP_DESCRIPTION=${CURRENT_APP_DESCRIPTION@Q}
GO_MODULE=${TARGET_GO_MODULE}
CONTAINER_REGISTRY=${TARGET_CONTAINER_REGISTRY}
APP_DOMAIN=${TARGET_APP_DOMAIN}
K8S_NAMESPACE=${TARGET_K8S_NAMESPACE}
BACKEND_IMAGE_NAME=${TARGET_BACKEND_IMAGE_NAME}
FRONTEND_IMAGE_NAME=${TARGET_FRONTEND_IMAGE_NAME}
BACKEND_IMAGE_REF=${TARGET_BACKEND_IMAGE_REF}
FRONTEND_IMAGE_REF=${TARGET_FRONTEND_IMAGE_REF}
APP_SUPPORT_EMAIL=${TARGET_APP_SUPPORT_EMAIL}
APP_RESET_URL=${TARGET_APP_RESET_URL}
LOCAL_FRONTEND_URL=${LOCAL_FRONTEND_URL}
LOCAL_API_URL=${LOCAL_API_URL}
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --app-name)
            TARGET_APP_SLUG="$2"
            shift 2
            ;;
        --display-name)
            TARGET_APP_DISPLAY_NAME="$2"
            shift 2
            ;;
        --go-module)
            TARGET_GO_MODULE="$2"
            shift 2
            ;;
        --container-registry)
            TARGET_CONTAINER_REGISTRY="$2"
            shift 2
            ;;
        --domain)
            TARGET_APP_DOMAIN="$2"
            shift 2
            ;;
        --k8s-namespace)
            TARGET_K8S_NAMESPACE="$2"
            shift 2
            ;;
        *)
            echo "Unknown argument: $1" >&2
            exit 1
            ;;
    esac
done

TARGET_APP_SLUG="${TARGET_APP_SLUG:-$(prompt_value "App name (slug)" "${CURRENT_APP_SLUG}")}"
TARGET_APP_DISPLAY_NAME="${TARGET_APP_DISPLAY_NAME:-$(prompt_value "Display name" "${CURRENT_APP_DISPLAY_NAME}")}"
TARGET_GO_MODULE="${TARGET_GO_MODULE:-$(prompt_value "Go module" "${CURRENT_GO_MODULE}")}"
TARGET_CONTAINER_REGISTRY="${TARGET_CONTAINER_REGISTRY:-$(prompt_value "Container registry" "${CURRENT_CONTAINER_REGISTRY}")}"
TARGET_APP_DOMAIN="${TARGET_APP_DOMAIN:-$(prompt_value "Public domain" "${CURRENT_APP_DOMAIN}")}"
TARGET_K8S_NAMESPACE="${TARGET_K8S_NAMESPACE:-$(prompt_value "Kubernetes namespace" "${CURRENT_K8S_NAMESPACE}")}"

TARGET_CONTAINER_REGISTRY="$(normalize_registry "${TARGET_CONTAINER_REGISTRY}")"
TARGET_BACKEND_IMAGE_NAME="${TARGET_APP_SLUG}-backend"
TARGET_FRONTEND_IMAGE_NAME="${TARGET_APP_SLUG}-frontend"
TARGET_BACKEND_IMAGE_REF="${TARGET_BACKEND_IMAGE_NAME}"
TARGET_FRONTEND_IMAGE_REF="${TARGET_FRONTEND_IMAGE_NAME}"

if [[ -n "${TARGET_CONTAINER_REGISTRY}" ]]; then
    TARGET_BACKEND_IMAGE_REF="${TARGET_CONTAINER_REGISTRY}/${TARGET_BACKEND_IMAGE_NAME}"
    TARGET_FRONTEND_IMAGE_REF="${TARGET_CONTAINER_REGISTRY}/${TARGET_FRONTEND_IMAGE_NAME}"
fi

TARGET_APP_SUPPORT_EMAIL="no-reply@${TARGET_APP_DOMAIN}"
TARGET_APP_RESET_URL="https://${TARGET_APP_DOMAIN}/reset-password?token="

if [[ -z "${TARGET_APP_SLUG}" || -z "${TARGET_APP_DISPLAY_NAME}" || -z "${TARGET_GO_MODULE}" ]]; then
    echo "app-name, display-name and go-module cannot be empty" >&2
    exit 1
fi

if [[ ! -d "${ROOT_DIR}/backend" || ! -d "${ROOT_DIR}/frontend" ]]; then
    echo "init-project.sh must run from the template repository root" >&2
    exit 1
fi

(
    cd "${ROOT_DIR}/backend"
    go mod edit -module "${TARGET_GO_MODULE}"
)

replace_in_text_files "${CURRENT_GO_MODULE}/" "${TARGET_GO_MODULE}/"
replace_in_text_files "${CURRENT_BACKEND_IMAGE_REF}" "${TARGET_BACKEND_IMAGE_REF}"
replace_in_text_files "${CURRENT_FRONTEND_IMAGE_REF}" "${TARGET_FRONTEND_IMAGE_REF}"
replace_in_text_files "${CURRENT_BACKEND_IMAGE_NAME}" "${TARGET_BACKEND_IMAGE_NAME}"
replace_in_text_files "${CURRENT_FRONTEND_IMAGE_NAME}" "${TARGET_FRONTEND_IMAGE_NAME}"
replace_in_text_files "${CURRENT_APP_RESET_URL}" "${TARGET_APP_RESET_URL}"
replace_in_text_files "${CURRENT_APP_SUPPORT_EMAIL}" "${TARGET_APP_SUPPORT_EMAIL}"
replace_in_text_files "${CURRENT_APP_DOMAIN}" "${TARGET_APP_DOMAIN}"
replace_in_text_files "${CURRENT_K8S_NAMESPACE}" "${TARGET_K8S_NAMESPACE}"
replace_in_text_files "${CURRENT_APP_DISPLAY_NAME}" "${TARGET_APP_DISPLAY_NAME}"
replace_in_text_files "${CURRENT_CONTAINER_REGISTRY}" "${TARGET_CONTAINER_REGISTRY}"
replace_in_text_files "${CURRENT_APP_SLUG}" "${TARGET_APP_SLUG}"

rewrite_project_env

CURRENT_MANIFEST="${ROOT_DIR}/k8s/${CURRENT_APP_SLUG}.yaml"
TARGET_MANIFEST="${ROOT_DIR}/k8s/${TARGET_APP_SLUG}.yaml"

if [[ -f "${CURRENT_MANIFEST}" && "${CURRENT_MANIFEST}" != "${TARGET_MANIFEST}" ]]; then
    mv "${CURRENT_MANIFEST}" "${TARGET_MANIFEST}"
fi

(
    cd "${ROOT_DIR}/backend"
    go mod tidy
)

if has_stale_project_refs; then
    echo "Template initialization left stale references to the previous project metadata." >&2
    exit 1
fi

echo "Initialized template for ${TARGET_APP_DISPLAY_NAME}"
