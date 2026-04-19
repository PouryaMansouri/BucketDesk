#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${VERSION:-dev}"
CLEAN_VERSION="${VERSION#v}"
ARCH="${ARCH:-arm64}"
OUT_DIR="$ROOT_DIR/dist/release"
STAGE_DIR="$OUT_DIR/dmg-${ARCH}"
APP_DIR="$STAGE_DIR/BucketDesk.app"
DMG_PATH="$OUT_DIR/BucketDesk_${VERSION}_macos_${ARCH}.dmg"

mkdir -p "$OUT_DIR"
rm -rf "$STAGE_DIR" "$DMG_PATH"
mkdir -p "$APP_DIR/Contents/MacOS" "$APP_DIR/Contents/Resources"

cat > "$APP_DIR/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleName</key>
  <string>BucketDesk</string>
  <key>CFBundleDisplayName</key>
  <string>BucketDesk</string>
  <key>CFBundleIdentifier</key>
  <string>com.pouryamansouri.bucketdesk</string>
  <key>CFBundleVersion</key>
  <string>${CLEAN_VERSION}</string>
  <key>CFBundleShortVersionString</key>
  <string>${CLEAN_VERSION}</string>
  <key>CFBundleExecutable</key>
  <string>BucketDesk</string>
  <key>CFBundleIconFile</key>
  <string>bucketdesk</string>
  <key>LSMinimumSystemVersion</key>
  <string>11.0</string>
  <key>LSUIElement</key>
  <false/>
</dict>
</plist>
PLIST

GOOS=darwin GOARCH="$ARCH" CGO_ENABLED=0 go build \
  -trimpath \
  -ldflags="-s -w -X main.version=$VERSION" \
  -o "$APP_DIR/Contents/MacOS/BucketDesk" \
  "$ROOT_DIR/cmd/bucketdesk"

chmod +x "$APP_DIR/Contents/MacOS/BucketDesk"
cp "$ROOT_DIR/LICENSE" "$APP_DIR/Contents/Resources/LICENSE"
cp "$ROOT_DIR/NOTICE" "$APP_DIR/Contents/Resources/NOTICE"
cp "$ROOT_DIR/assets/bucketdesk.icns" "$APP_DIR/Contents/Resources/bucketdesk.icns"

if [[ -n "${MACOS_CODESIGN_IDENTITY:-}" ]]; then
  codesign --force --deep --options runtime --timestamp --sign "$MACOS_CODESIGN_IDENTITY" "$APP_DIR"
fi

hdiutil create \
  -volname "BucketDesk" \
  -srcfolder "$STAGE_DIR" \
  -ov \
  -format UDZO \
  "$DMG_PATH"

if [[ -n "${MACOS_CODESIGN_IDENTITY:-}" ]]; then
  codesign --force --timestamp --sign "$MACOS_CODESIGN_IDENTITY" "$DMG_PATH"
fi

rm -rf "$STAGE_DIR"
echo "$DMG_PATH"
