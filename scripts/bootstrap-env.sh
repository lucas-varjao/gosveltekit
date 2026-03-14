#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

copy_if_missing() {
    local src="$1"
    local dest="$2"

    if [[ -f "$dest" ]]; then
        echo "Keeping existing ${dest#$ROOT_DIR/}"
        return
    fi

    cp "$src" "$dest"
    echo "Created ${dest#$ROOT_DIR/} from example"
}

copy_if_missing "${ROOT_DIR}/.env.example" "${ROOT_DIR}/.env"
copy_if_missing "${ROOT_DIR}/backend/.env.example" "${ROOT_DIR}/backend/.env"
copy_if_missing "${ROOT_DIR}/frontend/.env.example" "${ROOT_DIR}/frontend/.env"
