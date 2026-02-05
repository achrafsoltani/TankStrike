#!/usr/bin/env bash
set -euo pipefail

APP_NAME="tankstrike"
VERSION="1.0.0"
ARCH="amd64"
PKG_DIR="${APP_NAME}_${VERSION}_${ARCH}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "==> Building TankStrike v${VERSION} .deb package"

# Check dependencies
if ! command -v go &>/dev/null; then
    echo "Error: Go is not installed" >&2
    exit 1
fi

# Check that glow is available alongside
if [ ! -d "${SCRIPT_DIR}/../glow" ]; then
    echo "Error: glow must be cloned alongside TankStrike (expected at ../glow)" >&2
    echo "  git clone https://github.com/AchrafSoltani/glow.git ${SCRIPT_DIR}/../glow" >&2
    exit 1
fi

# Build the binary
echo "==> Compiling binary..."
cd "${SCRIPT_DIR}"
go build -ldflags="-s -w" -o "${APP_NAME}" main.go
BINARY_SIZE=$(stat --format="%s" "${APP_NAME}")
INSTALLED_SIZE=$(( BINARY_SIZE / 1024 ))

# Create package structure
echo "==> Creating package structure..."
rm -rf "${PKG_DIR}"
mkdir -p "${PKG_DIR}/DEBIAN"
mkdir -p "${PKG_DIR}/usr/bin"
mkdir -p "${PKG_DIR}/usr/share/applications"
mkdir -p "${PKG_DIR}/usr/share/icons/hicolor/128x128/apps"

# Install binary
cp "${APP_NAME}" "${PKG_DIR}/usr/bin/${APP_NAME}"
chmod 755 "${PKG_DIR}/usr/bin/${APP_NAME}"

# Desktop entry
cat > "${PKG_DIR}/usr/share/applications/${APP_NAME}.desktop" <<DESKTOP
[Desktop Entry]
Type=Application
Name=TankStrike
Comment=A modern Battle City remake
Exec=${APP_NAME}
Icon=${APP_NAME}
Terminal=false
Categories=Game;ActionGame;
Keywords=tank;battle;city;arcade;
DESKTOP

# Generate a simple icon (PPM → PNG via Go or fallback)
# If screenshot.png exists, use a cropped version; otherwise skip
if [ -f "${SCRIPT_DIR}/screenshot.png" ] && command -v ffmpeg &>/dev/null; then
    echo "==> Generating application icon..."
    ffmpeg -y -i "${SCRIPT_DIR}/screenshot.png" \
        -vf "crop=128:128:368:272" \
        "${PKG_DIR}/usr/share/icons/hicolor/128x128/apps/${APP_NAME}.png" \
        2>/dev/null || true
fi

# Control file
cat > "${PKG_DIR}/DEBIAN/control" <<CONTROL
Package: ${APP_NAME}
Version: ${VERSION}
Section: games
Priority: optional
Architecture: ${ARCH}
Depends: libx11-6 | xwayland
Installed-Size: ${INSTALLED_SIZE}
Maintainer: Achraf Soltani <contact@achrafsoltani.me>
Homepage: https://github.com/AchrafSoltani/TankStrike
Description: A modern Battle City remake
 TankStrike is a modern recreation of the classic Battle City (NES) game,
 built entirely in Go with the Glow engine. Features 10 levels, 4 enemy
 types, 6 power-ups, procedural audio, and particle effects. No external
 assets — everything is drawn with primitives.
CONTROL

# Build the .deb
echo "==> Packaging ${PKG_DIR}.deb..."
dpkg-deb --build --root-owner-group "${PKG_DIR}"

# Clean up
rm -rf "${PKG_DIR}"
rm -f "${APP_NAME}"

echo "==> Done: ${PKG_DIR}.deb"
echo "    Install with: sudo dpkg -i ${PKG_DIR}.deb"
