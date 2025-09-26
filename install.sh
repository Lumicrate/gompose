#!/usr/bin/env bash
set -euo pipefail

REPO="github.com/Lumicrate/gompose"
BIN_NAME="gompose"
BIN_DIR="$HOME/.local/bin"

echo "Installing $BIN_NAME from $REPO..."

if ! command -v go >/dev/null 2>&1; then
  echo "Go is not installed. Please install Go first: https://go.dev/dl/"
  exit 1
fi

echo "Running: go install $REPO@latest"
go install "$REPO@latest"

mkdir -p "$BIN_DIR"

GOBIN=$(go env GOPATH)/bin/$BIN_NAME

if [ ! -f "$GOBIN" ]; then
  echo "Failed to build $BIN_NAME binary at $GOBIN"
  exit 1
fi

cp "$GOBIN" "$BIN_DIR/"

echo "Installed $BIN_NAME to $BIN_DIR"

if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
  echo "⚠️  Note: $BIN_DIR is not in your PATH."
  echo "   Add this line to your shell config (~/.bashrc or ~/.zshrc):"
  echo "     export PATH=\"\$PATH:$BIN_DIR\""
  echo "   Then run: source ~/.bashrc (or restart your terminal)."
else
  echo "You can now run '$BIN_NAME' globally!"
fi
