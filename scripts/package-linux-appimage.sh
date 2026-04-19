#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${VERSION:-dev}"
ARCH="${ARCH:-x86_64}"
GOARCH="amd64"
OUT_DIR="$ROOT_DIR/dist/release"
TOOL_DIR="$ROOT_DIR/dist/tools"
APPDIR="$OUT_DIR/BucketDesk.AppDir"
APPIMAGETOOL="$TOOL_DIR/appimagetool-${ARCH}.AppImage"

if [[ "$ARCH" == "x86_64" ]]; then
  GOARCH="amd64"
else
  echo "Only x86_64 AppImage builds are supported in the default workflow." >&2
  exit 1
fi

mkdir -p "$OUT_DIR" "$TOOL_DIR"
rm -rf "$APPDIR"
mkdir -p "$APPDIR/usr/bin" "$APPDIR/usr/share/applications" "$APPDIR/usr/share/icons/hicolor/512x512/apps"

GOOS=linux GOARCH="$GOARCH" CGO_ENABLED=0 go build \
  -trimpath \
  -ldflags="-s -w -X main.version=$VERSION" \
  -o "$APPDIR/usr/bin/bucketdesk" \
  "$ROOT_DIR/cmd/bucketdesk"

cp "$ROOT_DIR/assets/bucketdesk.png" "$APPDIR/usr/share/icons/hicolor/512x512/apps/bucketdesk.png"
cp "$ROOT_DIR/assets/bucketdesk.png" "$APPDIR/bucketdesk.png"

cat > "$APPDIR/bucketdesk.desktop" <<DESKTOP
[Desktop Entry]
Type=Application
Name=BucketDesk
Comment=Bilingual MinIO and S3-compatible bucket manager
Exec=bucketdesk
Icon=bucketdesk
Categories=Utility;Network;
Terminal=false
DESKTOP

cat > "$APPDIR/AppRun" <<'APPRUN'
#!/usr/bin/env bash
HERE="$(dirname "$(readlink -f "$0")")"
exec "$HERE/usr/bin/bucketdesk" "$@"
APPRUN
chmod +x "$APPDIR/AppRun"

if [[ ! -x "$APPIMAGETOOL" ]]; then
  curl --fail --location --retry 5 --retry-all-errors --retry-delay 2 \
    "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-${ARCH}.AppImage" \
    -o "$APPIMAGETOOL"
  chmod +x "$APPIMAGETOOL"
fi

ARCH="$ARCH" APPIMAGE_EXTRACT_AND_RUN=1 "$APPIMAGETOOL" "$APPDIR" "$OUT_DIR/BucketDesk_${VERSION}_linux_${ARCH}.AppImage"
rm -rf "$APPDIR"
