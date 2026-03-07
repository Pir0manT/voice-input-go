#!/bin/bash
set -e

# Ensure Go is in PATH (for WSL with manual Go install)
export PATH="/usr/local/go/bin:$HOME/go/bin:$PATH"

APP_NAME="voice-input-go"
OUTPUT="${APP_NAME}"

# На Linux полная статическая линковка невозможна (GTK, ayatana-appindicator и др.)
# Линкуем динамически, требуются пакеты:
#   sudo apt install libportaudio2 libasound2-dev libayatana-appindicator3-dev libgtk-3-dev
LDFLAGS="-s -w"

echo "[1/2] Building ${OUTPUT}..."
go build -ldflags="$LDFLAGS" -o "$OUTPUT" ./cmd/${APP_NAME}/

echo "[2/2] Compressing with UPX..."
if command -v upx &>/dev/null; then
    upx --best "$OUTPUT"
else
    echo "  UPX not found, skipping compression"
    echo "  Install: sudo apt install upx-ucl"
fi

echo ""
echo "Done:"
ls -lh "$OUTPUT"
