#!/bin/bash
set -euo pipefail

# Build macOS .app bundle and .dmg installer
# Usage: ./build-dmg.sh <binary-path> <version>

BINARY="${1:?Usage: build-dmg.sh <binary> <version>}"
VERSION="${2:-0.0.0-dev}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
APP_NAME="Voice Input"
BUNDLE_NAME="${APP_NAME}.app"
DMG_NAME="voice-input-go-macos-arm64-${VERSION}.dmg"

echo "=== Building ${BUNDLE_NAME} v${VERSION} ==="

# --- 1. Generate .icns from SVG ---
echo "Generating app icon..."
ICONSET_DIR=$(mktemp -d)/AppIcon.iconset
mkdir -p "$ICONSET_DIR"

SVG_PATH="${SCRIPT_DIR}/../internal/tray/icon.svg"

# Render SVG to large PNG (1024x1024) via rsvg-convert or sips
if command -v rsvg-convert &>/dev/null; then
    rsvg-convert -w 1024 -h 1024 "$SVG_PATH" -o /tmp/icon_1024.png
elif command -v python3 &>/dev/null; then
    # Fallback: use cairosvg via python
    python3 -c "
import subprocess, sys
try:
    import cairosvg
    cairosvg.svg2png(url='$SVG_PATH', write_to='/tmp/icon_1024.png', output_width=1024, output_height=1024)
except ImportError:
    subprocess.run([sys.executable, '-m', 'pip', 'install', 'cairosvg', '-q'])
    import cairosvg
    cairosvg.svg2png(url='$SVG_PATH', write_to='/tmp/icon_1024.png', output_width=1024, output_height=1024)
"
else
    echo "ERROR: need rsvg-convert (brew install librsvg) or python3 with cairosvg"
    exit 1
fi

# Generate all required sizes for .iconset
for size in 16 32 64 128 256 512; do
    sips -z $size $size /tmp/icon_1024.png --out "${ICONSET_DIR}/icon_${size}x${size}.png" >/dev/null
    double=$((size * 2))
    if [ $double -le 1024 ]; then
        sips -z $double $double /tmp/icon_1024.png --out "${ICONSET_DIR}/icon_${size}x${size}@2x.png" >/dev/null
    fi
done
cp /tmp/icon_1024.png "${ICONSET_DIR}/icon_512x512@2x.png"

# Convert iconset to icns
iconutil -c icns "$ICONSET_DIR" -o /tmp/AppIcon.icns
echo "Icon created: /tmp/AppIcon.icns"

# --- 2. Build .app bundle ---
echo "Building app bundle..."
rm -rf "${BUNDLE_NAME}"
mkdir -p "${BUNDLE_NAME}/Contents/MacOS"
mkdir -p "${BUNDLE_NAME}/Contents/Resources"

# Copy binary
cp "$BINARY" "${BUNDLE_NAME}/Contents/MacOS/voice-input-go"
chmod +x "${BUNDLE_NAME}/Contents/MacOS/voice-input-go"

# Copy icon
cp /tmp/AppIcon.icns "${BUNDLE_NAME}/Contents/Resources/AppIcon.icns"

# Generate Info.plist with version
sed "s/\${VERSION}/${VERSION}/g" "${SCRIPT_DIR}/Info.plist" > "${BUNDLE_NAME}/Contents/Info.plist"

echo "Bundle created: ${BUNDLE_NAME}"

# --- 3. Build .dmg ---
echo "Building DMG..."
rm -f "$DMG_NAME"

# Create temp dir for DMG contents
DMG_DIR=$(mktemp -d)
cp -R "${BUNDLE_NAME}" "${DMG_DIR}/"
ln -s /Applications "${DMG_DIR}/Applications"

hdiutil create -volname "Voice Input" \
    -srcfolder "$DMG_DIR" \
    -ov -format UDZO \
    "$DMG_NAME"

rm -rf "$DMG_DIR"

echo "=== Done: ${DMG_NAME} ==="
ls -lh "$DMG_NAME"
