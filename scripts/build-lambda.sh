#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="$ROOT_DIR/backend/dist"

mkdir -p "$OUT_DIR"

cd "$ROOT_DIR/backend"
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o "$OUT_DIR/bootstrap" ./cmd/lambda
cd "$OUT_DIR"
zip -q lambda.zip bootstrap

echo "Built Lambda artifact: $OUT_DIR/lambda.zip"
