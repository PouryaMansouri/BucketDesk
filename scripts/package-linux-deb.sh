#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${VERSION:-0.1.0}"
ARCH="${ARCH:-amd64}"
GOARCH="$ARCH"
OUT_DIR="$ROOT_DIR/dist/release"
PKG_DIR="$OUT_DIR/bucketdesk_${VERSION}_linux_${ARCH}"

if [[ "$ARCH" == "amd64" ]]; then
  GOARCH="amd64"
elif [[ "$ARCH" == "arm64" ]]; then
  GOARCH="arm64"
else
  echo "Unsupported ARCH: $ARCH" >&2
  exit 1
fi

rm -rf "$PKG_DIR"
mkdir -p "$PKG_DIR/DEBIAN" "$PKG_DIR/usr/bin" "$PKG_DIR/usr/share/doc/bucketdesk"

GOOS=linux GOARCH="$GOARCH" CGO_ENABLED=0 go build \
  -trimpath \
  -ldflags="-s -w -X main.version=$VERSION" \
  -o "$PKG_DIR/usr/bin/bucketdesk" \
  "$ROOT_DIR/cmd/bucketdesk"

cat > "$PKG_DIR/DEBIAN/control" <<CONTROL
Package: bucketdesk
Version: ${VERSION#v}
Section: utils
Priority: optional
Architecture: $ARCH
Maintainer: Pourya Mansouri
Description: Bilingual MinIO and S3-compatible bucket manager
 BucketDesk lets users manage scoped buckets without access to the MinIO Console.
CONTROL

cp "$ROOT_DIR/README.md" "$PKG_DIR/usr/share/doc/bucketdesk/README.md"
cp "$ROOT_DIR/LICENSE" "$PKG_DIR/usr/share/doc/bucketdesk/LICENSE"
cp "$ROOT_DIR/NOTICE" "$PKG_DIR/usr/share/doc/bucketdesk/NOTICE"

dpkg-deb --build "$PKG_DIR" "$OUT_DIR/bucketdesk_${VERSION}_linux_${ARCH}.deb"
rm -rf "$PKG_DIR"
