#!/usr/bin/env bash
# Builds "TiBCO Obfuscation Tool.app" for macOS.
# Run this in macOS Terminal (not the Claude VM):  cd ~/dev/tibco-gui && ./build-macos.sh
set -euo pipefail

APP_NAME="TiBCO Obfuscation Tool"
APP_ID="com.diogo.tibcoobfuscationtool"

cd "$(dirname "$0")"

# 1. Locate Go
if ! command -v go >/dev/null 2>&1; then
  for p in /usr/local/go/bin /opt/homebrew/bin "$HOME/go/bin"; do
    [ -x "$p/go" ] && export PATH="$p:$PATH" && break
  done
fi
command -v go >/dev/null 2>&1 || { echo "ERROR: Go not found. Install from https://go.dev/dl or 'brew install go'."; exit 1; }
echo "Using $(go version)"

# 2. Xcode command-line tools (needed because Fyne uses CGO)
if ! xcode-select -p >/dev/null 2>&1; then
  echo "ERROR: Xcode command-line tools missing. Run: xcode-select --install"; exit 1
fi

export PATH="$PATH:$(go env GOPATH)/bin"

# 3. Resolve dependencies and install the fyne packager if absent
go mod tidy
command -v fyne >/dev/null 2>&1 || go install fyne.io/fyne/v2/cmd/fyne@latest

# 4. Package the .app (native arch). Add arches for a universal binary below.
rm -rf "$APP_NAME.app"
fyne package -os darwin -name "$APP_NAME" -appID "$APP_ID" -icon Icon.png

echo
echo "Built: $(pwd)/$APP_NAME.app"
echo "Open it with:  open \"$APP_NAME.app\""
echo "Install it with:  cp -R \"$APP_NAME.app\" /Applications/"
