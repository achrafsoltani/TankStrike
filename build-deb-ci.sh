#!/usr/bin/env bash
set -euo pipefail

# CI .deb builder — takes version as argument, expects binary already built
VERSION="${1:?Usage: build-deb-ci.sh VERSION}"
APP_NAME="tankstrike"
ARCH="amd64"
PKG_DIR="${APP_NAME}_${VERSION}_${ARCH}"

echo "==> Packaging TankStrike v${VERSION} .deb"

if [ ! -f "${APP_NAME}" ]; then
    echo "Error: binary '${APP_NAME}' not found — run 'go build' first" >&2
    exit 1
fi

BINARY_SIZE=$(stat --format="%s" "${APP_NAME}")
INSTALLED_SIZE=$(( BINARY_SIZE / 1024 ))

# Create package structure
rm -rf "${PKG_DIR}"
mkdir -p "${PKG_DIR}/DEBIAN"
mkdir -p "${PKG_DIR}/usr/bin"
mkdir -p "${PKG_DIR}/usr/share/applications"

cp "${APP_NAME}" "${PKG_DIR}/usr/bin/${APP_NAME}"
chmod 755 "${PKG_DIR}/usr/bin/${APP_NAME}"

# Desktop entry
cat > "${PKG_DIR}/usr/share/applications/${APP_NAME}.desktop" <<'DESKTOP'
[Desktop Entry]
Type=Application
Name=TankStrike
Comment=A modern Battle City remake
Exec=tankstrike
Icon=tankstrike
Terminal=false
Categories=Game;ActionGame;
Keywords=tank;battle;city;arcade;
DESKTOP

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
 types, 6 power-ups, procedural audio, and particle effects.
CONTROL

dpkg-deb --build --root-owner-group "${PKG_DIR}"
rm -rf "${PKG_DIR}"

echo "==> Done: ${PKG_DIR}.deb"
