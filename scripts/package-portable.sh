#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${VERSION:-dev}"
OUT_DIR="$ROOT_DIR/dist/release"

mkdir -p "$OUT_DIR"

build_one() {
  local goos="$1"
  local goarch="$2"
  local ext=""
  local archive_ext="tar.gz"

  if [[ "$goos" == "windows" ]]; then
    ext=".exe"
    archive_ext="zip"
  fi

  local name="bucketdesk_${VERSION}_${goos}_${goarch}"
  local work="$OUT_DIR/$name"
  rm -rf "$work"
  mkdir -p "$work"

  echo "Building $name"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o "$work/bucketdesk$ext" \
    "$ROOT_DIR/cmd/bucketdesk"

  cp "$ROOT_DIR/README.md" "$work/README.md"
  cp "$ROOT_DIR/LICENSE" "$work/LICENSE"
  cp "$ROOT_DIR/NOTICE" "$work/NOTICE"

  if [[ "$archive_ext" == "zip" ]]; then
    (cd "$OUT_DIR" && zip -qr "$name.zip" "$name")
  else
    (cd "$OUT_DIR" && tar -czf "$name.tar.gz" "$name")
  fi

  rm -rf "$work"
}

build_one darwin amd64
build_one darwin arm64
build_one linux amd64
build_one linux arm64
build_one windows amd64
build_one windows arm64

echo "Portable packages are in $OUT_DIR"
