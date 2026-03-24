#!/usr/bin/env bash
set -euo pipefail

APP_NAME="countdown"
SRC_DIR="$(cd "$(dirname "$0")" && pwd)"
BUILD_DIR="$SRC_DIR/build"
INSTALL_DIR="$HOME/.local/bin"

# ── Ensure Go is installed ─────────────────────────────────────────────────
if ! command -v go &>/dev/null; then
  echo "❌  Go is not installed."
  echo "    Install it from https://go.dev/dl/ then re-run this script."
  exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✅  Go found: $GO_VERSION"

# ── Build ──────────────────────────────────────────────────────────────────
mkdir -p "$BUILD_DIR"
echo "🔨  Building $APP_NAME …"

cd "$SRC_DIR"

# go.mod (create if missing)
if [ ! -f go.mod ]; then
  go mod init countdown-menubar
fi

CGO_ENABLED=1 GOOS=darwin go build -o "$BUILD_DIR/$APP_NAME" .

echo "✅  Build succeeded → $BUILD_DIR/$APP_NAME"

# ── Optional: install to PATH ──────────────────────────────────────────────
read -r -p "📦  Install to $INSTALL_DIR so you can run it anywhere? [Y/n] " yn
yn="${yn:-Y}"
if [[ "$yn" =~ ^[Yy] ]]; then
  mkdir -p "$INSTALL_DIR"
  cp "$BUILD_DIR/$APP_NAME" "$INSTALL_DIR/$APP_NAME"
  chmod +x "$INSTALL_DIR/$APP_NAME"
  echo "✅  Installed to $INSTALL_DIR/$APP_NAME"
  echo ""
  echo "    Make sure $INSTALL_DIR is in your PATH:"
  echo '    export PATH="$HOME/.local/bin:$PATH"'
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀  HOW TO USE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "  Run with a target date:"
echo '    ./build/countdown 2025-12-31 "New Year 2026"'
echo ""
echo "  Run in the background (persists after terminal close):"
echo '    nohup ./build/countdown 2025-12-31 "New Year" &'
echo ""
echo "  Stop it:"
echo '    pkill countdown'
echo ""
